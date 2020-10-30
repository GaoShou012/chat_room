package frontier

import (
	"sync"
	"time"

	"github.com/golang/glog"
)

// Sender 定义消息发送器.
type Sender interface {
	// SendTo 发送消息到指定通道.
	SendTo(Sendable, Message)
}

// NewSender 初始化带缓冲消息发送器.
func NewSender(capcity int) Sender {
	s := &sender{jobs: make(chan *job, capcity)}
	s.pool.New = func() interface{} {
		return new(job)
	}
	go s.do()
	return s
}

// job 定义消息发送请求.
type job struct {
	ch  Sendable
	msg Message
}

// sender 定义消息发送器.
type sender struct {
	jobs chan *job
	pool sync.Pool
}

// SendTo 发送消息到指定通道.
func (s *sender) SendTo(ch Sendable, m Message) {
	j := s.pool.Get().(*job)
	j.ch, j.msg = ch, m
	s.jobs <- j
}

func (s *sender) do() {
	for job := range s.jobs {
		m := job.msg
		start := time.Now()
		if err := job.ch.Send(m); err != nil {
			failedSendCounter.WithLabelValues(err.Error()).Inc()
			glog.V(1).Infof("send message(router: %v, receiver: %d) to %s failed, %v", m.GetRId(), m.GetRType(), job.ch, err)
		}
		sendDurationsHistogram.Observe(time.Since(start).Seconds())
		s.pool.Put(job)
	}
}

// NewBalanceSender 创建平衡消息发送器.
//
// 依据 Channel 哈希值进行消息发送平衡.
//
// Channel 与 Sender 有 sticky 特性.
func NewBalanceSender(senders ...Sender) Sender {
	return &bsender{senders: senders, mask: len(senders)}
}

type bsender struct {
	senders []Sender
	mask    int
}

// SendTo 发送消息到指定通道.
func (b *bsender) SendTo(ch Sendable, m Message) {
	b.senders[ch.Hash()%b.mask].SendTo(ch, m)
}
