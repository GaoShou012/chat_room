package room

import (
	"container/list"
	"github.com/golang/glog"
	"os"
	"runtime"
	"sync"
	proto_room "wchatv1/proto/room"
)

var Sender sender

type pointJob struct {
	client  *Client
	message *proto_room.Message
}
type broadcastJob struct {
	wg           *sync.WaitGroup
	clients      *list.List
	clientsIndex *sync.Map
	message      *proto_room.Message
}

type sender struct {
	point struct {
		parallel int
		pool     sync.Pool
		jobs     []chan *pointJob
	}

	broadcast struct {
		pool sync.Pool
		jobs chan *broadcastJob
	}
}

func (s *sender) Init() {
	s.initBroadcast()
	s.initPoint(100000)
}

func (s *sender) initBroadcast() {
	parallel := runtime.NumCPU()
	jobsSize := parallel * 10

	s.broadcast.jobs = make(chan *broadcastJob, jobsSize)
	for i := 0; i < parallel; i++ {
		go func() {
			for {
				job := <-s.broadcast.jobs
				clients, clientsIndex, message, wg := job.clients, job.clientsIndex, job.message, job.wg

				if message == nil {
					glog.Errorln("message is nil", message)
					os.Exit(1)
					continue
				}
				if message.From == nil {
					glog.Errorln("message.From is nil", message)
					os.Exit(1)
					continue
				}

				// 过滤发送者
				ruledOut := make(map[*Client]bool)
				ruledOutSender(clientsIndex, message.From.Path, ruledOut)

				// 压缩数据
				msg, err := Codec.Encode(message)
				if err != nil {
					glog.Errorln(err)
					wg.Done()
					continue
				}

				// 遍历客户端，发送数据
				for ele := clients.Front(); ele != nil; ele = ele.Next() {
					cli := ele.Value.(*Client)
					if _, ok := ruledOut[cli]; ok {
						continue
					}
					if cli.isCache {
						cli.PutCache(message)
						continue
					}

					// 如果同步完成，消息广播，直接发送
					// 如果同步未完成，检查消息类型，非持久化的消息，直接发送

					// 如果是持久化的消息，检查消息ID，是否比客户端的要新，如果是就发送消息，并且同步完成
					// 如果不是此消息可以忽略
					if cli.isSyncDone {
						cli.conn.Sender(msg)
						//SenderParallelConnId.Send(cli.conn, msg)
					} else {
						if message.Type != "message.broadcast" && message.Type != "message.p2p" {
							cli.conn.Sender(msg)
							//SenderParallelConnId.Send(cli.conn, msg)
							continue
						}
						if cli.IsNewMessage(message) {
							cli.isSyncDone = true
							cli.conn.Sender(msg)
							//SenderParallelConnId.Send(cli.conn, msg)
							continue
						}
					}
				}
				wg.Done()
			}
		}()
	}
}
func (s *sender) initPoint(jobsSize int) {
	parallel := runtime.NumCPU()

	s.point.parallel = parallel
	s.point.pool.New = func() interface{} {
		return new(pointJob)
	}
	s.point.jobs = make([]chan *pointJob, parallel)
	for i := 0; i < parallel; i++ {
		s.point.jobs[i] = make(chan *pointJob, jobsSize)
	}

	for i := 0; i < parallel; i++ {
		jobs := s.point.jobs[i]
		go func(jobs chan *pointJob) {
			for {
				job := <-jobs
				cli, message := job.client, job.message
				s.point.pool.Put(job)
				if cli.isClose {
					continue
				}

				data, err := Codec.Encode(message)
				if err != nil {
					glog.Errorln(err)
					continue
				}

				cli.conn.Sender(data)
				//SenderParallelConnId.Send(cli.conn, data)
			}
		}(jobs)
	}
}

func (s *sender) Point(client *Client, message *proto_room.Message) {
	j := s.point.pool.Get().(*pointJob)
	j.client, j.message = client, message
	s.point.jobs[client.conn.GetId()%s.point.parallel] <- j
}

func (s *sender) Broadcast(clients *list.List, clientsIndex *sync.Map, message *proto_room.Message) *sync.WaitGroup {
	job := &broadcastJob{
		wg:           &sync.WaitGroup{},
		clients:      clients,
		clientsIndex: clientsIndex,
		message:      message,
	}
	job.wg.Add(1)
	s.broadcast.jobs <- job
	return job.wg
}
