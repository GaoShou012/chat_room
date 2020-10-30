package tenant_api

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"wchatv1/utils/api"
)

type Users struct {
	RedisClient *redis.ClusterClient
}

/*
	Get
	查询用户列表
*/
func (c *Users) GetList(ctx *gin.Context) {
	var params struct {
		RoomCode string
		api.Page
	}
}
