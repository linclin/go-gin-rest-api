package middleware

import (
	"fmt"
	"go-gin-rest-api/models"
	"go-gin-rest-api/pkg/global"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

// 全局异常处理中间件
func Exception(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			// 将异常写入日志
			global.Log.Error(fmt.Sprintf("未知panic异常: %v\n堆栈信息: %v", err, string(debug.Stack())))
			models.FailWithDetailed(fmt.Sprintf("未知panic异常: %v\n堆栈信息: %v", err, string(debug.Stack())), models.CustomError[models.InternalServerError], c)
			c.Abort()
			return
		}
	}()
	c.Next()
}
