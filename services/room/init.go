package room

import (
	"github.com/go-redis/redis/v7"
	"runtime"
)

var Codec *codec
var RedisClusterClient *redis.ClusterClient

func InitCodec() {
	Codec = &codec{}
	Codec.Init()
}

func InitRedisClusterClient(addr []string, password string) {
	poolSize := 0
	minIdleConns := runtime.NumCPU() * 10

	// join,leave processor parallel
	poolSize += runtime.NumCPU() * 10
	// record parallel
	poolSize += 100

	RedisClusterClient = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:              addr,
		MaxRedirects:       0,
		ReadOnly:           false,
		RouteByLatency:     false,
		RouteRandomly:      false,
		ClusterSlots:       nil,
		OnNewNode:          nil,
		Dialer:             nil,
		OnConnect:          nil,
		Username:           "",
		Password:           password,
		MaxRetries:         0,
		MinRetryBackoff:    0,
		MaxRetryBackoff:    0,
		DialTimeout:        0,
		ReadTimeout:        0,
		WriteTimeout:       0,
		NewClient:          nil,
		PoolSize:           poolSize,
		MinIdleConns:       minIdleConns,
		MaxConnAge:         0,
		PoolTimeout:        0,
		IdleTimeout:        0,
		IdleCheckFrequency: 0,
		TLSConfig:          nil,
	})
}
