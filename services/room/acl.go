package room

import (
	"fmt"
	proto_room "wchatv1/proto/room"
)

func Acl(message *proto_room.Message) (ctx interface{}, err error) {
	passToken := message.From.PassToken

	switch passToken.UserType {
	case "user":
		return AclUser(message)
	case "manager":
		return AclManager(message)
	case "system":
		break
	default:
		return nil, fmt.Errorf("未知的用户类型")
	}

	return
}

func AclUser(message *proto_room.Message) (interface{}, error) {
	passToken := message.From.PassToken

	switch message.Type {
	case "message.broadcast":
		roomAcl, err := GetRoomAcl(passToken.TenantCode, passToken.RoomCode, message.Type)
		if err != nil {
			return nil, err
		}
		userAcl, err := GetUserAcl(passToken.TenantCode, passToken.RoomCode, passToken.UserId, message.Type)
		if err != nil {
			return nil, err
		}

		tenantUserAcl, err := GetTenantUserAcl(passToken.TenantCode, passToken.UserId, message.Type)
		if err != nil {
			return nil, err
		}

		if roomAcl != "" && roomAcl != "1" {
			return nil, fmt.Errorf("房间禁止发言")
		}
		if userAcl == "0" {
			return nil, fmt.Errorf("您没有此房间得发言权限")
		}
		if tenantUserAcl == "0" {
			return nil, fmt.Errorf("您已经被禁止发言")
		}
		break
	case "message.cancel":
		msg, err := GetMessageById(passToken.TenantCode, passToken.RoomCode, message.ServerMsgId)
		if err != nil {
			return nil, err
		}
		if msg == nil {
			return nil, fmt.Errorf("消息ID不存在")
		}
		if msg.Type != "message.broadcast" && msg.Type != "message.p2p" {
			return nil, fmt.Errorf("操作错误")
		}
		if msg.From.PassToken.UserId != passToken.UserId {
			return nil, fmt.Errorf("消息ID错误，或者权限不足")
		}
		//if (msg.ServerMsgTimestamp + (time.Second * 30)) > uint64(time.Now().UnixNano()/1e6) {
		//
		//}
		return msg, nil
	default:
		return nil, fmt.Errorf("ACL规则不存在")
	}

	return nil, nil
}

func AclManager(message *proto_room.Message) (interface{}, error) {
	passToken := message.From.PassToken

	switch message.Type {
	case "message.broadcast":
		break
	case "message.cancel":
		msg, err := GetMessageById(passToken.TenantCode, passToken.RoomCode, message.ServerMsgId)
		if err != nil {
			return nil, err
		}
		if msg == nil {
			return nil, fmt.Errorf("消息ID不存在")
		}
		if msg.Type != "message.broadcast" && msg.Type != "message.p2p" {
			return nil, fmt.Errorf("消息类型错误")
		}
		return msg, nil
	case "acl.message.broadcast":
		break
	case "remove.user":

		break
	default:
		return nil, fmt.Errorf("ACL规则不存在")
	}

	return nil, nil
}
