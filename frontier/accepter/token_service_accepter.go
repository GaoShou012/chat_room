package accepter

import (
	"context"
	"io"
	"io/ioutil"
	"time"

	"github.com/golang/glog"

	"wchat.im/frontier"
	mpb "wchat.im/proto/message"
	tpb "wchat.im/proto/token"
)

func NewTokenServiceAccepter(ts tpb.TokenServiceClient) frontier.Accepter {
	return &tsAccepter{ts: ts}

}

type tsAccepter struct {
	ts tpb.TokenServiceClient
}

// TokenConnect 实现 Token 连接请求.
//
// 不允许连接时 Context 为 nil.
func (t *tsAccepter) TokenConnect(ch frontier.Channel, token string) (frontier.Context, error) {
	glog.V(4).Infof("token: %s", token)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := t.ts.CheckToken(ctx, &tpb.CheckTokenRequest{UserToken: token})
	if err != nil {
		glog.Errorf("validate token failed, %v", err)
		return nil, err
	}

	return &tkctx{
		userId:   resp.UserId,
		userName: resp.UserName,
		roomId:   resp.RoomId,
	}, nil
}

type tkctx struct {
	userId   string
	userName string
	roomId   string
}

func (c *tkctx) JoinTo(ch frontier.Channel, j frontier.Joinable) error {
	j.JoinRoom(c.roomId, ch)
	return nil
}

// 序列化消息.
func (c *tkctx) WriteTo(m frontier.Message, w io.Writer) (err error) {
	// TODO: 支持各种序列化协议
	mp := m.(*mpb.Message)
	_, err = w.Write(mp.Body)
	return
}

// 反序列化消息.
func (c *tkctx) ReadFrom(r io.Reader) (frontier.Message, error) {
	// TODO: 支持各种序列化协议

	m := &mpb.Message{}
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	m.Body = data
	m.RType = mpb.RType_ROOM
	m.RId = c.roomId
	return m, nil
}
