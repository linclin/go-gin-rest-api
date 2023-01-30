package sys

import (
	"go-gin-rest-api/api/v1/sys"
	"go-gin-rest-api/middleware"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

// 角色
func InitRoleRouter(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) (R gin.IRoutes) {
	router := r.Group("role").Use(authMiddleware.MiddlewareFunc()).Use(middleware.CasbinMiddleware)
	{
		router.POST("/list", sys.GetRoles)
		router.GET("/get/:id", sys.GetRoleById)
		router.POST("/create", sys.CreateRole)
		router.PATCH("/update/:id", sys.UpdateRoleById)
		router.DELETE("/delete/:id", sys.DeleteRoleById)
		//角色权限
		router.GET("/perm/get/:id", sys.GetRolePermById)
		router.POST("/perm/create/:id", sys.CreateRolePerm)
		router.DELETE("/perm/delete/:id", sys.DeleteRolePermById)
		//角色用户
		router.GET("/users/get/:id", sys.GetRoleUsersById)
		router.POST("/users/create/:id", sys.CreateRoleUser)
		router.DELETE("/users/delete/:id", sys.DeleteRoleUserById)
	}
	return router
}
