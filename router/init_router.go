package router

import (
	"github.com/gin-gonic/gin"
	"review-bot/api"
	"review-bot/internal/middlewares"
)

func InitRouter() *gin.Engine {
	r := gin.New()

	// 中间件注册
	middlewares.InitMiddleware(r)

	g := r.Group("")
	// 路由注册
	InitSysRouter(g)

	// TODO 注册业务路由

	return r
}

// InitSysRouter 路由注册
func InitSysRouter(r *gin.RouterGroup) *gin.RouterGroup {

	// 测试路由
	r.GET("/ping", func(context *gin.Context) {
		context.JSON(200, gin.H{
			"message": "pong!",
		})
	})

	r.GET("/test/qa", api.QaTest)
	r.GET("/test/qa/list", api.QaListTest)
	return r
}
