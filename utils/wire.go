//+build wireinject

package utils

import (
	"github.com/go-redis/redis/v7"
	"github.com/google/wire"
	"wchatv1/utils/dis_data"
)

func NewDisSortedSet(redisClient *redis.ClusterClient) *DisSortedSet {
	wire.Build(
		wire.Struct(new(DisSortedSet), "*"),
	)
	return nil
}

func NewDisHash() *dis_data.DisHash {

}