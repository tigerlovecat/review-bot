package middlewares

import (
	"github.com/gin-gonic/gin"
	"log"
	"review-bot/pkg/app"
	"review-bot/pkg/app/utils"
)

func CheckToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		userID, bo := utils.CheckToken(token)
		if bo {
			c.Set("userID", userID)
			c.Next()
		} else {
			// 验证不通过时终止请求
			c.Abort()
			log.Println("--- 校验 token 失败 ---")
			app.XError(c, "ServiceTokenLimitedErr")
			return
		}
	}
}
