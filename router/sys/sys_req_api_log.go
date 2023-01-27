package sys

import (
	"go-gin-rest-api/api/v1/sys"
	"go-gin-rest-api/middleware"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// 系统路由
func InitReqApiLogRouter(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) (R gin.IRoutes) {
	router := r.Group("reqapilog").Use(authMiddleware.MiddlewareFunc()).Use(middleware.CasbinMiddleware)
	{
		router.POST("/list", sys.GetReqApiLog)
		router.GET("/get/:requestid", sys.GetReqApiLogById)
	}
	return router
}
