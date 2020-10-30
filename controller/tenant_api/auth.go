package tenant_api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"wchatv1/config"
	"wchatv1/utils"
	utils_api "wchatv1/utils/api"
)

const (
	TokenKey = "fdm5DmTKp7Jep1GV"
)

type Auth struct {
	RedisClient *redis.ClusterClient
}

func (c *Auth) VerifyByLocalTesting(ctx *gin.Context) {
	var params Operator
	if err := ctx.BindHeader(&params); err != nil {
		utils_api.RspState(ctx, 1, err)
		ctx.Abort()
		return
	}
	SetOperator(ctx, &params)
}

func (c *Auth) VerifyByEncryptedToken(ctx *gin.Context) {
	token := ctx.GetHeader(config.HeaderToken)
	operator := &Operator{}
	{
		res, err := utils.AesEncrypt([]byte(token), []byte(TokenKey), nil, utils.AesModeCBCPk5)
		if err != nil {
			utils_api.RspState(ctx, 1, err)
			ctx.Abort()
			return
		}
		if err := json.Unmarshal(res, operator); err != nil {
			utils_api.RspState(ctx, 1, err)
			ctx.Abort()
			return
		}
	}

	{
		// 从缓存中，检查租户是否有效
		key := fmt.Sprintf("im:tenant:enable:%s", operator.TenantCode)
		res, err := c.RedisClient.Get(key).Result()
		if err != nil {
			utils_api.RspState(ctx, 1, err)
			ctx.Abort()
			return
		}
		if res != "1" {
			utils_api.RspState(ctx, 1, "租户已经被禁用")
			ctx.Abort()
			return
		}
	}

	SetOperator(ctx, operator)
}
