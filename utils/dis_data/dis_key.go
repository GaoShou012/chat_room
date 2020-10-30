package dis_data

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"strconv"
	"time"
)

type DisKey struct {
	RedisClient *redis.ClusterClient
}

func (d *DisKey) Add(key string, member string) error {
	_, err := d.RedisClient.HSet(key, member, fmt.Sprintf("%d", time.Now().Unix())).Result()
	return err
}

func (d *DisKey) Del(key string, member string) error {
	_, err := d.RedisClient.HDel(key, member).Result()
	return err
}

func (d *DisKey) GetAll(key string) (map[string]string, error) {
	return d.RedisClient.HGetAll(key).Result()
}

func (d *DisKey) GetValid(key string, timeout int64) ([]string, error) {
	// 读取所有的成员
	res, err := d.RedisClient.HGetAll(key).Result()
	if err != nil {
		return nil, err
	}

	// 设置超时时间
	exp := time.Now().Unix() - timeout
	var rows []string

	// 遍历所有会员
	// 比较时间戳筛选出有效的会员
	for key, val := range res {
		timestamp, err := strconv.Atoi(val)
		if err != nil {
			return nil, err
		}
		if exp > int64(timestamp) {
			continue
		}
		rows = append(rows, key)
	}

	return rows, nil
}

func (d *DisKey) GetInvalid(key string, timeout int64) ([]string, error) {
	// 读取所有的成员
	res, err := d.RedisClient.HGetAll(key).Result()
	if err != nil {
		return nil, err
	}

	// 设置超时时间
	exp := time.Now().Unix() - timeout
	var rows []string

	// 遍历所有会员
	// 比较时间戳筛选出无效的会员
	for key, val := range res {
		timestamp, err := strconv.Atoi(val)
		if err != nil {
			return nil, err
		}
		if exp < int64(timestamp) {
			continue
		}
		rows = append(rows, key)
	}

	return rows, nil
}
