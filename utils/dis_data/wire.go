//+build wireinject

package dis_data

import (
	"github.com/go-redis/redis/v7"
	"github.com/google/wire"
)

func NewDisKey(redisClient *redis.ClusterClient) *DisKey {
	wire.Build(
		Provider,
	)
	return nil
}

func newDisKey() *DisKey {
	wire.Build(
		testingProvider,
		Provider,
	)
	return nil
}
