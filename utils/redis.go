package utils

import (
	"github.com/go-redis/redis"
	"runtime"
	"time"
)

var RedisClusterClient *redis.ClusterClient

func InitRedisClusterClient(addr []string, password string) {
	cli := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:              addr,
		MaxRedirects:       0,
		ReadOnly:           false,
		RouteByLatency:     false,
		RouteRandomly:      false,
		ClusterSlots:       nil,
		OnNewNode:          nil,
		OnConnect:          nil,
		Password:           password,
		MaxRetries:         0,
		MinRetryBackoff:    0,
		MaxRetryBackoff:    0,
		DialTimeout:        0,
		ReadTimeout:        time.Second*3,
		WriteTimeout:       0,
		PoolSize:           runtime.NumCPU()*3,
		MinIdleConns:       runtime.NumCPU(),
		MaxConnAge:         0,
		PoolTimeout:        0,
		IdleTimeout:        0,
		IdleCheckFrequency: 0,
		TLSConfig:          nil,
	})

	RedisClusterClient = cli
}