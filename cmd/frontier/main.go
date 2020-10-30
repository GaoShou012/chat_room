package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/GaoShou012/frontier"
	sarama_cluster "github.com/bsm/sarama-cluster"
	"github.com/gobwas/ws"
	"github.com/golang/glog"
	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	uuid "github.com/satori/go.uuid"
	"log"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"
	"wchatv1/config"
	proto_room "wchatv1/proto/room"
	"wchatv1/src/room"
	"wchatv1/utils"
)

/*
	监听kafka,group是根据frontierId进行监听
	所以每个frontierId都必须要不一样
*/

var (
	prometheusAddr string
)

var loadFlags = micro.Action(func(c *cli.Context) error {
	prometheusAddr = c.String("prometheus_address")
	return nil
})

func main() {
	service := micro.NewService(
		micro.Flags(
			&cli.StringFlag{Name: "prometheus_address", Usage: "The prometheus service"},
		),
	)
	service.Init(loadFlags)

	go utils.Prometheus(prometheusAddr)

	utils.Micro.Init(service)
	utils.Micro.LoadSource()
	utils.Micro.LoadConfigMust(config.FrontierConfig)
	utils.Micro.LoadConfigMust(config.KafkaClusterConfig)
	utils.Micro.LoadConfigMust(config.RedisClusterConfig)
	utils.Micro.LoadConfigMust(config.RoomServiceConfig)

	dynamicParams := &frontier.DynamicParams{
		LogLevel:         config.FrontierConfig.LogLevel,
		HeartbeatTimeout: config.FrontierConfig.HeartbeatTimeout,
		WriterBufferSize: config.FrontierConfig.WriterBufferSize,
		ReaderBufferSize: config.FrontierConfig.ReaderBufferSize,
		WriterTimeout:    time.Millisecond * 40,
		ReaderTimeout:    time.Microsecond * 10,
	}
	handler := &frontier.Handler{
		OnRequest: func(conn frontier.Conn, uri []byte) error {
			var token string
			var lastMessageId string

			u, err := url.Parse(string(uri))
			if err != nil {
				glog.Errorln(err)
				return err
			}

			if str := u.Query().Get("token"); str != "" {
				token = str
			} else {
				return errors.New("token is empty")
			}
			if str := u.Query().Get("mid"); str != "" {
				lastMessageId = str
			}

			req := &proto_room.ViewPassTokenReq{Token: token}
			rsp, err := config.RoomServiceConfig.ServiceClient().ViewPassToken(context.TODO(), req)
			if err != nil {
				glog.Errorln(err)
				return err
			}
			room.Agent.Join(conn, rsp.PassToken, lastMessageId)
			return nil
		},
		OnHost:          nil,
		OnHeader:        nil,
		OnBeforeUpgrade: nil,
		OnOpen: func(conn frontier.Conn) error {
			return nil
		},
		OnMessage: func(conn frontier.Conn, message []byte) {
			room.Agent.OnMessage(conn, message)
		},
		OnClose: func(conn frontier.Conn) {
			room.Agent.Leave(conn)
		},
	}

	id := uuid.NewV4().String()
	addr := service.Server().Options().Address
	maxConnections := 1000000

	go func() {
		fmt.Println("frontier heartbeat")
		ticker := time.NewTicker(time.Second * 10)
		for {
			req := &proto_room.FrontierPingReq{FrontierId: id}
			_, err := config.RoomServiceConfig.ServiceClient().FrontierPing(context.TODO(), req)
			if err != nil {
				glog.Errorln(err)
			}
			<-ticker.C
		}
	}()

	go func() {
		fmt.Println("watching frontier config")
		watcher, err := utils.Micro.Config.Watch("micro", "config", "frontier")
		if err != nil {
			panic(err)
		}
		for {
			val, err := watcher.Next()
			if err != nil {
				glog.Errorln(err)
				continue
			}
			fmt.Println("New Frontier Config", string(val.Bytes()))
			err = json.Unmarshal(val.Bytes(), config.FrontierConfig)
			if err != nil {
				glog.Errorln(err)
				continue
			}

			dynamicParams.LogLevel = config.FrontierConfig.LogLevel
			dynamicParams.HeartbeatTimeout = config.FrontierConfig.HeartbeatTimeout
			dynamicParams.WriterBufferSize = config.FrontierConfig.WriterBufferSize
			dynamicParams.ReaderBufferSize = config.FrontierConfig.ReaderBufferSize
		}
	}()

	{
		conf := sarama_cluster.NewConfig()
		conf.Consumer.Offsets.AutoCommit.Enable = true
		conf.Consumer.Offsets.AutoCommit.Interval = time.Second
		conf.Consumer.Offsets.CommitInterval = time.Second
		conf.Consumer.Return.Errors = true
		conf.Group.Return.Notifications = true

		addr, topics := config.KafkaClusterConfig.Addr, []string{config.RoomServiceConfig.Topic}
		consumer, err := sarama_cluster.NewConsumer(addr, id, topics, conf)
		if err != nil {
			panic(err)
		}
		
		// trap SIGINT to trigger a shutdown.
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt)

		// consume errors
		go func() {
			for err := range consumer.Errors() {
				log.Printf("Error: %s\n", err.Error())
			}
		}()

		// consume notifications
		go func() {
			for ntf := range consumer.Notifications() {
				log.Printf("Rebalanced: %+v\n", ntf)
			}
		}()

		// consume messages, watch signals
		go func() {
			for {
				select {
				case message, ok := <-consumer.Messages():
					if ok {
						//fmt.Fprintf(os.Stdout, "%s/%d/%d\t%s\t%s\n", msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value)
						room.Agent.BroadcastCache <- message.Value
						consumer.MarkOffset(message, "") // mark message as processed
					}
				case <-signals:
					return
				}
			}
		}()
	}

	f := &frontier.Frontier{
		Id:              id,
		Addr:            addr,
		MaxConnections:  maxConnections,
		DynamicParams:   dynamicParams,
		Protocol:        &frontier.ProtocolWs{MessageType: ws.OpText},
		Handler:         handler,
		SenderParallel:  100,
		SenderCacheSize: 1000,
	}
	if err := f.Init(); err != nil {
		panic(err)
	}
	if err := f.Start(); err != nil {
		panic(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		switch s := <-c; s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			glog.Infof("got signal %s; stop server", s)
		case syscall.SIGHUP:
			glog.Infof("got signal %s; go to deamon", s)
			continue
		}
		if err := f.Stop(); err != nil {
			glog.Errorf("stop server error: %v", err)
		}
		break
	}
}
