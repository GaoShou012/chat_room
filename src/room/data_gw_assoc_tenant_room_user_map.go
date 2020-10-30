package room

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"strconv"
)

func KeyOfGwAssocTenantRoomsUsers(gwId string, tenantCode string, roomCode string) string {
	return fmt.Sprintf("im:gw-t-r-user-list:%s:%s:%s", gwId, tenantCode, roomCode)
}

type GwAssocTenantRoomUserMapItem struct {
	UserId uint64
	Uuid   string
}

type GwAssocTenantRoomUserMap struct {
	RedisClient *redis.ClusterClient
}

/*
	获取网关 关联的租户-房间 用户数量
*/
func (g *GwAssocTenantRoomUserMap) Len(gwId string, tenantCode string, roomCode string) (int64, error) {
	key := KeyOfGwAssocTenantRoomsUsers(gwId, tenantCode, roomCode)
	return g.RedisClient.HLen(key).Result()
}

/*
	设置网关 关联的租户-房间-用户 uuid
*/
func (g *GwAssocTenantRoomUserMap) Set(gwId string, tenantCode string, roomCode string, userId uint64, uuid string) error {
	key := KeyOfGwAssocTenantRoomsUsers(gwId, tenantCode, roomCode)
	_, err := g.RedisClient.HSet(key, fmt.Sprintf("%d", userId), uuid).Result()
	return err
}

/*
	获取网关 关联的 租户-房间-用户 uuid
*/
func (g *GwAssocTenantRoomUserMap) Get(gwId string, tenantCode string, roomCode string, userId uint64) (string, error) {
	key := KeyOfGwAssocTenantRoomsUsers(gwId, tenantCode, roomCode)
	return g.RedisClient.HGet(key, fmt.Sprintf("%d", userId)).Result()
}

/*
	移除网关 关联的 租户-房间-用户
*/
func (g *GwAssocTenantRoomUserMap) Del(gwId string, tenantCode string, roomCode string, userId uint64) error {
	key := KeyOfGwAssocTenantRoomsUsers(gwId, tenantCode, roomCode)
	_, err := g.RedisClient.HDel(key, fmt.Sprintf("%d", userId)).Result()
	return err
}

/*
	获取网关 关联的 租户-房间-用户
*/
func (g *GwAssocTenantRoomUserMap) GetAll(gwId string, tenantCode string, roomCode string) ([]*GwAssocTenantRoomUserMapItem, error) {
	var users []*GwAssocTenantRoomUserMapItem
	key := KeyOfGwAssocTenantRoomsUsers(gwId, tenantCode, roomCode)
	res, err := g.RedisClient.HGetAll(key).Result()
	if err != nil {
		return nil, err
	}
	for key, val := range res {
		userId, err := strconv.Atoi(key)
		if err != nil {
			return nil, err
		}
		user := &GwAssocTenantRoomUserMapItem{
			UserId: uint64(userId),
			Uuid:   val,
		}
		users = append(users, user)
	}

	return users, nil
}

/*
	移除所有 网关关联的 租户-房间-用户
*/
func (g *GwAssocTenantRoomUserMap) DelAll(gwId string, tenantCode string, roomCode string) error {
	key := KeyOfGwAssocTenantRoomsUsers(gwId, tenantCode, roomCode)
	_, err := g.RedisClient.Del(key).Result()
	return err
}
