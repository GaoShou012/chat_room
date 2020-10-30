package room

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/golang/glog"
	"strconv"
	"time"
)

/*
	租户-房间-用户 有序集合
	可以使用 range 进行分页查询，range不会受 分数 变化

*/

func keyOfTenantsRoomsUsersSortedSet(tenantCode string, roomCode string) string {
	return fmt.Sprintf("im:rooms:user-list:%s:%s", tenantCode, roomCode)
}

type TenantRoomUserSortedSet struct {
	RedisClient *redis.ClusterClient
}

// 删除一个用户从 tenant-room-有序列表
func (t *TenantRoomUserSortedSet) Del(tenantCode string, roomCode string, userId uint64) error {
	key := keyOfTenantsRoomsUsersSortedSet(tenantCode, roomCode)
	_, err := t.RedisClient.ZRem(key, userId).Result()
	return err
}

// 添加一个用户到 tenant-room-有序列表
func (t *TenantRoomUserSortedSet) Set(tenantCode string, roomCode string, userId uint64) error {
	key := keyOfTenantsRoomsUsersSortedSet(tenantCode, roomCode)
	add := &redis.Z{
		Score:  float64(time.Now().UnixNano()),
		Member: userId,
	}
	_, err := t.RedisClient.ZAdd(key, add).Result()
	return err
}

// tenant-room-user 数量
func (t *TenantRoomUserSortedSet) Count(tenantCode string, roomCode string) (int64, error) {
	key := keyOfTenantsRoomsUsersSortedSet(tenantCode, roomCode)
	return t.RedisClient.ZCard(key).Result()
}

func (t *TenantRoomUserSortedSet) Range(tenantCode string, roomCode string, page uint64, pageSize uint64) ([]uint64, error) {
	key := keyOfTenantsRoomsUsersSortedSet(tenantCode, roomCode)
	start := page * pageSize
	stop := start + pageSize - 1
	res, err := t.RedisClient.ZRange(key, int64(start), int64(stop)).Result()
	if err != nil {
		return nil, err
	}

	var users []uint64
	for _, v := range res {
		id, err := strconv.Atoi(v)
		if err != nil {
			glog.Errorln(err)
			continue
		}
		users = append(users, uint64(id))
	}

	return users, nil
}
