package app

import (
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)

// Error 失败处理
func Error(c *gin.Context, code int, err error, msg string) {
	var res Response
	res.Msg = err.Error()
	if msg != "" {
		res.Msg = msg
	}
	//logger.Errorf(res.Msg, zap.String("msg", msg), zap.Error(err))
	c.JSON(http.StatusOK, res.ReturnError(code))
}

// OK 通常成功数据处理
func OK(c *gin.Context, data interface{}, msg string) {
	var res Response
	res.Data = data
	if msg != "" {
		res.Msg = msg
	}
	c.JSON(http.StatusOK, res.ReturnOK())
}

// PageOK 分页数据处理
func PageOK(c *gin.Context, result interface{}, count int64, offset int, limit int, msg string) {
	var res PageResponse
	if reflect.ValueOf(result).IsNil() {
		res.Data.List = []int{}
	} else {
		res.Data.List = result
	}
	res.Data.Count = count
	res.Data.Offset = offset
	res.Data.Limit = limit
	if msg != "" {
		res.Msg = msg
	}
	c.JSON(http.StatusOK, res.ReturnOK())
}
