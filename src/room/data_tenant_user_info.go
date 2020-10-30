package room

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/mitchellh/mapstructure"
)

/*
	租户-用户信息
	set,get 对redis cluster进行储存

	Developer:高手
*/

type TenantUserInfo struct {
	TenantCode string
	UserType   string
	UserId     uint64
	Username   string
	Nickname   string
	UserThumb  string
	UserTags   string
}

func (t *TenantUserInfo) Set(tenantCode string, userId uint64, redisClient *redis.ClusterClient) error {
	key := fmt.Sprintf("im:rooms:user-info:%s:%d", tenantCode, userId)
	_, err := redisClient.HMSet(key,
		"tenantCode", t.TenantCode,
		"userType", t.UserType,
		"userId", t.UserId,
		"username", t.Username,
		"nickname", t.Nickname,
		"userThumb", t.UserThumb,
		"useTags", t.UserTags,
	).Result()
	return err
}
func (t *TenantUserInfo) Get(tenantCode string, userId uint64, redisClient *redis.ClusterClient) error {
	key := fmt.Sprintf("im:rooms:user-info:%s:%d", tenantCode, userId)
	res, err := redisClient.HGetAll(key).Result()
	if err != nil {
		return err
	}
	if err := mapstructure.WeakDecode(res, t); err != nil {
		return err
	}
	return nil
}
