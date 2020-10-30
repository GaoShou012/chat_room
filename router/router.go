package router

import "github.com/gin-gonic/gin"

type Router interface {
	Cors() gin.HandlerFunc
	Route(engine *gin.Engine)
}

func Run(tcp string, router Router, mode string) {
	// 默认mode
	if mode == "" {
		mode = gin.ReleaseMode
	}
	gin.SetMode(mode)

	// 初始化
	r := gin.Default()
	r.Use(router.Cors())
	router.Route(r)

	// 启动
	if err := r.Run(tcp); err != nil {
		panic(err)
	}
}

func New(router Router) *gin.Engine {
	r := gin.Default()
	r.Use(router.Cors())
	router.Route(r)
	return r
}
