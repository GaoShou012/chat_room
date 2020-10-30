package room

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"wchatv1/common"
	"wchatv1/config"
	proto_room "wchatv1/proto/room"
)

func MakePassToken(ctx *gin.Context) {
	params := &proto_room.MakePassTokenReq{}
	err := ctx.BindJSON(&params)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	rsp, err := config.RoomServiceConfig.ServiceClient().MakePassToken(context.TODO(), params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    "200",
		"msg":     "success",
		"success": true,
		"retry":   false,
		"token":   rsp.Token,
		"data":    "",
	})
}

func SetUserAcl(ctx *gin.Context) {
	var params struct {
		TenantCode string
		UserId     uint64
		Key        string
		Val        string
	}
	err := ctx.BindJSON(&params)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	req := &proto_room.SetTenantUserAclReq{
		TenantCode: params.TenantCode,
		UserId:     params.UserId,
		Key:        params.Key,
		Val:        params.Val,
	}
	rsp, err := config.RoomServiceConfig.ServiceClient().SetTenantUserAcl(context.TODO(), req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	if rsp.Code != 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    "200",
			"msg":     "success",
			"success": false,
			"retry":   false,
			"data":    "",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    "200",
		"msg":     "设置成功",
		"success": true,
		"retry":   false,
		"data":    "",
	})
}

func GetUserAcl(ctx *gin.Context) {
	var params struct {
		TenantCode string
		UserId     uint64
		Key        string
	}
	err := ctx.Bind(&params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	req := &proto_room.GetTenantUserAclReq{
		TenantCode: params.TenantCode,
		UserId:     params.UserId,
		Key:        params.Key,
	}
	rsp, err := config.RoomServiceConfig.ServiceClient().GetTenantUserAcl(context.TODO(), req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	if rsp.Code != 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    "200",
			"msg":     rsp.Desc,
			"success": false,
			"retry":   false,
			"data":    "",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    "200",
		"msg":     "success",
		"success": true,
		"retry":   false,
		"data":    rsp.Val,
	})
}

func Select(ctx *gin.Context) {
	var params struct {
		common.Page
	}
	if err := ctx.BindJSON(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	if err := params.PageCheck(); err != nil {
		common.ResponseError(ctx, err)
		return
	}
}

func SetVirtualUserCounter(ctx *gin.Context) {
	var params struct {
		TenantCode string
		RoomCode   string
		Count      uint64
	}
	if err := ctx.BindJSON(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	req := &proto_room.SetVirtualUserCountReq{
		TenantCode: params.TenantCode,
		RoomCode:   params.RoomCode,
		Count:      params.Count,
	}
	rsp, err := config.RoomServiceConfig.ServiceClient().SetVirtualUserCount(context.TODO(), req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	if rsp.Code != 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    "200",
			"msg":     rsp.Desc,
			"success": false,
			"retry":   false,
			"data":    "",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    "200",
		"msg":     "success",
		"success": true,
		"retry":   false,
		"data":    "",
	})
}

func GetVirtualUserCount(ctx *gin.Context) {
	var params struct {
		TenantCode string
		RoomCode   string
	}
	err := ctx.Bind(&params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	req := &proto_room.GetVirtualUserCountReq{
		TenantCode: params.TenantCode,
		RoomCode:   params.RoomCode,
	}
	rsp, err := config.RoomServiceConfig.ServiceClient().GetVirtualUserCount(context.TODO(), req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	if rsp.Code != 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    "200",
			"msg":     rsp.Desc,
			"success": false,
			"retry":   false,
			"data":    "",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    "200",
		"msg":     "success",
		"success": true,
		"retry":   false,
		"data":    rsp.Count,
	})
}
