package app

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ErrorMessage struct {
	Code int    `yaml:"code"`
	Msg  string `yaml:"msg"`
}

var errorMessages = map[string]ErrorMessage{}

func init() {
	okMessages := map[string]ErrorMessage{
		"OK": {Code: 200, Msg: "成功"},
	}
	// 参数异常返回 10000 +
	paramsErrorMessages := map[string]ErrorMessage{
		"ParamsHandlerErr": {Code: 10001, Msg: "参数解析异常"},
	}
	// 业务异常返回 20000 +
	entryErrorMessages := map[string]ErrorMessage{}
	// 数据库操作类 30000 +
	dbErrorMessages := map[string]ErrorMessage{
		"DBGetErr":         {Code: 30001, Msg: "数据库查询失败"},
		"DBUpdateErr":      {Code: 30002, Msg: "数据库更新失败"},
		"DBDeleteErr":      {Code: 30003, Msg: "数据库删除失败"},
		"MappingDataErr":   {Code: 30004, Msg: "映射数据异常"},
		"DataUnmarshalErr": {Code: 30005, Msg: "数据解析异常"},
	}

	mergeErrorMessages(errorMessages, okMessages, paramsErrorMessages, entryErrorMessages, dbErrorMessages)
}

func XError(c *gin.Context, key string) {
	var res Response
	errMsg := getErrorMessage(key) // 获取对应的错误消息
	res.Msg = errMsg.Msg
	code := errMsg.Code
	c.JSON(http.StatusOK, res.ReturnError(code))
}

func XOK(c *gin.Context, data interface{}, msg string) {
	var res Response
	res.Data = data
	if msg != "" {
		res.Msg = msg
	}
	c.JSON(http.StatusOK, res.ReturnOK())
}

func getErrorMessage(key string) ErrorMessage {
	errMsg, found := errorMessages[key]
	if !found {
		return ErrorMessage{Code: 500, Msg: "Unknown error"}
	}
	return errMsg
}

func mergeErrorMessages(maps ...map[string]ErrorMessage) {
	for _, m := range maps {
		for key, value := range m {
			errorMessages[key] = value
		}
	}
}
