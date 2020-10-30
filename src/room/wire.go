//+build wireinject

package room

import (
	"github.com/go-redis/redis/v7"
	"github.com/google/wire"
)

func newTenantRoomUser() *TenantRoomUser {
	wire.Build(
		testingProvider,
		Provider,
	)
	return nil
}

func NewTenantRoomUser(redisClient *redis.ClusterClient) *TenantRoomUser {
	wire.Build(
		Provider,
	)
	return nil
}
