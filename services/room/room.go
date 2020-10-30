package room

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/golang/glog"
	"strconv"
	"strings"
	"time"
	proto_room "wchatv1/proto/room"
)

func SetVirtualUserCount(tenantCode string, roomCode string, count uint64) error {
	key := fmt.Sprintf("im:tenant:virtual:user:count:%s", tenantCode)
	_, err := RedisClusterClient.HSet(key, roomCode, fmt.Sprintf("%d", count)).Result()
	return err
}
func GetVirtualUserCount(tenantCode string, roomCode string) (uint64, error) {
	key := fmt.Sprintf("im:tenant:virtual:user:count:%s", tenantCode)
	val, err := RedisClusterClient.HGet(key, roomCode).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		} else {
			return 0, err
		}
	}
	num, err := strconv.Atoi(val)
	if err != nil {
		return 0, err
	}

	return uint64(num), nil
}

func SetTenantUserAcl(tenantCode string, userId uint64, field string, value string) error {
	key := fmt.Sprintf("im:tenant:user:acl:%s:%d", tenantCode, userId)
	_, err := RedisClusterClient.HSet(key, field, value).Result()
	return err
}
func GetTenantUserAcl(tenantCode string, userId uint64, field string) (string, error) {
	key := fmt.Sprintf("im:tenant:user:acl:%s:%d", tenantCode, userId)
	res, err := RedisClusterClient.HGet(key, field).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil
		} else {
			return "", err
		}
	}
	return res, nil
}

func SetRoomAcl(tenantCode string, roomCode string, field string, value string) error {
	key := fmt.Sprintf("im:%s:%s:acl", tenantCode, roomCode)
	_, err := RedisClusterClient.HSet(key, field, value).Result()
	return err
}
func GetRoomAcl(tenantCode string, roomCode string, field string) (string, error) {
	key := fmt.Sprintf("im:%s:%s:acl", tenantCode, roomCode)
	res, err := RedisClusterClient.HGet(key, field).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil
		} else {
			return "", err
		}
	}
	return res, nil
}
func SetUserAcl(tenantCode string, roomCode string, userId uint64, field string, value string) error {
	key := fmt.Sprintf("im:%s:%s:%d:acl", tenantCode, roomCode, userId)
	_, err := RedisClusterClient.HSet(key, field, value).Result()
	return err
}
func GetUserAcl(tenantCode string, roomCode string, userId uint64, field string) (string, error) {
	key := fmt.Sprintf("im:%s:%s:%d:acl", tenantCode, roomCode, userId)
	res, err := RedisClusterClient.HGet(key, field).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil
		} else {
			return "", err
		}
	}
	return res, nil
}

func GetRoomPermission(tenantCode string, roomCode string, k string) (string, error) {
	key := fmt.Sprintf("im:%s:%s:permission", tenantCode, roomCode)
	res, err := RedisClusterClient.HGet(key, k).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil
		} else {
			return "", err
		}
	}
	return res, nil
}
func GetUserPermission(tenantCode string, roomCode string, userId uint64, k string) (string, error) {
	key := fmt.Sprintf("im:%s:%s:%d:permission", tenantCode, roomCode, userId)
	res, err := RedisClusterClient.HGet(key, k).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil
		} else {
			return "", err
		}
	}
	return res, nil
}

func FrontierPing(frontierId string) {
	key := fmt.Sprintf("im:frontier:heartbeat")
	RedisClusterClient.HSet(key, frontierId, time.Now().Unix()).Result()
}
func FrontierHeartbeatList() (map[string]string, error) {
	key := fmt.Sprintf("im:frontier:heartbeat")
	return RedisClusterClient.HGetAll(key).Result()
}
func IncrUsersCount(tenantCode string, roomCode string, frontierId string) (int64, error) {
	key := fmt.Sprintf("im:rooms:users:counter:%s:%s:%s", tenantCode, roomCode, frontierId)
	return RedisClusterClient.Incr(key).Result()
}
func DecrUsersCount(tenantCode string, roomCode string, frontierId string) (int64, error) {
	key := fmt.Sprintf("im:rooms:users:counter:%s:%s:%s", tenantCode, roomCode, frontierId)
	return RedisClusterClient.Decr(key).Result()
}
func GetUsersCountOld(tenantCode string, roomCode string) (count int64, err error) {
	now := time.Now().Unix()
	m, err := FrontierHeartbeatList()
	if err != nil {
		return
	}
	pipe := RedisClusterClient.TxPipeline()
	for key, val := range m {
		timestamp, err := strconv.Atoi(val)
		if err != nil {
			glog.Errorln(err)
			continue
		}
		if now-int64(timestamp) > 30 {
			continue
		}
		pipe.Get(fmt.Sprintf("im:rooms:users:counter:%s:%s:%s", tenantCode, roomCode, key))
	}
	res, err := pipe.Exec()
	if err != nil {
		if err == redis.Nil {
			err = nil
		} else {
			return
		}
	}
	for _, val := range res {
		if val.Err() == redis.Nil {
			continue
		}

		row := strings.Split(val.String(), " ")
		num, err := strconv.Atoi(row[2])
		if err != nil {
			glog.Errorln(err)
			continue
		}
		count += int64(num)
	}

	return
}
func GetUsersCount(tenantCode string, roomCode string) (count int64, err error) {
	now := time.Now().Unix()
	m, err := FrontierHeartbeatList()
	if err != nil {
		return
	}

	for key, val := range m {
		timestamp, err := strconv.Atoi(val)
		if err != nil {
			glog.Errorln(err)
			continue
		}
		if now-int64(timestamp) > 30 {
			continue
		}
		usersCounterKey := fmt.Sprintf("im:rooms:users:counter:%s:%s:%s", tenantCode, roomCode, key)
		str, err := RedisClusterClient.Get(usersCounterKey).Result()
		if err != nil {
			glog.Errorln(err)
			continue
		}

		num, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			glog.Errorln(err)
			continue
		}

		count += num
	}

	return
}

func Broker(tenantCode string, roomCode string, message *proto_room.Message) (string, error) {
	j, err := json.Marshal(message)
	if err != nil {
		return "", err
	}
	stream := fmt.Sprintf("%s:%s", tenantCode, roomCode)
	values := make(map[string]interface{})
	values["payload"] = j
	xAddArgs := &redis.XAddArgs{
		Stream:       stream,
		MaxLen:       200,
		MaxLenApprox: 0,
		ID:           "*",
		Values:       values,
	}
	return RedisClusterClient.XAdd(xAddArgs).Result()
}

func GetLastMessageId(tenantCode string, roomCode string) (string, error) {
	stream := fmt.Sprintf("im:stream:%s:%s", tenantCode, roomCode)
	res, err := RedisClusterClient.XRevRangeN(stream, "+", "-", 1).Result()
	if err != nil {
		return "", err
	}
	if res == nil {
		return "", nil
	}
	return res[0].ID, nil
}

func Stream(tenantCode string, roomCode string, message *proto_room.Message) (string, error) {
	j, err := Codec.EncodeMessage(message)
	if err != nil {
		return "", err
	}
	stream := fmt.Sprintf("im:stream:%s:%s", tenantCode, roomCode)
	values := make(map[string]interface{})
	values["payload"] = j
	xAddArgs := &redis.XAddArgs{
		Stream:       stream,
		MaxLen:       50,
		MaxLenApprox: 0,
		ID:           "*",
		Values:       values,
	}
	str, err := RedisClusterClient.XAdd(xAddArgs).Result()
	return str, err
}

func GetRecord(tenantCode string, roomCode string, lastMessageId string, count uint64) ([]redis.XStream, error) {
	stream := fmt.Sprintf("im:stream:%s:%s", tenantCode, roomCode)
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
			return nil, err
		}
	}
	return res, err
}

func DelMessageById(tenantCode string, roomCode string, msgId string) {
	stream := fmt.Sprintf("im:stream:%s:%s", tenantCode, roomCode)
	RedisClusterClient.XDel(stream, msgId)
}
func GetMessageById(tenantCode string, roomCode string, msgId string) (*proto_room.Message, error) {
	stream := fmt.Sprintf("im:stream:%s:%s", tenantCode, roomCode)
	res, err := RedisClusterClient.XRange(stream, msgId, msgId).Result()
	if err != nil {
		return nil, err
	}
	if len(res) != 1 {
		return nil, fmt.Errorf("消息不存在")
	}
	payload, ok := res[0].Values["payload"].(string)
	if !ok {
		RedisClusterClient.XDel(stream, msgId)
		return nil, fmt.Errorf("Assert payload is failed")
	}

	msg := &proto_room.Message{}
	if err := json.Unmarshal([]byte(payload), msg); err != nil {
		RedisClusterClient.XDel(stream, msgId)
		return nil, err
	}
	msg.ServerMsgId = res[0].ID
	return msg, nil
}
