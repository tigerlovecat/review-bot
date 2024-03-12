package middlewares

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func InitMiddleware(r *gin.Engine) {
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// 初始化全部变量
	logger := InitLogger()
	r.Use(LoggerMiddleware(logger))

	// 配置限流参数
	// 配置限流 - 这里以每秒钟1个请求为例
	limiter := rate.NewLimiter(1000, 10)
	// 应用限流中间件
	r.Use(rateLimit(limiter))
}
