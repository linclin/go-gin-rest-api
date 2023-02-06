package sys

import (
	"go-gin-rest-api/api/v1/sys"
	"go-gin-rest-api/middleware"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// 数据审计日志
func InitChangeLogRouter(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) (R gin.IRoutes) {
	router := r.Group("changelog").Use(authMiddleware.MiddlewareFunc()).Use(middleware.CasbinMiddleware)
	{
		router.POST("/list", sys.GetChangeLog)
	}
	return router
}
