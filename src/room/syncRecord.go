package room

import (
	"container/list"
	"context"
	"github.com/golang/glog"
	"github.com/micro/go-micro/v2/client"
	"runtime"
	"time"
	"wchatv1/config"
	proto_room "wchatv1/proto/room"
)

const (
	PullNumberOfMessagesPerTime   = 50
	PublishNumberOfMessagePreTime = 20
)

var SyncRecord syncRecord

type syncRecord struct {
	pullingRecordClients *list.List
	pullingRecordChannel chan *Client

	publishRecordClients *list.List
	publishRecordChannel chan *Client

	publishCacheClients *list.List
	publishCacheChannel chan *Client
}

func (s *syncRecord) NewClient(cli *Client) {
	ClientPullingRecordGauge.Inc()

	s.pullingRecordClients.PushBack(cli)

	//defer s.publishRecordClientsMutex.Unlock()
	//s.publishRecordClientsMutex.Lock()
	s.publishRecordClients.PushBack(cli)
}

// pullerParallel 同步消息记录的并行数
func (s *syncRecord) Init(pullerParallel int) {
	s.pullingRecordClients = list.New()
	s.pullingRecordChannel = make(chan *Client, pullerParallel)
	s.publishRecordClients = list.New()
	s.publishRecordChannel = make(chan *Client, runtime.NumCPU())

	s.publishCacheClients = list.New()
	s.publishCacheChannel = make(chan *Client, runtime.NumCPU())

	go func() {
		for {
			ele := s.pullingRecordClients.Front()
			if ele == nil {
				time.Sleep(time.Millisecond * 100)
				continue
			}
			s.pullingRecordClients.Remove(ele)

			cli := ele.Value.(*Client)
			if cli.isClose {
				continue
			}

			s.pullingRecordChannel <- cli
		}
	}()

	for i := 0; i < pullerParallel; i++ {
		go func() {
			for {
				cli := <-s.pullingRecordChannel
				ctx := cli.conn.GetContext().(*ConnContext)
				req := &proto_room.RecordReq{
					TenantCode:    ctx.TenantCode,
					RoomCode:      ctx.RoomCode,
					LastMessageId: cli.lastMessageId,
					Count:         PullNumberOfMessagesPerTime,
				}

				serviceClient := config.RoomServiceConfig.ServiceClient()
				var opss client.CallOption = func(o *client.CallOptions) {
					o.RequestTimeout = time.Second * 10
					o.DialTimeout = time.Second * 10
				}
				rsp, err := serviceClient.Record(context.TODO(), req, opss)
				if err != nil {
					glog.Errorln(err)
					continue
				}
				for _, message := range rsp.Record {
					cli.SetLastMessageId(message.ServerMsgId)
					cli.PutRecord(message)
				}
				if len(rsp.Record) < PullNumberOfMessagesPerTime {
					cli.isPullingRecordDone = true
					if timestamp, num, err := MessageIdToInt(cli.lastMessageId); err != nil {
						cli.lastMessageIdTimestamp, cli.lastMessageIdNum = 0, 0
					} else {
						cli.lastMessageIdTimestamp, cli.lastMessageIdNum = timestamp, num
					}
					continue
				}

				s.pullingRecordClients.PushBack(cli)
			}
		}()
	}

	// 定时把客户端投放到，历史消息推送器
	// 把客户端的需要的历史消息推送给客户端
	go func() {
		ticker := time.NewTicker(time.Millisecond * 250)
		for {
			<-ticker.C
			var delEle []*list.Element
			for ele := s.publishRecordClients.Front(); ele != nil; ele = ele.Next() {
				cli := ele.Value.(*Client)
				if cli.isClose {
					delEle = append(delEle, ele)
					continue
				}
				if cli.isPublishRecordDone {
					delEle = append(delEle, ele)
					continue
				}
				s.publishRecordChannel <- cli
			}
			for _, ele := range delEle {
				ClientPullingRecordGauge.Dec()
				s.publishRecordClients.Remove(ele)
			}
		}
	}()

	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			for {
				cli := <-s.publishRecordChannel

				// 读取历史消息，并投递到Sender进行推送
				for i := 0; i < PublishNumberOfMessagePreTime; i++ {
					msg := cli.GetRecord()
					if msg == nil {
						break
					} else {
						Sender.Point(cli, msg)
					}
				}

				// 拉取历史消息已经完成,推送历史消息已经完成,标记同步历史消息已经完成
				if cli.record.Len() == 0 && cli.isPullingRecordDone {
					cli.isPublishRecordDone = true
				}

				// 客户端投递到 line 切换区
				if cli.isPublishRecordDone {
					ClientPublishCacheGauge.Inc()
					s.publishCacheClients.PushBack(cli)
				}
			}
		}()
	}

	go func() {
		ticker := time.NewTicker(time.Millisecond * 200)
		for {
			<-ticker.C
			var delEle []*list.Element
			for ele := s.publishCacheClients.Front(); ele != nil; ele = ele.Next() {
				cli := ele.Value.(*Client)
				if cli.isClose {
					delEle = append(delEle, ele)
					continue
				}
				if cli.isPublishCacheDone {
					delEle = append(delEle, ele)
					continue
				}
				s.publishCacheChannel <- cli
			}
			for _, ele := range delEle {
				s.publishCacheClients.Remove(ele)
			}
		}
	}()

	// publish cache
	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			for {
				cli := <-s.publishCacheChannel

				for i := 0; i < PublishNumberOfMessagePreTime; i++ {
					msg := cli.GetCache()
					if msg == nil {
						break
					}

					if msg.Type != "message.broadcast" && msg.Type != "message.cancel" {
						Sender.Point(cli, msg)
						continue
					}

					if cli.IsNewMessage(msg) {
						Sender.Point(cli, msg)
						continue
					}
				}
				if cli.cache.Len() == 0 {
					cli.isPublishCacheDone = true
				}

				if cli.isPublishCacheDone {
					Agent.ClientToDirectMode(cli)
				}
			}
		}()
	}
}
