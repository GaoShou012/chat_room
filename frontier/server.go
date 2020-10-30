package frontier

import (
	"github.com/golang/glog"
)

// NewServer 创建 frontier 服务.
func NewServer(name string, buckets []*Bucket) *Server {
	return &Server{
		nameIdentifier: nameIdentifier(name),
		buckets:        buckets,
		bktmask:        len(buckets),
	}
}

// Server 定义 frontier Server.
type Server struct {
	nameIdentifier
	buckets []*Bucket
	bktmask int
}

var _ Joinable = new(Server)
var _ Leavable = new(Server)
var _ Router = new(Server)

// JoinRoom 加入指定房间.
func (s *Server) JoinRoom(id string, ch Channel) {
	idx := Hash(id) % s.bktmask
	s.buckets[idx].JoinRoom(id, ch)
}

// Leave 断开长连服务关联.
func (s *Server) Leave(ch interface{}) {
	ch.(Anchorable).DelAnchors()
}

// Route 实现消息路由功能.
func (s *Server) Route(m Message) error {
	idx := Hash(m.GetRId()) % s.bktmask
	err := s.buckets[idx].Route(m)
	if err != nil {
		glog.Errorf("%s: route message error: %v", s, err)

		failedSendCounter.WithLabelValues(err.Error()).Inc()
	}
	return err
}
