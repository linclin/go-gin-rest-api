package middleware

import (
	"go-gin-rest-api/models"
	"go-gin-rest-api/pkg/global"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Casbin中间件, 基于RBAC的权限访问控制模型
func CasbinMiddleware(c *gin.Context) {
	AppId, _ := c.Get("AppId")
	// 请求URL路径作为casbin访问资源obj(需先清除path前缀)
	obj := c.Request.URL.Path
	// 请求方式作为casbin访问动作act
	act := c.Request.Method
	// 检查策略
	permResult, err := global.CasbinACLEnforcer.Enforce(AppId, obj, "*", "*", act)
	if err != nil || !permResult {
		c.JSON(http.StatusForbidden, models.Resp{
			Code: models.Forbidden,
			Data: models.ForbiddenMsg,
			Msg:  models.CustomError[models.Forbidden],
		})
		c.Abort()
		return
	} else {
		// 处理请求
		c.Next()
	}

}
