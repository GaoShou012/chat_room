package room

import (
	"fmt"
	"github.com/go-redis/redis/v7"
)

func KeyOfTenantRoomUserMap(tenantCode string, roomCode string) string {
	return fmt.Sprintf("im:rooms:tenant-room-user-map:%s:%s", tenantCode, roomCode)
}

type TenantRoomUserMap struct {
	RedisClusterClient *redis.ClusterClient
}

// 读取用户uuid
func (t *TenantRoomUserMap) Get(tenantCode string, roomCode string, userId uint64) (string, error) {
	key := KeyOfTenantRoomUserMap(tenantCode, roomCode)
	return t.RedisClusterClient.HGet(key, fmt.Sprintf("%d", userId)).Result()
}

// 设置用户uuid
func (t *TenantRoomUserMap) Set(tenantCode string, roomCode string, userId uint64, uuid string) error {
	key := KeyOfTenantRoomUserMap(tenantCode, roomCode)
	_, err := t.RedisClusterClient.HSet(key, fmt.Sprintf("%d", userId), uuid).Result()
	return err
}

// 移除用户uuid
func (t *TenantRoomUserMap) Del(tenantCode string, roomCode string, userId uint64) error {
	key := KeyOfTenantRoomUserMap(tenantCode, roomCode)
	_, err := t.RedisClusterClient.HDel(key, fmt.Sprintf("%d", userId)).Result()
	return err
}

// 获取用户数量
func (t *TenantRoomUserMap) Len(tenantCode string, roomCode string) (int64, error) {
	key := KeyOfTenantRoomUserMap(tenantCode, roomCode)
	num, err := t.RedisClusterClient.HLen(key).Result()
	return num, err
}
