// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package tenant_management

import (
	"github.com/go-redis/redis/v7"
	"github.com/jinzhu/gorm"
	"wchatv1/controller/tenant_management"
	"wchatv1/src/room"
)

// Injectors from wire.go:

func NewHttpService(db *gorm.DB, redisClient *redis.ClusterClient) *HttpService {
	auth := &tenant_api.Auth{
		RedisClient: redisClient,
	}
	tenantRoomUserMap := &room.TenantRoomUserMap{
		RedisClusterClient: redisClient,
	}
	tenantRoomUserSortedSet := &room.TenantRoomUserSortedSet{
		RedisClient: redisClient,
	}
	tenantRoomUser := &room.TenantRoomUser{
		RedisClient:             redisClient,
		TenantRoomUserMap:       tenantRoomUserMap,
		TenantRoomUserSortedSet: tenantRoomUserSortedSet,
	}
	users := &tenant_api.Users{
		RedisClient:    redisClient,
		TenantRoomUser: tenantRoomUser,
	}
	httpService := &HttpService{
		Auth:  auth,
		Users: users,
	}
	return httpService
}
