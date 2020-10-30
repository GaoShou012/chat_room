package utils

import (
	"github.com/go-redis/redis/v7"
	"github.com/google/wire"
	"github.com/jinzhu/gorm"
	"runtime"
)

var testingProvider = wire.NewSet(
	newRedis,
	newMysql,
)

func newRedis() *redis.ClusterClient {
	addr := []string{
		"192.168.56.101:9001",
		"192.168.56.101:9002",
		"192.168.56.101:9003",
		"192.168.56.101:9004",
		"192.168.56.101:9005",
		"192.168.56.101:9006",
	}
	password := ""
	minIdleConns := runtime.NumCPU() * 10
	poolSize := runtime.NumCPU() * 20

	client := redis.NewClusterClient(&redis.ClusterOptions{
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
	return client
}
func newMysql() *gorm.DB {
	return nil
}
