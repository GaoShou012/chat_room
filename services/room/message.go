package room

import (
	"fmt"
	proto_room "wchatv1/proto/room"
)

func MsgFormatStateOk(msgType string, cliMsgId uint64) *proto_room.Message {
	msg := &proto_room.Message{
		Type:               fmt.Sprintf("%s.state.ok", msgType),
		Content:            "",
		ContentType:        "",
		ClientMsgId:        cliMsgId,
		ServerMsgId:        "",
		ServerMsgTimestamp: 0,
		From:               nil,
		To:                 nil,
		RuledOut:           nil,
	}
	return msg
}
func MsgFormatStateNok(msgType string, cliMsgId uint64, content string) *proto_room.Message {
	msg := &proto_room.Message{
		Type:               fmt.Sprintf("%s.state.nok", msgType),
		Content:            content,
		ContentType:        "",
		ClientMsgId:        cliMsgId,
		ServerMsgId:        "",
		ServerMsgTimestamp: 0,
		From:               nil,
		To:                 nil,
		RuledOut:           nil,
	}
	return msg
}

func MsgFormatMessageCancel(msgId string, passToken *proto_room.PassToken) *proto_room.Message {
	msg := &proto_room.Message{
		Type:        "message.cancel",
		Content:     fmt.Sprintf("%s 取消了一条消息", passToken.UserName),
		ContentType: "text",
		ClientMsgId: 0,
		ServerMsgId: msgId,
		From: &proto_room.From{
			Path:      "",
			PassToken: passToken,
		},
		To:       &proto_room.To{Path: fmt.Sprintf("%s.%s", passToken.TenantCode, passToken.RoomCode)},
		RuledOut: nil,
	}
	return msg
}

func NotificationUserJoin(tenantCode string, roomCode string, passToken *proto_room.PassToken) *proto_room.Message {
	msg := &proto_room.Message{
		Type:        "notification.user.join",
		Content:     "",
		ContentType: "",
		ClientMsgId: 0,
		ServerMsgId: "",
		From: &proto_room.From{
			Path:      "system",
			PassToken: passToken,
		},
		To:       &proto_room.To{Path: fmt.Sprintf("%s.%s", tenantCode, roomCode)},
		RuledOut: nil,
	}
	return msg
}

func NotificationUserLeave(tenantCode string, roomCode string, passToken *proto_room.PassToken) *proto_room.Message {
	msg := &proto_room.Message{
		Type:        "notification.user.leave",
		Content:     "",
		ContentType: "",
		ClientMsgId: 0,
		ServerMsgId: "",
		From: &proto_room.From{
			Path:      "system",
			PassToken: passToken,
		},
		To:       &proto_room.To{Path: fmt.Sprintf("%s.%s", tenantCode, roomCode)},
		RuledOut: nil,
	}
	return msg
}

func NotificationUsersCount(tenantCode string, roomCode string, count int64) *proto_room.Message {
	msg := &proto_room.Message{
		Type:        "notification.users.count",
		Content:     fmt.Sprintf("%d", count),
		ContentType: "text",
		ClientMsgId: 0,
		ServerMsgId: "",
		From: &proto_room.From{
			Path:      "system",
			PassToken: nil,
		},
		To:       &proto_room.To{Path: fmt.Sprintf("%s.%s", tenantCode, roomCode)},
		RuledOut: nil,
	}
	return msg
}

func NotificationRemoveUser(tenantCode string, roomCode string, content string) *proto_room.Message {
	msg := &proto_room.Message{
		Type:               "notification.user.remove",
		Content:            content,
		ContentType:        "text",
		ClientMsgId:        0,
		ServerMsgId:        "",
		ServerMsgTimestamp: 0,
		From: &proto_room.From{
			Path:      "system",
			PassToken: nil,
		},
		To:       &proto_room.To{Path: fmt.Sprintf("%s.%s", tenantCode, roomCode)},
		RuledOut: nil,
	}
	return msg
}
