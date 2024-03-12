package middlewares

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"net/http"
)

// 限流中间件
func rateLimit(limiter *rate.Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"status":  429,
				"message": "Too many requests",
			})
			return
		}
		c.Next()
	}
}
