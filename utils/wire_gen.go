// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package utils

import (
	"github.com/go-redis/redis/v7"
)

import (
	_ "github.com/go-sql-driver/mysql"
)

// Injectors from wire.go:

func NewDisSortedSet(redisClient *redis.ClusterClient) *DisSortedSet {
	disSortedSet := &DisSortedSet{
		RedisClient: redisClient,
	}
	return disSortedSet
}
