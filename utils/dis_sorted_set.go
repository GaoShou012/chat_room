package utils

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis/v7"
	"time"
)

type DisSortedSet struct {
	RedisClient *redis.ClusterClient
}

func (d *DisSortedSet) Add(key string, member string) error {
	z := &redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: member,
	}
	_, err := d.RedisClient.ZAdd(key, z).Result()
	return err
}
func (d *DisSortedSet) GetAll(key string, timeout int64) ([]string, error) {
	opt := &redis.ZRangeBy{
		Min:    fmt.Sprintf("%d", time.Now().Unix()-timeout),
		Max:    "-1",
		Offset: 0,
		Count:  0,
	}
	//return d.RedisClient.ZRangeByLex(key, opt).Result()
	//return d.RedisClient.ZRange(key,-1,time.Now().Unix()-timeout).Result()
	return d.RedisClient.ZRangeByScore(key, opt).Result()
	//return d.RedisClient.ZRange(key, time.Now().Unix()-timeout, time.Now().Unix()+100).Result()
	//return d.RedisClient.ZRange(key,0,-1).Result()
}

func (d *DisSortedSet) GetValid(key string, timeout int64) ([]string, error) {
	res, err := d.RedisClient.ZRangeWithScores(key, 0, -1).Result()
	if err != nil {
		return nil, err
	}
	exp := float64(time.Now().Unix() - timeout)
	var rows []string
	for _, row := range res {
		member, ok := row.Member.(string)
		if !ok {
			return nil, errors.New("assert member failed")
		}
		if row.Score > exp {
			rows = append(rows, member)
		}
	}
	return rows, nil
}

func (d *DisSortedSet) Count(key string, timeout int64) (int64, error) {
	return d.RedisClient.ZCount(key, fmt.Sprintf("%d", time.Now().Unix()-timeout), "0").Result()
}
