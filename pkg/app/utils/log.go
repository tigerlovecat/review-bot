package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"sync"
	"time"
)

// 封装日志
var (
	logs      []string
	logsMutex sync.Mutex
)

// Logger 自定义中间件，记录请求的开始和结束日志
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// 在上下文中设置开始时间，以便在请求处理函数中使用
		c.Set("startTime", startTime)

		// 处理请求
		c.Next()

		endTime := time.Now()
		latencyTime := endTime.Sub(startTime)

		// 记录请求的结束日志
		requestLog := fmt.Sprintf("[GIN] %v | %13v | %15s | %-7s %#v\n",
			endTime.Format("2006/01/02 - 15:04:05"),
			latencyTime,
			c.ClientIP(),
			c.Request.Method,
			c.Request.URL.Path,
		)

		logsMutex.Lock()
		logs = append(logs, requestLog)
		logsMutex.Unlock()
	}
}

// OutPutLog 模拟函数中的日志输出
func OutPutLog(message string) {
	currentTime := time.Now()
	logTime := currentTime.Format("2006/01/02 - 15:04:05")

	logsMutex.Lock()
	logs = append(logs, fmt.Sprintf("[FUNCTION] %v: %s\n", logTime, message))
	logsMutex.Unlock()
}
