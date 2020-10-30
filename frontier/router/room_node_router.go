package router

import (
	"context"
	"io"
	"sync"
	"time"

	"github.com/golang/glog"

	"wchat.im/frontier"
	mpb "wchat.im/proto/message"
	pb "wchat.im/proto/room/node"
)

// NewRoomNodeRouter 创建与房间节点的通信.
func NewRoomNodeRouter(client pb.RoomNodeServiceClient, routeTo frontier.Router) *RoomNodeRouter {
	rn := &RoomNodeRouter{
		client:   client,
		reports:  make(chan *pb.RoomReport, 10000),
		messages: make(chan *mpb.Message, 10000),
	}
	go rn.serve()
	return rn
}

// RoomNodeRouter 代表与房间节点的连接.
type RoomNodeRouter struct {
	client   pb.RoomNodeServiceClient
	reports  chan *pb.RoomReport
	messages chan *mpb.Message

	routeTo frontier.Router
}

var _ frontier.Router = new(RoomNodeRouter)

func (n *RoomNodeRouter) serve() {
	for m := range n.messages {
		glog.V(4).Infof("receive message from room node: %v", m)

		// TODO: 处理房间命令.

		if err := n.routeTo.Route(m); err != nil {
			glog.Errorf("route room node message error: %v", err)
		}
	}
}

func (n *RoomNodeRouter) Route(m frontier.Message) error {
	// TODO: 支持向上路由.
	return nil
}

// Report 上报报告到 RoomNode.
func (n *RoomNodeRouter) Report(r *pb.RoomReport) {
	select {
	case n.reports <- r:
		return
	default:
	}

	for {
		t := time.NewTimer(2 * time.Second)
		select {
		case n.reports <- r:
			if !t.Stop() {
				<-t.C
			}
			return
		case <-t.C:
			select {
			case e := <-n.reports:
				// TODO(sven): 对 r 进行脱敏输出.
				glog.Errorf("failed to report to room node, %+v", r)
				r = e
			default:
			}
		}
	}
}

func (n *RoomNodeRouter) Serve(ctx context.Context) {
	var r *pb.RoomReport
	for {
		stream, err := n.client.Streaming(ctx)
		if err != nil {
			glog.Errorf("open Streaming to room node failed, %v", err)
			time.Sleep(1 * time.Second)
			continue
		}
		glog.V(1).Info("connected to Streaming of room node.")

		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()

			if r != nil {
				if err := stream.Send(r); err != nil {
					glog.Errorf("send report to room node failed, %v", err)
					return
				}
			}
			r = nil
			for {
				select {
				case <-stream.Context().Done():
					return
				case r = <-n.reports:
					if err := stream.Send(r); err != nil {
						glog.Errorf("send report to room node failed, %v", err)
						return
					}
					r = nil
				}
			}
		}()

		go func() {
			defer wg.Done()

		RECV:
			for {
				c, err := stream.Recv()
				if err == io.EOF {
					return
				}
				if err != nil {
					glog.Errorf("receive from room node failed, %v", err)
					return
				}
				select {
				case n.messages <- c:
				default:
					for {
						t := time.NewTimer(1 * time.Second)
						select {
						case n.messages <- c:
							if !t.Stop() {
								<-t.C
							}
							continue RECV
						case <-t.C:
							// TODO(sven): 对 c 进行脱敏输出.
							glog.Errorf("send room node command to chan buffer timeout, %v", c)
						}
					}
				}
			}
		}()
		wg.Wait()
	}
}
