package tenant_management

import (
	"github.com/gin-gonic/gin"
	"net/http"
	controller_tenant_management "wchatv1/controller/tenant_management"
)

type HttpService struct {
	Auth  *controller_tenant_management.Auth
	Users *controller_tenant_management.Users
}

func (r *HttpService) Cors() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//请求方法
		method := ctx.Request.Method

		// 允许任何源
		ctx.Header("Access-Control-Allow-Origin", "*")
		//服务器支持的所有跨域请求的方法,为了避免浏览次请求的多次'预检'请求
		ctx.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
		//ctx.Header("Access-Control-Allow-Headers", "Token,Content-Type")
		ctx.Header("Access-Control-Allow-Headers", "*")
		// 跨域关键设置 让浏览器可以解析
		ctx.Header("Access-Control-Allow-Credentials", "true")

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			ctx.JSON(http.StatusOK, "Options Request!")
		}
		// 处理请求
		ctx.Next() //  处理请求
	}
}

func (r *HttpService) Route(engine *gin.Engine) {
	// 未经验证的API
	unauthenticated := engine
	{
		unauthenticated.POST("/login", r.Auth.Login)
	}

	// 已经验证的API
	authenticated := engine
	authenticated.Use(r.Auth.VerifyByEncryptedToken)
	{
		authenticated.GET("chatRoom/user/list", r.Users.GetList)
	}
}
