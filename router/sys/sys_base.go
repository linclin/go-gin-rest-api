package sys

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// 基础路由
func InitBaseRouter(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) (R gin.IRoutes) {
	router := r.Group("base")
	{
		router.POST("/auth", authMiddleware.LoginHandler)
		router.POST("/logout", authMiddleware.LogoutHandler)
		router.POST("/refresh_token", authMiddleware.RefreshHandler)
	}
	return router
}
