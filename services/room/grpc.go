package room

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/golang/glog"
	"strconv"
	"sync"
	proto_room "wchatv1/proto/room"
)

var _ proto_room.RoomServiceHandler = &Service{}

type Service struct {
	Key                 []byte
	BroadcastToFrontier chan *proto_room.Message

	mutex sync.Mutex
}

func (s *Service) SetVirtualUserCount(ctx context.Context, req *proto_room.SetVirtualUserCountReq, rsp *proto_room.SetVirtualUserCountRsp) error {
	err := SetVirtualUserCount(req.TenantCode, req.RoomCode, req.Count)
	if err != nil {
		rsp.Code = 1
		rsp.Desc = err.Error()
		return nil
	}
	return nil
}

func (s *Service) GetVirtualUserCount(ctx context.Context, req *proto_room.GetVirtualUserCountReq, rsp *proto_room.GetVirtualUserCountRsp) error {
	num, err := GetVirtualUserCount(req.TenantCode, req.RoomCode)
	if err != nil {
		rsp.Code = 1
		rsp.Desc = err.Error()
		return nil
	}
	rsp.Code = 0
	rsp.Count = num
	return nil
}

func (s *Service) SetTenantUserAcl(ctx context.Context, req *proto_room.SetTenantUserAclReq, rsp *proto_room.SetTenantUserAclRsp) error {
	err := SetTenantUserAcl(req.TenantCode, req.UserId, req.Key, req.Val)
	if err != nil {
		rsp.Code = 1
		rsp.Desc = err.Error()
		return nil
	}
	rsp.Code = 0
	rsp.Desc = ""
	return nil
}

func (s *Service) GetTenantUserAcl(ctx context.Context, req *proto_room.GetTenantUserAclReq, rsp *proto_room.GetTenantUserAclRsp) error {
	val, err := GetTenantUserAcl(req.TenantCode, req.UserId, req.Key)
	if err != nil {
		rsp.Code = 1
		rsp.Desc = err.Error()
		return nil
	}
	rsp.Code = 0
	rsp.Desc = "ok"
	rsp.Val = val
	fmt.Println(rsp)
	return nil
}

func (s *Service) FrontierPing(ctx context.Context, req *proto_room.FrontierPingReq, rsp *proto_room.FrontierPingRsp) error {
	FrontierPing(req.FrontierId)
	return nil
}

func (s *Service) MakePassToken(ctx context.Context, req *proto_room.MakePassTokenReq, rsp *proto_room.MakePassTokenRsp) error {
	passToken := &proto_room.PassToken{
		TenantCode: req.TenantCode,
		RoomCode:   req.RoomCode,
		UserType:   req.UserType,
		UserId:     req.UserId,
		UserName:   req.UserName,
		UserThumb:  req.UserThumb,
		UserTags:   req.UserTags,
	}

	str, err := Codec.EncodePassToken(s.Key, passToken)
	if err != nil {
		return err
	}

	rsp.Token = str
	return nil
}

func (s *Service) ViewPassToken(ctx context.Context, req *proto_room.ViewPassTokenReq, rsp *proto_room.ViewPassTokenRsp) error {
	passToken, err := Codec.DecodePassToken(s.Key, req.Token)
	if err != nil {
		return err
	}
	rsp.PassToken = passToken
	return nil
}

func (s *Service) Join(ctx context.Context, req *proto_room.JoinReq, rsp *proto_room.JoinRsp) error {
	passToken := req.PassToken

	// 增加在线人数数量
	// 发出加入通知
	frontierId := req.FrontierId
	tenantCode := passToken.TenantCode
	roomCode := passToken.RoomCode
	_, err := IncrUsersCount(tenantCode, roomCode, frontierId)
	if err != nil {
		return err
	}
	num, err := GetUsersCount(tenantCode, roomCode)
	if err != nil {
		return err
	}

	// 获取虚拟用户数量
	{
		tmp, err := GetVirtualUserCount(tenantCode, roomCode)
		if err != nil {
			glog.Errorln(err)
		}
		num += int64(tmp)
	}

	//s.BroadcastToFrontier <- NotificationUserJoin(tenantCode, roomCode, passToken)
	s.BroadcastToFrontier <- NotificationUsersCount(tenantCode, roomCode, num)
	return nil
}

func (s *Service) Leave(ctx context.Context, req *proto_room.LeaveReq, rsp *proto_room.LeaveRsp) error {
	frontierId := req.FrontierId
	passToken := req.PassToken
	tenantCode := passToken.TenantCode
	roomCode := passToken.RoomCode
	_, err := DecrUsersCount(tenantCode, roomCode, frontierId)
	if err != nil {
		return err
	}
	num, err := GetUsersCount(tenantCode, roomCode)
	if err != nil {
		return err
	}

	// 获取虚拟用户数量
	{
		tmp, err := GetVirtualUserCount(tenantCode, roomCode)
		if err != nil {
			glog.Errorln(err)
		}
		num += int64(tmp)
	}

	//s.BroadcastToFrontier <- NotificationUserLeave(tenantCode, roomCode, passToken)
	s.BroadcastToFrontier <- NotificationUsersCount(tenantCode, roomCode, num)
	return nil
}

func (s *Service) GetUsersCount(ctx context.Context, req *proto_room.GetUsersCountReq, rsp *proto_room.GetUsersCountRsp) error {
	count, err := GetRoomUsersCount(req.TenantCode, req.RoomCode)
	if err != nil {
		return err
	}

	rsp.Count = count
	return nil
}

func (s *Service) SetUsersCount(ctx context.Context, req *proto_room.SetUsersCountReq, rsp *proto_room.SetUsersCountRsp) error {
	err := SetRoomUsersCount(req.FrontierId, req.TenantCode, req.RoomCode, req.Count)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) Broadcast(ctx context.Context, req *proto_room.BroadcastReq, rsp *proto_room.BroadcastRsp) error {
	if req.Message.From == nil {
		return fmt.Errorf("message.from is nil")
	}

	message := req.Message
	passToken := message.From.PassToken

	aclCtx, err := Acl(message)
	if err != nil {
		rsp.Message = MsgFormatStateNok(req.Message.Type, req.Message.ClientMsgId, err.Error())
		return nil
	}

	switch req.Message.Type {
	case "message.broadcast":
		message.To = &proto_room.To{Path: fmt.Sprintf("%s.%s", passToken.TenantCode, passToken.RoomCode)}
		mid, err := Stream(passToken.TenantCode, passToken.RoomCode, message)
		if err != nil {
			return err
		}
		message.ServerMsgId = mid
		s.BroadcastToFrontier <- message

		rsp.Message = MsgFormatStateOk(req.Message.Type, req.Message.ClientMsgId)
		rsp.Message.ServerMsgId = mid
		break
	case "message.cancel":
		targetMsg, ok := aclCtx.(*proto_room.Message)
		if !ok {
			return fmt.Errorf("aclCtx Assert payload is failed")
		}
		DelMessageById(passToken.TenantCode, passToken.RoomCode, targetMsg.ServerMsgId)

		msgId := targetMsg.ServerMsgId
		msg := MsgFormatMessageCancel(msgId, passToken)
		Stream(passToken.TenantCode, passToken.RoomCode, msg)
		s.BroadcastToFrontier <- msg

		rsp.Message = MsgFormatStateOk(req.Message.Type, req.Message.ClientMsgId)
		rsp.Message.ServerMsgId = msgId
		break
	case "acl.message.broadcast":
		permissionKey := "message.broadcast"
		permissionVal := message.Content

		if message.To.Path == "all" {
			if permissionVal != "0" && permissionVal != "1" {
				rsp.Message = MsgFormatStateNok(req.Message.Type, req.Message.ClientMsgId, "Content内容不符合要求 0 or 1")
				return nil
			}
			err := SetRoomAcl(
				passToken.TenantCode, passToken.RoomCode,
				permissionKey, permissionVal,
			)
			if err != nil {
				return err
			}
		} else {
			if message.To == nil {
				rsp.Message = MsgFormatStateNok(req.Message.Type, req.Message.ClientMsgId, "To内容不能为nil")
				return nil
			}
			userId, err := strconv.Atoi(message.To.Path)
			if err != nil {
				rsp.Message = MsgFormatStateNok(req.Message.Type, req.Message.ClientMsgId, err.Error())
				return nil
			}
			err = SetUserAcl(
				passToken.TenantCode, passToken.RoomCode, uint64(userId),
				permissionKey, permissionVal,
			)
			if err != nil {
				return err
			}
		}

		rsp.Message = MsgFormatStateOk(req.Message.Type, req.Message.ClientMsgId)
		break
	default:
		rsp.Message = MsgFormatStateNok(req.Message.Type, req.Message.ClientMsgId, "未处理的操作")
	}

	//// TODO (修正广播消息上下文路由)
	//if req.Message.Type == "broadcast" {
	//	passToken := req.Message.From.PassToken
	//	switch req.Message.From.PassToken.UserType {
	//	case "user":
	//		req.Message.To.Path = fmt.Sprintf("room.%s.%s", passToken.TenantCode, passToken.RoomCode)
	//		break
	//	case "manager":
	//		req.Message.To.Path = fmt.Sprintf("room.%s.%s", passToken.TenantCode, passToken.RoomCode)
	//		break
	//	}
	//}
	//
	//// TODO （检查Redis，用户是否已经被禁言）
	//// TODO (鉴权、敏感词过滤)
	//if err := CheckAuthority(req.Message); err != nil {
	//	rsp.Code = 1
	//	rsp.Message = err.Error()
	//	return rsp, nil
	//}
	//if err := CheckWords(req.Message.Content); err != nil {
	//	rsp.Code = 1
	//	rsp.Message = err.Error()
	//	return rsp, nil
	//}
	return nil
}

func (s *Service) Record(ctx context.Context, req *proto_room.RecordReq, rsp *proto_room.RecordRsp) error {
	stream := fmt.Sprintf("im:stream:%s:%s", req.TenantCode, req.RoomCode)
	lastMessageId, count := req.LastMessageId, req.Count
	if lastMessageId == "" {
		lastMessageId = "0"
	}
	if count == 0 || count > 1000 {
		count = 20
	}

	res, err := RedisClusterClient.XRead(&redis.XReadArgs{
		Streams: []string{stream, lastMessageId},
		Count:   int64(count),
		Block:   -1,
	}).Result()
	if err != nil {
		if err == redis.Nil {
			return nil
		}
		return err
	}

	var record []*proto_room.Message
	for _, val := range res {
		for _, message := range val.Messages {
			payload, ok := message.Values["payload"].(string)
			if !ok {
				glog.Errorln("Assert payload is failed")
				RedisClusterClient.XDel(stream, message.ID)
				continue
			}

			msg, err := Codec.DecodeMessage([]byte(payload))
			if err != nil {
				glog.Errorln(err)
				RedisClusterClient.XDel(stream, message.ID)
				continue
			}

			if msg.Type == "message.broadcast" || msg.Type == "message.p2p" {
				msg.ServerMsgId = message.ID
			}

			record = append(record, msg)
		}
	}

	rsp.Record = record
	return nil
}

func (s *Service) Info(ctx context.Context, req *proto_room.InfoReq, rsp *proto_room.InfoRsp) error {
	for _, roomCode := range req.RoomCode {
		count, err := GetUsersCount(req.TenantCode, roomCode)
		if err != nil {
			glog.Errorln(err)
		}
		rsp.List[roomCode] = uint64(count)
	}
	return nil
}