package sys

import (
	"go-gin-rest-api/api/v1/sys"
	"go-gin-rest-api/middleware"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// 用户权限
func InitUserRouter(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) (R gin.IRoutes) {
	router := r.Group("user").Use(authMiddleware.MiddlewareFunc()).Use(middleware.CasbinMiddleware)
	{
		router.POST("/login", sys.Login)
		router.GET("/info", sys.GetUserInfo)
		router.GET("/perm/:user", sys.GetPermission)
		router.POST("/auth/:user", sys.AuthPermission)
	}
	return router
}
