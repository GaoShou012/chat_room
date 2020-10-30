package tenant

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"wchatv1/common"
	"wchatv1/models"
	"wchatv1/utils"
)

func Insert(ctx *gin.Context) {
	var params struct {
		Enable     bool
		TenantCode string
		TenantKey  string
	}
	err := ctx.BindJSON(&params)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	tenant := &models.Tenants{
		Enable:     &params.Enable,
		TenantCode: &params.TenantCode,
		TenantKey:  &params.TenantKey,
	}
	res := utils.DB.Model(tenant).Create(tenant)
	if res.Error != nil {
		common.ResponseError(ctx, res.Error)
		return
	}

	common.ResponseOk(ctx)
}

func Update(ctx *gin.Context) {
	var params struct {
		Id     uint64
		Enable bool
	}
	err := ctx.BindJSON(&params)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	tenant := &models.Tenants{
		Enable: &params.Enable,
	}

	db := utils.DB.Model(tenant)
	db = db.Where("id = ?", params.Id)
	res := db.Updates(tenant)
	if res.Error != nil {
		common.ResponseError(ctx, res.Error)
		return
	}

	common.ResponseOk(ctx)
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

	var total uint64
	var data []models.Tenants

	// TODO SQL条件绑定
	db := utils.DB.Model(models.Tenants{})

	// TODO 查询数据
	db.Count(&total)
	res := db.Offset(params.PageOffset()).Limit(params.PageSize).Find(&data)
	if res.Error != nil {
		common.ResponseError(ctx,res.Error)
		return
	}

	common.ResponseSearchData(ctx,total,data)
}
