package sys

import (
	"go-gin-rest-api/models"
	"go-gin-rest-api/models/sys"
	"go-gin-rest-api/pkg/global"

	"github.com/casbin/casbin/v2"
	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// @Summary [系统内部]用户登录
// @Id Login
// @Tags [系统内部]角色
// @version 1.0
// @Accept application/x-json-stream
// @Param  code query string	true "code"
// @Param  state query string	true "state"
// @Success 200 object models.Resp 返回列表
// @Failure 400 object models.Resp 查询失败
// @Security ApiKeyAuth
// @Router /api/v1/user/login [post]
func Login(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")
	token, err := casdoorsdk.GetOAuthToken(code, state)
	if err != nil {
		models.FailWithMessage(err.Error(), c)
		return
	}
	models.OkWithData(token, c)
}

// @Summary [系统内部]用户登录
// @Id GetUserInfo
// @Tags [系统内部]角色
// @version 1.0
// @Accept application/x-json-stream
// @Success 200 object models.Resp 返回列表
// @Failure 400 object models.Resp 查询失败
// @Security ApiKeyAuth
// @Router /api/v1/user/info [get]
func GetUserInfo(c *gin.Context) {
	token := c.GetHeader("X-Auth-Token")
	claims, err := casdoorsdk.ParseJwtToken(token)
	if err != nil {
		models.FailWithMessage(err.Error(), c)
		return
	}
	models.OkWithData(claims.User, c)
}

// @Summary [系统内部]获取指定用户全部权限
// @Id GetPermission
// @Tags [系统内部]角色
// @version 1.0
// @Accept application/x-json-stream
// @Param	user path string	true "用户名"
// @Success 200 object models.Resp 返回列表
// @Failure 400 object models.Resp 查询失败
// @Security ApiKeyAuth
// @Router /api/v1/user/perm/{user} [get]
func GetPermission(c *gin.Context) {
	user := c.Param("user")
	group, err := global.CasbinACLEnforcer.GetImplicitRolesForUser(user)
	if err != nil {
		models.FailWithMessage(err.Error(), c)
	}
	perm, err := casbin.CasbinJsGetPermissionForUser(global.CasbinACLEnforcer, user)
	if err != nil {
		models.FailWithMessage(err.Error(), c)
	} else {
		models.OkWithData(map[string]interface{}{"group": group, "perm": perm}, c)
	}
}

// @Summary [系统内部]用户操作鉴权
// @Id AuthPermission
// @Tags [系统内部]角色
// @version 1.0
// @Accept application/x-json-stream
// @Param	user	path 	string	true "用户名"
// @Param	body	body 	[]sys.RolePermission	true "权限"
// @Success 200 object models.Resp 返回列表
// @Failure 400 object models.Resp 查询失败
// @Security ApiKeyAuth
// @Router /api/v1/user/auth/{user} [post]
func AuthPermission(c *gin.Context) {
	user := c.Param("user")
	var role_perms sys.RolePermission
	err := c.ShouldBindJSON(&role_perms)
	if err != nil {
		// 获取validator.ValidationErrors类型的errors
		errs, ok := err.(validator.ValidationErrors)
		var errInfo interface{}
		if !ok {
			// 非validator.ValidationErrors类型错误直接返回
			errInfo = err
		} else {
			// validator.ValidationErrors类型错误则进行翻译
			errInfo = errs.Translate(global.Translator) // 翻译校验错误提示
		}
		models.FailWithDetailed(errInfo, models.CustomError[models.NotOk], c)
		return
	}
	ok, reason, err := global.CasbinACLEnforcer.EnforceEx(user, role_perms.Obj, role_perms.Obj1, role_perms.Obj2, role_perms.Action)
	if err != nil {
		models.FailWithMessage(err.Error(), c)
	} else {
		models.OkWithData(map[string]interface{}{"auth": ok, "data": reason}, c)
	}

}
