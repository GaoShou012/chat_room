package tenant_api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"wchatv1/config"
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
		// 解密token
		err := operator.DecryptByJwt(token)
		if err != nil {
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
			if err == redis.Nil {
				utils_api.RspState(ctx, 2, "缓存丢失租户信息")
				ctx.Abort()
				return
			}
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

func (c *Auth) Login(ctx *gin.Context) {
	var params struct {
		Username string
		Password string
	}
	{
		err := ctx.BindJSON(&params)
		if err != nil {
			utils_api.RspState(ctx, 1, err)
			ctx.Abort()
			return
		}
	}

	var token string
	var operator *Operator

	{
		// jwt 加密数据，返回token
		operator = &Operator{TenantCode: params.Username}
		str, err := operator.EncryptByJwt([]byte(TokenKey))
		if err != nil {
			utils_api.RspState(ctx, 1, err)
			return
		}
		token = str
	}

	{
		// 登录成功后，设置缓存状态
		key := fmt.Sprintf("im:tenant:enable:%s", operator.TenantCode)
		_, err := c.RedisClient.Set(key, "1", 0).Result()
		if err != nil {
			utils_api.RspState(ctx, 1, err)
			return
		}
	}

	{
		// 返回token
		utils_api.RspData(ctx, 0, nil, token)
		return
	}
}
