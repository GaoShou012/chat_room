package frontier

import (
	"container/list"
	"fmt"
	"hash/fnv"
	"io"
	"sync"

	"github.com/golang/glog"
	pb "wchat.im/proto/message"
)

// Context 定义消息上下文.
type Context interface {
	// WriteTo 序列化消息.
	WriteTo(Message, io.Writer) error
	// ReadFrom 反序列化消息.
	ReadFrom(io.Reader) (Message, error)
	// JoinTo 加入 Channel 到 Joinable.
	JoinTo(Channel, Joinable) error
}

// Joinable 定义 Channel 加入接口.
type Joinable interface {
	// JoinRoom 定义加入房间.
	JoinRoom(id string, ch Channel)
}

// Leavable 定义离开能力.
type Leavable interface {
	// Leave 离开当前订阅.
	Leave(interface{})
}

// Accepter 定义连接接收器.
type Accepter interface {
	// TokenConnect 实现 Token 连接请求.
	//
	// 不允许连接时 Context 为 nil.
	TokenConnect(ch Channel, token string) (Context, error)
}

// Message 定义推送消息.
type Message interface {
	// 消息投递类型.
	GetRType() pb.RType
	// 消息投递标识.
	//
	// P2P  : 用户号
	// GROUP: 群组号
	// ROOM : 房间号
	GetRId() string
}

// 默认哈希函数.
var Hash = func(s string) int {
	h := fnv.New32()
	h.Write([]byte(s))
	return int(h.Sum32())
}

// Identifier 表示定位符.
type Identifier interface {
	// 用于消息发送 sticky.
	Hash() int
	// 唯一标识.
	String() string
}

type nameIdentifier string

func (n nameIdentifier) String() string {
	return string(n)
}

func (n nameIdentifier) Hash() int {
	return Hash(string(n))
}

var _ Identifier = new(nameIdentifier)

// Sendable 定义消息发送能力.
type Sendable interface {
	Identifier
	// Send 发送消息.
	Send(Message) error
}

// Receivable 定义消息接收能力.
type Receivable interface {
	// Receive 接收消息.
	Receive() (Message, error)
}

// Router 定义消息路由能力.
type Router interface {
	// 进行消息路由.
	Route(Message) error
}

// Anchorable 定义加入能力.
type Anchorable interface {
	// SetAnchor 设置锚点.
	SetAnchor(Leavable, interface{})
	// IsAnchored 判断锚点是否已设置.
	IsAnchored(Leavable) bool
	// DelAnchor 删除锚点.
	DelAnchor(Leavable) bool
	// DelAnchors 删除所有锚点.
	DelAnchors()
}

// Anchor 实现 Anchorable.
type Anchor struct {
	anchors map[Leavable]interface{}
}

var _ Anchorable = new(Anchor)

func (j *Anchor) SetAnchor(l Leavable, an interface{}) {
	if j.anchors == nil {
		j.anchors = make(map[Leavable]interface{})
	}
	j.anchors[l] = an
}

func (j *Anchor) IsAnchored(l Leavable) bool {
	_, exist := j.anchors[l]
	return exist
}

func (j *Anchor) DelAnchor(l Leavable) bool {
	an, exist := j.anchors[l]
	if !exist {
		return false
	}
	l.Leave(an)
	delete(j.anchors, l)
	return exist
}

func (j *Anchor) DelAnchors() {
	for l, an := range j.anchors {
		l.Leave(an)
	}
	j.anchors = nil
}

// Channel 代表消息通道.
type Channel interface {
	// Context 返回连接上下文.
	Context() Context
	// Accept 接受连接请求.
	Accept(Accepter) error
	// Close 关闭 Channel.
	Close() error

	Sendable
	Receivable
	Anchorable
}

// NewRoom 创建指定名称房间.
func NewRoom(name string, sender Sender) *Room {
	return &Room{
		nameIdentifier: nameIdentifier(name),
		sender:         sender,
		subs:           list.New(),
	}
}

// Room 定义一个房间.
type Room struct {
	nameIdentifier
	Anchor

	sender Sender
	mut    sync.RWMutex
	subs   *list.List
	// 房间是否废弃.
	discard bool
}

var _ Anchorable = new(Room)
var _ Leavable = new(Room)
var _ Router = new(Room)

// Join 加入房间.
//
// true: 成功, false: 失败.
func (r *Room) Join(ch Channel) bool {
	if ch.IsAnchored(r) {
		return true
	}
	r.mut.Lock()
	defer r.mut.Unlock()

	if r.discard {
		return false
	}
	ch.SetAnchor(r, r.subs.PushFront(ch))
	return true
}

// Leave 离开房间.
//
// 最后一个人离开后，房间废弃.
func (r *Room) Leave(an interface{}) {
	r.mut.Lock()
	r.subs.Remove(an.(*list.Element))
	r.discard = r.subs.Len() <= 0
	r.mut.Unlock()
	if r.discard {
		r.DelAnchors()
	}
}

// Route 路由消息给房间内人员.
func (r *Room) Route(m Message) error {
	r.mut.RLock()
	for e := r.subs.Front(); e != nil; e = e.Next() {
		r.sender.SendTo(e.Value.(Sendable), m)
	}
	r.mut.RUnlock()
	return nil
}

// NewBucket 创建消息通道容器.
func NewBucket(name string, sender Sender, capcity int) *Bucket {
	return &Bucket{
		nameIdentifier: nameIdentifier(name),
		sender:         sender,
		rooms:          make(map[string]*Room, capcity),
	}
}

// Bucket 定义消息通道容器.
type Bucket struct {
	nameIdentifier

	sender Sender
	mut    sync.RWMutex
	rooms  map[string]*Room
}

var _ Joinable = new(Bucket)
var _ Leavable = new(Bucket)
var _ Router = new(Bucket)

// JoinRoom 加入指定房间.
func (b *Bucket) JoinRoom(id string, ch Channel) {
	b.mut.RLock()
	r, ok := b.rooms[id]
	b.mut.RUnlock()

	if ok && r.Join(ch) {
		return
	}

	b.mut.Lock()
	defer b.mut.Unlock()

	for {
		v, ok := b.rooms[id]
		switch {
		case !ok || v == r:
			v = NewRoom(id, b.sender)
			v.SetAnchor(b, v)
			v.Join(ch)
			b.rooms[id] = v
			return
		default:
			if v.Join(ch) {
				return
			}
		}
	}
}

// LeaveRoom 离开指定房间.
func (b *Bucket) LeaveRoom(id string, ch Channel) {
	b.mut.RLock()
	r, ok := b.rooms[id]
	b.mut.RUnlock()
	if !ok {
		return
	}
	ch.DelAnchor(r)
}

// Leave 解除与 Bucket 绑定关系.
func (b *Bucket) Leave(an interface{}) {
	switch an.(type) {
	case *Room:
		id := an.(*Room).String()
		b.mut.Lock()
		if b.rooms[id] == an {
			delete(b.rooms, id)
		}
		b.mut.Unlock()
	}
}

// Route 实现消息路由.
func (b *Bucket) Route(m Message) error {
	switch m.GetRType() {
	case pb.RType_ROOM:
		b.mut.RLock()
		r, ok := b.rooms[m.GetRId()]
		b.mut.RUnlock()
		if !ok {
			glog.V(4).Infof("no channels in room: %s", m.GetRId())
			return nil
		}
		return r.Route(m)
	default:
		return fmt.Errorf("%s: unsupported route type, %v", b, m.GetRType())
	}
	return nil
}
