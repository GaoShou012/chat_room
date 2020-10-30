package room

import (
	"container/list"
	"github.com/GaoShou012/frontier"
	"time"
	proto_room "wchatv1/proto/room"
)

type Client struct {
	state int
	conn  frontier.Conn

	isClose       bool
	isSyncDone    bool
	isCache       bool
	cache         *list.List
	record        *list.List
	lastMessageId string

	isPullingRecordDone bool
	isPublishRecordDone bool
	isPublishCacheDone  bool

	lastMessageIdTimestamp int
	lastMessageIdNum       int

	startPublishCacheTime time.Time
}

func (c *Client) IsClose() bool {
	return c.isClose
}

// 弹出一条历史消息
func (c *Client) GetRecord() *proto_room.Message {
	ele := c.record.Front()
	if ele == nil {
		return nil
	}
	c.record.Remove(ele)
	return ele.Value.(*proto_room.Message)
}

// 放入一条历史消息
func (c *Client) PutRecord(message *proto_room.Message) {
	c.record.PushBack(message)
}

// 设置最后拉取的历史消息ID
func (c *Client) SetLastMessageId(id string) {
	c.lastMessageId = id
}

// 读取最后拉取历的史消息ID
func (c *Client) GetLastMessageId() string {
	return c.lastMessageId
}

func (c *Client) PopCache() (*proto_room.Message, error) {
	ele := c.cache.Front()
	if ele == nil {
		return nil, nil
	}
	return ele.Value.(*proto_room.Message), nil
}

// 客户端 历史流水线 切换到 缓存流水线
func (c *Client) RecordLineToCacheLine() (bool, error) {
	timestamp, num, err := MessageIdToInt(c.lastMessageId)
	if err != nil {
		return false, nil
	}

	for {
		ele := c.cache.Front()
		if ele == nil {
			return true, nil
		}

		msg := ele.Value.(*proto_room.Message)
		t, n, err := MessageIdToInt(msg.ServerMsgId)
		if err != nil {
			c.cache.Remove(ele)
			continue
		}

		if t <= timestamp && n <= num {
			c.cache.Remove(ele)
		}

		if t >= timestamp && n >= num {
			return true, nil
		}
	}
}

func (c *Client) GetCache() *proto_room.Message {
	ele := c.cache.Front()
	if ele == nil {
		return nil
	}
	c.cache.Remove(ele)
	return ele.Value.(*proto_room.Message)
}
func (c *Client) PutCache(message *proto_room.Message) {
	c.cache.PushBack(message)
}

func (c *Client) IsNewMessage(message *proto_room.Message) bool {
	t, n, err := MessageIdToInt(message.ServerMsgId)
	if err != nil {
		return false
	}

	if t < c.lastMessageIdTimestamp {
		return false
	}

	if t == c.lastMessageIdTimestamp && n <= c.lastMessageIdNum {
		return false
	}

	return true
}

func (c *Client) CheckMessageContinuity(message *proto_room.Message) bool {
	if message.Type != "message.broadcast" && message.Type != "message.p2p" {
		return true
	}

	t, n, err := MessageIdToInt(message.ServerMsgId)
	if err != nil {
		return false
	}

	if t < c.lastMessageIdTimestamp {
		return false
	}
	if t == c.lastMessageIdTimestamp && n <= c.lastMessageIdNum {
		return false
	}

	return true
}

func (c *Client) WhetherCache(message *proto_room.Message) bool {
	if c.lastMessageId == "0" {
		return true
	}
	if message.Type != "message.broadcast" && message.Type != "message.p2p" {
		return true
	}
	timestamp, num, err := MessageIdToInt(c.lastMessageId)
	if err != nil {
		return false
	}
	t, n, err := MessageIdToInt(message.ServerMsgId)
	if err != nil {
		return false
	}
	if t >= timestamp && n >= num {
		return true
	}

	return false
}
