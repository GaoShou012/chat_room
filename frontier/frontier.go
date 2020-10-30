package frontier

import (
	"net"
	"time"

	"github.com/golang/glog"
	"github.com/mailru/easygo/netpoll"
)

// NewFrontier 创建长连接入层.
func NewFrontier(ln net.Listener, ca Accepter, ja Joinable, la Leavable, rt Router, opts ...Option) *Frontier {
	poller, err := netpoll.New(nil)
	if err != nil {
		panic(err)
	}
	desc := netpoll.Must(netpoll.HandleListener(ln, netpoll.EventRead|netpoll.EventOneShot))
	f := &Frontier{
		ln: ln,
		ca: ca,
		ja: ja,
		la: la,
		rt: rt,

		accept: make(chan error, 1),
		desc:   desc,
		poller: poller,
	}
	for _, opt := range opts {
		opt(f)
	}
	if f.runner == nil {
		f.runner = func(f func()) {
			go f()
		}
	}
	if f.channelBuilder == nil {
		f.channelBuilder = NewWSChannel
	}
	return f
}

// Option 代表 Frontier 初始化选项.
type Option func(*Frontier)

// WithRunner 设置异步任务执行器.
func WithRunner(fn func(func())) Option {
	return func(f *Frontier) {
		f.runner = fn
	}
}

// WithChannelBuilder 设置消息通道构建器.
func WithChannelBuilder(fn func(net.Conn) Channel) Option {
	return func(f *Frontier) {
		f.channelBuilder = fn
	}
}

// Frontier 代表长连接入层.
type Frontier struct {
	ln net.Listener
	ca Accepter
	ja Joinable
	la Leavable
	rt Router

	accept chan error
	desc   *netpoll.Desc
	poller netpoll.Poller

	runner         func(func())
	channelBuilder func(net.Conn) Channel
}

// Start 启动长连服务.
func (f *Frontier) Start() error {
	return f.poller.Start(f.desc, func(e netpoll.Event) {
		f.runner(func() {
			conn, err := f.ln.Accept()
			if err != nil {
				f.accept <- err
				return
			}
			glog.V(3).Infof("accpect: %s", conn.RemoteAddr().String())
			f.accept <- nil
			f.handle(conn)
		})
		err := <-f.accept
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				glog.Errorf("accept error: %v; retrying in %v", err, 5*time.Millisecond)
				time.Sleep(5 * time.Millisecond)
			} else {
				glog.Errorf("accept error: %v", err)
			}
		}
		f.poller.Resume(f.desc)
	})
}

// Stop 停止长连服务.
func (f *Frontier) Stop() error {
	if err := f.poller.Stop(f.desc); err != nil {
		return err
	}
	return f.ln.Close()
}

func (f *Frontier) handle(conn net.Conn) {
	channel := f.channelBuilder(conn)
	if err := channel.Accept(f.ca); err != nil {
		glog.Errorf("%s: channel accept error: %v", channel, err)
		if err := channel.Close(); err != nil {
			glog.Errorf("%s: close channel failed, %v", channel, err)
		}
		return
	}
	if err := channel.Context().JoinTo(channel, f.ja); err != nil {
		glog.Errorf("%s: channel join error: %v", channel, err)

		f.la.Leave(channel)
		if err := channel.Close(); err != nil {
			glog.Errorf("%s: close channel failed, %v", channel, err)
		}
		return
	}

	desc := netpoll.Must(netpoll.HandleRead(conn))
	f.poller.Start(desc, func(ev netpoll.Event) {
		if ev&(netpoll.EventReadHup|netpoll.EventHup) != 0 {
			glog.V(3).Infof("%s: receive: %v; close connection", channel, ev)

			f.poller.Stop(desc)
			f.runner(func() {
				f.la.Leave(channel)
				if err := channel.Close(); err != nil {
					glog.Errorf("%s: close channel failed, %v", channel, err)
				}
			})
			return
		}
		// Here we can read some new message from connection.
		// We can not read it right here in callback, because then we will
		// block the poller's inner loop.
		// We do not want to spawn a new goroutine to read single message.
		// But we want to reuse previously spawned goroutine.
		f.runner(func() {
			m, err := channel.Receive()
			if err != nil {
				glog.V(3).Infof("%s: receive: %v; close connection", channel, err)

				f.poller.Stop(desc)
				f.la.Leave(channel)
				if err := channel.Close(); err != nil {
					glog.Errorf("%s: close channel failed, %v", channel, err)
				}
				return
			}
			f.rt.Route(m)
		})
	})
}
