package sys

import (
	"go-gin-rest-api/api/v1/sys"
	"go-gin-rest-api/middleware"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// 系统路由
func InitSystemRouter(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) (R gin.IRoutes) {
	router := r.Group("system").Use(authMiddleware.MiddlewareFunc()).Use(middleware.CasbinMiddleware)
	{
		router.POST("/list", sys.GetSystems)
		router.GET("/get/:id", sys.GetSystemById)
		router.POST("/create", sys.CreateSystem)
		router.PATCH("/update/:id", sys.UpdateSystemById)
		router.DELETE("/delete/:id", sys.DeleteSystemById)
		//系统权限
		router.GET("/perm/get/:id", sys.GetSystemPermById)
		router.POST("/perm/create/:id", sys.CreateSystemPerm)
		router.DELETE("/perm/delete/:id", sys.DeleteSystemPermById)

	}
	return router
}
