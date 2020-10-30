package room

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"strconv"
)

/*
	房间的虚假用户数量
*/

func keyOfTenantRoomFakeUser() string {
	return fmt.Sprintf("im:chat-room:fake_user_count")
}

type TenantRoomFakeUser struct {
	RedisClient *redis.ClusterClient
}

/*
	设置房间的虚假用户数量
*/
func (t *TenantRoomFakeUser) SetCount(tenantCode string, roomCode string, count uint64) error {
	key := keyOfTenantRoomFakeUser()
	_, err := t.RedisClient.HSet(key, fmt.Sprintf("%s:%s", tenantCode, roomCode), fmt.Sprintf("%d", count)).Result()
	return err
}

/*
	读取房间的虚假用户数量
*/
func (t *TenantRoomFakeUser) GetCount(tenantCode string, roomCode string) (uint64, error) {
	key := keyOfTenantRoomFakeUser()
	res, err := t.RedisClient.HGet(key, fmt.Sprintf("%s:%s", tenantCode, roomCode)).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		} else {
			return 0, err
		}
	}
	num, err := strconv.Atoi(res)
	if err != nil {
		return 0, err
	}
	return uint64(num), nil
}
