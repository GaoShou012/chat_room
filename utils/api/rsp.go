package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type State struct {
	Code int64  `json:"code"`
	Desc string `json:"desc"`
}

func AssertDesc(desc interface{}) string {
	var descString string
	switch desc.(type) {
	case string:
		descString = desc.(string)
	case error:
		descString = desc.(error).Error()
	case nil:
		descString = ""
	default:
		descString = "Assert desc failed"
	}
	return descString
}

func RspAuthFailed(ctx *gin.Context, code int64, desc interface{}) {
	ctx.JSON(http.StatusOK, gin.H{
		"authFailed": &State{
			Code: code,
			Desc: AssertDesc(desc),
		},
	})
}

func RspState(ctx *gin.Context, code int64, desc interface{}) {
	ctx.JSON(http.StatusOK, gin.H{
		"state": &State{
			Code: code,
			Desc: AssertDesc(desc),
		},
	})
}

func RspData(ctx *gin.Context, code int64, desc interface{}, data interface{}) {
	ctx.JSON(http.StatusOK, gin.H{
		"state": &State{
			Code: code,
			Desc: AssertDesc(desc),
		},
		"data": data,
	})
}

func RspRows(ctx *gin.Context, code int64, desc interface{}, total int64, data interface{}) {
	ctx.JSON(http.StatusOK, gin.H{
		"state": &State{
			Code: code,
			Desc: AssertDesc(desc),
		},
		"total": total,
		"data":  data,
	})
}
