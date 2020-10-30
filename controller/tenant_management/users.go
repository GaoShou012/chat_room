package tenant_api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"wchatv1/src/room"
	"wchatv1/utils/api"
)

type Users struct {
	RedisClient    *redis.ClusterClient
	TenantRoomUser *room.TenantRoomUser
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
	{
		err := ctx.Bind(&params)
		if err != nil {
			api.RspState(ctx, 1, err)
			return
		}
	}

	{
		if err := params.PageCheck(); err != nil {
			api.RspState(ctx, 1, err)
			return
		}
	}

	operator, err := GetOperator(ctx)
	if err != nil {
		api.RspState(ctx, 1, err)
		return
	}

	tenantCode := operator.TenantCode
	roomCode := params.RoomCode
	page := params.Page.Page
	pageSize := params.Page.PageSize

	var count int64
	var rows []interface{}
	{
		fmt.Println(tenantCode, roomCode)
		fmt.Println(page, pageSize)
		num, err := c.TenantRoomUser.Count(tenantCode, roomCode)
		if err != nil {
			api.RspState(ctx, 1, err)
			return
		}
		count = num

		if num == 0 {
			api.RspRows(ctx, 0, nil, count, rows)
			return
		}
	}

	{
		res, err := c.TenantRoomUser.Range(tenantCode, roomCode, page, pageSize)
		if err != nil {
			api.RspState(ctx, 1, err)
			return
		}
		api.RspRows(ctx, 0, nil, count, res)
		return
	}
}
