package room

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis/v7"
	"strings"
	"time"
)

/*
	网关关联租户房间
	网关宕机后，可以通过此数据结构，移除相应的房间
	Developer: 高手
*/

func KeyOfGwRooms(gwId string) string {
	return fmt.Sprintf("im:frontier-rooms:%s", gwId)
}

type GwAssocTenantRoomItem struct {
	TenantCode string
	RoomCode   string
}

type GwAssocTenantRoomMap struct {
	RedisClient *redis.ClusterClient
}

/*
	绑定一个租户的房间 & frontier 之间的关系
*/
func (g *GwAssocTenantRoomMap) Set(gwId string, tenantCode string, roomCode string) error {
	key := KeyOfGwRooms(gwId)
	roomKey := fmt.Sprintf("%s:%s", tenantCode, roomCode)
	_, err := g.RedisClient.HSet(key, roomKey, time.Now()).Result()
	return err
}

/*
	解绑一个租户的房间 & frontier 之间的关系
*/
func (g *GwAssocTenantRoomMap) Del(gwId string, tenantCode string, roomCode string) error {
	key := KeyOfGwRooms(gwId)
	roomKey := fmt.Sprintf("%s:%s", tenantCode, roomCode)
	_, err := g.RedisClient.HDel(key, roomKey).Result()
	return err
}

/*
	获取 frontier 关联的所有租户房间
*/
func (g *GwAssocTenantRoomMap) GetAll(gwId string) ([]*GwAssocTenantRoomItem, error) {
	key := KeyOfGwRooms(gwId)
	var rooms []*GwAssocTenantRoomItem
	{
		res, err := g.RedisClient.HGetAll(key).Result()
		if err != nil {
			if err == redis.Nil {
				return rooms, nil
			} else {
				return nil, err
			}
		}
		for _, val := range res {
			arr := strings.Split(val, ":")
			if len(arr) != 2 {
				return nil, errors.New("parse val error")
			}

			item := &GwAssocTenantRoomItem{
				TenantCode: arr[0],
				RoomCode:   arr[1],
			}
			rooms = append(rooms, item)
		}
	}
	return rooms, nil
}

/*
	删除 frontier 关联的所有租户房间
*/
func (g *GwAssocTenantRoomMap) DelAll(gwId string) error {
	key := KeyOfGwRooms(gwId)
	_, err := g.RedisClient.Del(key).Result()
	return err
}
