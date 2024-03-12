package middlewares

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"path"
	"sync"
	"time"
)

var (
	logsMutex sync.Mutex
)

const (
	red    = 31
	yellow = 33
	blue   = 36
	gray   = 37
)

type LogFormatter struct {
}

func (t LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	//根据不同的level展示颜色
	var levelColor int
	switch entry.Level {
	case logrus.DebugLevel, logrus.TraceLevel:
		levelColor = gray
	case logrus.WarnLevel:
		levelColor = yellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		levelColor = red
	default:
		levelColor = blue
	}
	//字节缓冲区
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}
	//自定义日期格式
	timestamp := entry.Time.Format("2006-01-02 15:04:06")
	if entry.HasCaller() {
		//自定义文件路径
		funcVal := entry.Caller.Function
		fileVal := fmt.Sprintf("%s:%d", path.Base(entry.Caller.File), entry.Caller.Line)
		//自定义输出格式
		fmt.Fprintf(b, "[%s] \033[%dm[%s]\033[0m [ %s %s ] \033[%dm%s\033[0m \n", timestamp, levelColor, entry.Level, fileVal, funcVal, levelColor, entry.Message)
	} else {
		fmt.Fprintf(b, "[%s] \033[%dm[%s]\033[0m  \033[%dm%s\033[0m \n", timestamp, levelColor, entry.Level, levelColor, entry.Message)
	}
	return b.Bytes(), nil
}

// InitLogger 初始化全局变量
func InitLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	logger.SetReportCaller(true)
	logger.SetFormatter(&LogFormatter{})
	return logger
}

// LoggerMiddleware 中间件
func LoggerMiddleware(logger logrus.FieldLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// 创建新的日志记录器对象并放入上下文
		reqLogger := logger.WithFields(logrus.Fields{
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
			"params": c.Request.URL.Query().Encode(),
		})
		c.Set("logger", reqLogger)

		// 打印请求开始信息
		reqLogger.Infof("----------- {%v} 请求开始 -----------", c.Request.URL.Path)

		// 请求处理
		c.Next()

		// 打印请求结束信息
		reqLogger.WithFields(logrus.Fields{
			"duration": time.Since(start),
		}).Infof("----------- {%v} 请求结束 -----------", c.Request.URL.Path)
	}
}
