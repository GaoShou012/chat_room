package main

import (
	"fmt"
	"github.com/micro/cli/v2"
	"github.com/micro/go-micro/v2"
	"time"
	"wchatv1/config"
	proto_room "wchatv1/proto/room"
	service_room "wchatv1/services/room"
	"wchatv1/utils"
)

var (
	prometheusAddr string
)

var loadFlags = micro.Action(func(c *cli.Context) error {
	prometheusAddr = c.String("prometheus_address")
	return nil
})

func main() {
	service := micro.NewService(
		micro.Name(config.RoomServiceConfig.ServiceName()),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*10),
		micro.Flags(
			&cli.StringFlag{Name: "prometheus_address", Usage: "The prometheus service"},
		),
	)
	service.Init(loadFlags)

	go utils.Prometheus(prometheusAddr)

	utils.Micro.Init(service)
	utils.Micro.LoadSource()
	utils.Micro.LoadConfigMust(config.RedisClusterConfig)
	utils.Micro.LoadConfigMust(config.KafkaClusterConfig)
	utils.Micro.LoadConfigMust(config.RoomServiceConfig)

	service_room.InitCodec()
	service_room.InitRedisClusterClient(config.RedisClusterConfig.Addr, config.RedisClusterConfig.Password)
	broadcaster := &service_room.BroadcastToFrontierByKafka{}
	if err := broadcaster.Init(config.KafkaClusterConfig.Addr, config.RoomServiceConfig.Topic); err != nil {
		panic(err)
	}

	handler := &service_room.Service{
		Key:                 []byte("sadljkfslkjfa"),
		BroadcastToFrontier: broadcaster.Bucket(),
	}

	fmt.Println("启动服务")
	if err := proto_room.RegisterRoomServiceHandler(service.Server(), handler); err != nil {
		panic(err)
	}
	if err := service.Run(); err != nil {
		panic(err)
	}
}
