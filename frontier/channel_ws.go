package frontier

import (
	"bytes"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/golang/glog"
)

var (
	ReadMessageTimeout  = 100 * time.Millisecond
	WriteMessageTimeout = 100 * time.Millisecond
)

// NewWSChannel 创建 Websocket 长连推送通道.
func NewWSChannel(conn net.Conn) Channel {
	return &WSChannel{
		hash: Hash(conn.RemoteAddr().String()),
		conn: conn,
	}
}

// WSChannel 代表 websocket 长连接.
type WSChannel struct {
	Anchor
	hash int
	ctx  Context
	conn net.Conn
}

var _ Channel = new(WSChannel)

func (c *WSChannel) String() string {
	return c.conn.RemoteAddr().String()
}

// Hash 返回消息通道哈希值.
func (c *WSChannel) Hash() int {
	return c.hash
}

// Context 返回连接上下文.
func (c *WSChannel) Context() Context {
	return c.ctx
}

// Authenticate 实现连接认证.
func (c *WSChannel) Accept(a Accepter) error {
	var token string
	hs, err := ws.Upgrader{
		OnRequest: func(uri []byte) error {
			u, err := url.Parse(string(uri))
			if err != nil {
				return err
			}
			if t := u.Query().Get("token"); t != "" {
				token = t
			}
			return nil
		},
		OnHeader: func(key, value []byte) error {
			if token == "" && bytes.Equal(key, []byte("Authorization")) {
				if bytes.HasPrefix(value, []byte("Bearer ")) {
					t := bytes.TrimPrefix(value, []byte("Bearer "))
					if len(t) > 0 {
						token = string(t)
					}
				}
			}
			return nil
		},
		OnBeforeUpgrade: func() (_ ws.HandshakeHeader, err error) {
			c.ctx, err = a.TokenConnect(c, token)
			if err != nil {
				return nil, err
			}
			if c.ctx == nil {
				return nil, ws.RejectConnectionError(ws.RejectionStatus(http.StatusUnauthorized),
					ws.RejectionReason("authorize failed"))
			}
			return
		},
	}.Upgrade(c.conn)
	if err != nil {
		return err
	}
	glog.V(3).Infof("%s: established websocket connection: %+v", c, hs)
	return nil
}

func (c *WSChannel) Close() (err error) {
	defer func() {
		if e := c.conn.Close(); err == nil {
			err = e
		}
	}()
	if c.ctx == nil {
		return
	}
	w := wsutil.NewWriter(c.conn, ws.StateServerSide, ws.OpClose)
	if _, e := w.Write([]byte("close")); err == nil {
		err = e
	}
	if e := w.Flush(); err == nil {
		err = e
	}
	return
}

// Send 发送消息.
func (c *WSChannel) Send(m Message) error {
	glog.V(4).Infof("%s: send %s", c, m)

	if err := c.conn.SetWriteDeadline(time.Now().Add(WriteMessageTimeout)); err != nil {
		glog.Errorf("%s: send error: %v", c, err)
		return err
	}
	w := wsutil.NewWriter(c.conn, ws.StateServerSide, ws.OpBinary)
	if err := c.ctx.WriteTo(m, w); err != nil {
		glog.Errorf("%s: send error: %v", c, err)
		return err
	}
	return w.Flush()
}

// Receive 接收消息.
func (c *WSChannel) Receive() (Message, error) {
	if err := c.conn.SetReadDeadline(time.Now().Add(ReadMessageTimeout)); err != nil {
		return nil, err
	}
	h, r, err := wsutil.NextReader(c.conn, ws.StateServerSide)
	if err != nil {
		return nil, err
	}
	if h.OpCode.IsControl() {
		return nil, wsutil.ControlFrameHandler(c.conn, ws.StateServerSide)(h, r)
	}
	return c.ctx.ReadFrom(r)
}
