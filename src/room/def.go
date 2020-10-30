package room

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	// 通知客户端，消息状态
	TypeMessageStateOk  = "message.broadcast.state.ok"
	TypeMessageStateNok = "message.broadcast.state.nok"
	// 消息取消
	TypeMessageCancel = "message.cancel"
	// 消息点赞
	TypeMessageLike = "message.like"

	// 通知
	TypeNotification = "notification"
	// 数据
	TypeData = "data.room.users.count"

	ContentTypeText  = "text"
	ContextTypeImage = "image"
	ContextTypeUrl   = "url"
)

// Code = $tenantCode:%roomCode
type Code string

//type Client struct {
//	Conn     frontier.Conn
//	RoomCode Code
//}

func ClientIndex(ctx *ConnContext) string {
	return fmt.Sprintf("%s.%d", ctx.TenantCode, ctx.UserId)
}
func RoomCode(ctx *ConnContext) Code {
	return Code(fmt.Sprintf("%s.%s", ctx.TenantCode, ctx.RoomCode))
}

func MessageIdToInt(messageId string) (timestamp int, num int, err error) {
	arr := strings.Split(messageId, "-")
	timestamp, err = strconv.Atoi(arr[0])
	if err != nil {
		return
	}
	num, err = strconv.Atoi(arr[1])
	if err != nil {
		return
	}

	return
}


