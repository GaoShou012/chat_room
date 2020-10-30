package router

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/glog"
	"wchat.im/frontier"
	pb "wchat.im/proto/message"
	rpb "wchat.im/proto/room"
)

func NewRoomServiceRouter(rs rpb.RoomServiceClient, timeout time.Duration) frontier.Router {
	return &rsRouter{rs: rs, timeout: timeout}
}

type rsRouter struct {
	rs      rpb.RoomServiceClient
	timeout time.Duration
}

func (r *rsRouter) Route(m frontier.Message) error {
	rm, ok := m.(*pb.Message)
	if !ok {
		return fmt.Errorf("message format not support by room service: %v", m)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if _, err := r.rs.SendMessage(ctx, rm); err != nil {
		glog.Errorf("route message to room service error: %v", rm)
		return err
	}
	return nil
}
