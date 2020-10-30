package common

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func ResponseOk(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"code":    "200",
		"msg":     "success",
		"success": true,
		"retry":   false,
		"data":    "",
	})
}

func ResponseError(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusOK, gin.H{
		"code":    "200",
		"msg":     err.Error(),
		"success": false,
		"retry":   false,
		"data":    "",
	})
}

func ResponseSearchData(ctx *gin.Context, count uint64, data interface{}) {
	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "",
		"count":   count,
		"data":    data,
	})
}
