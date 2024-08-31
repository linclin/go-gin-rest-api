package sys

import (
	"go-gin-rest-api/models"
	"go-gin-rest-api/models/sys"
	"go-gin-rest-api/pkg/global"

	"github.com/a8m/rql"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/cast"
)

// @Summary [系统内部]获取角色列表
// @Id 1 GetRoles
// @Tags [系统内部]角色
// @version 1.0
// @Accept application/x-json-stream
// @Param	body body  models.Req	true "RQL查询json"
// @Success 200 object models.Resp 返回列表
// @Failure 400 object models.Resp 查询失败
// @Security ApiKeyAuth
// @Router /api/v1/role/list [post]
func GetRoles(c *gin.Context) {
	rqlQueryParser, err := rql.NewParser(rql.Config{
		Model: sys.SysRole{},
	})
	if err != nil {
		models.FailWithMessage(err.Error(), c)
		return
	}
	// 绑定参数
	var rqlQuery rql.Query
	err = c.ShouldBindJSON(&rqlQuery)
	if err != nil {
		models.FailWithMessage(err.Error(), c)
		return
	}
	rqlParams, err := rqlQueryParser.ParseQuery(&rqlQuery)
	if err != nil {
		models.FailWithMessage(err.Error(), c)
		return
	}

	list := make([]sys.SysRole, 0)
	query := global.Mysql
	query = query.Where(rqlParams.FilterExp, rqlParams.FilterArgs...)
	count := int64(0)
	err = query.Limit(rqlParams.Limit).Offset(rqlParams.Offset).Order(rqlParams.Sort).Find(&list).Count(&count).Error
	if err != nil {
		models.FailWithMessage(err.Error(), c)
		return
	}
	// 返回分页数据
	var resp models.PageData
	// 设置分页参数
	resp.PageInfo.Offset = rqlParams.Offset
	resp.PageInfo.Limit = rqlParams.Limit
	resp.PageInfo.Total = count
	resp.PageInfo.SortBy = rqlParams.Sort
	// 设置数据列表
	resp.List = list
	models.OkWithData(resp, c)
}

// @Summary [系统内部]根据ID获取角色
// @Id 2 GetRoleById
// @Tags [系统内部]角色
// @version 1.0
// @Accept application/x-json-stream
// @Param	id		path 	string	true		"角色ID"
// @Success 200 object models.Resp 返回列表
// @Failure 400 object models.Resp 查询失败
// @Security ApiKeyAuth
// @Router /api/v1/role/get/{id} [get]
func GetRoleById(c *gin.Context) {
	var sysrole sys.SysRole
	id := cast.ToInt(c.Param("id"))
	err := global.Mysql.Where("id = ?", id).First(&sysrole).Error
	if err != nil {
		models.FailWithMessage(err.Error(), c)
	} else {
		models.OkWithData(sysrole, c)
	}
}

// @Summary [系统内部]创建角色
// @Id 3 CreateRole
// @Tags [系统内部]角色
// @version 1.0
// @Accept application/x-json-stream
// @Param	body		body 	sys.SysRole	true		"角色"
// @Success 201 object models.Resp 返回创建
// @Failure 500 object models.Resp 创建失败
// @Security ApiKeyAuth
// @Router /api/v1/role/create [post]
func CreateRole(c *gin.Context) {
	// 绑定参数
	var role sys.SysRole
	err := c.ShouldBindJSON(&role)
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
	err = global.Mysql.Create(&role).Error
	if err != nil {
		models.FailWithDetailed(err, models.CustomError[models.NotOk], c)
	} else {
		models.OkWithData(role, c)
	}
}

// @Summary [系统内部]更新角色
// @Id 4 UpdateRoleById
// @Tags [系统内部]角色
// @version 1.0
// @Accept application/x-json-stream
// @Param	id		path 	string	true		"角色ID"
// @Param	body		body 	sys.SysRole	true		"角色"
// @Success 200 object models.Resp 返回创建
// @Failure 500 object models.Resp 创建失败
// @Security ApiKeyAuth
// @Router /api/v1/role/update/{id} [patch]
func UpdateRoleById(c *gin.Context) {
	id := cast.ToInt(c.Param("id"))
	var role sys.SysRole
	query := global.Mysql.Where("id = ?", id).First(&role)
	if query.Error != nil {
		models.FailWithDetailed("记录不存在", models.CustomError[models.NotOk], c)
		return
	}
	// 绑定参数
	err := c.ShouldBindJSON(&role)
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
	err = query.Updates(role).Error
	if err != nil {
		models.FailWithDetailed(err, models.CustomError[models.NotOk], c)
		return
	} else {
		models.OkWithData(role, c)
	}
}

// @Summary [系统内部]删除角色
// @Id DeleteRoleById
// @Tags [系统内部]角色
// @version 1.0
// @Accept application/x-json-stream
// @Param	id		path 	string	true		"角色ID"
// @Success 200 object models.Resp 返回创建
// @Failure 500 object models.Resp 创建失败
// @Security ApiKeyAuth
// @Router /api/v1/role/delete/{id} [delete]
func DeleteRoleById(c *gin.Context) {
	id := cast.ToInt(c.Param("id"))
	var sysrole sys.SysRole
	err := global.Mysql.Where("id = ?", id).First(&sysrole).Error
	if err != nil {
		models.FailWithMessage(err.Error(), c)
		return
	}
	err = global.Mysql.Where("id = ?", id).Delete(&sys.SysRole{}).Error
	if err != nil {
		models.FailWithDetailed(err, models.CustomError[models.NotOk], c)
	} else {
		global.CasbinACLEnforcer.DeleteRole("group_" + sysrole.Name)
		models.OkResult(c)
	}
}

// @Summary [系统内部]获取角色所有权限
// @Id GetRolePermById
// @Tags [系统内部]角色
// @version 1.0
// @Accept application/x-json-stream
// @Param	id		path 	string	true		"角色ID"
// @Success 200 object models.Resp 返回列表
// @Failure 400 object models.Resp 查询失败
// @Security ApiKeyAuth
// @Router /api/v1/role/perm/get/{id} [get]
func GetRolePermById(c *gin.Context) {
	var role sys.SysRole
	id := cast.ToInt(c.Param("id"))
	err := global.Mysql.Where("id = ?", id).First(&role).Error
	if err != nil {
		models.FailWithDetailed(err, models.CustomError[models.NotOk], c)
	} else {
		filteredNamedPolicy, err := global.CasbinACLEnforcer.GetPermissionsForUser("group_" + role.Name)
		if err != nil {
			models.FailWithDetailed(err, models.CustomError[models.NotOk], c)
		}
		models.OkWithData(filteredNamedPolicy, c)
	}
}

// @Summary [系统内部]创建角色权限
// @Id CreateRolePerm
// @Tags [系统内部]角色
// @version 1.0
// @Accept application/x-json-stream
// @Param	id		path 	string	true "角色ID"
// @Param	body	body 	[]sys.RolePermission	true "角色权限"
// @Success 200 object models.Resp 返回创建
// @Failure 500 object models.Resp 创建失败
// @Security ApiKeyAuth
// @Router /api/v1/role/perm/create/{id} [post]
func CreateRolePerm(c *gin.Context) {
	var role sys.SysRole
	var role_perms []sys.RolePermission
	id := cast.ToInt(c.Param("id"))
	query := global.Mysql.Where("id = ?", id).First(&role)
	if query.Error != nil {
		models.FailWithDetailed(query.Error, models.CustomError[models.NotOk], c)
		return
	}
	// 绑定参数
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
	for _, perm := range role_perms {
		if _, err := global.CasbinACLEnforcer.AddPolicy("group_"+role.Name, perm.Obj, perm.Obj1, perm.Obj2, perm.Action, "allow"); err != nil {
			models.FailWithDetailed("授权错误:"+"group_"+role.Name+" "+perm.Obj+" "+" "+perm.Obj1+" "+" "+perm.Obj2+" "+perm.Action+" 错误："+err.Error(), models.CustomError[models.NotOk], c)
			return
		}
	}
	models.OkResult(c)
}

// @Summary [系统内部]删除角色权限
// @Id DeleteRolePermById
// @Tags [系统内部]角色
// @version 1.0
// @Accept application/x-json-stream
// @Param	id		path 	string	true		"角色ID"
// @Param	body	body 	[]sys.RolePermission	 true "角色权限"
// @Success 204 object models.Resp 返回创建
// @Failure 500 object models.Resp 创建失败
// @Security ApiKeyAuth
// @Router /api/v1/role/perm/delete/{id} [delete]
func DeleteRolePermById(c *gin.Context) {
	id := cast.ToInt(c.Param("id"))
	var role sys.SysRole
	var role_perms []sys.RolePermission
	query := global.Mysql.Where("id = ?", id).First(&role)
	if query.Error != nil {
		models.FailWithDetailed("记录不存在", models.CustomError[models.NotOk], c)
		return
	}
	// 绑定参数
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
	for _, perm := range role_perms {
		if _, err := global.CasbinACLEnforcer.RemoveNamedPolicy("p", "group_"+role.Name, perm.Obj, perm.Obj1, perm.Obj2, perm.Action, "allow"); err != nil {
			models.FailWithDetailed("删除API系统授权错误:"+"group_"+role.Name+role.Name+" "+perm.Obj+" "+" "+perm.Obj1+" "+" "+perm.Obj2+" "+perm.Action+" 错误："+err.Error(), models.CustomError[models.NotOk], c)
			return
		}
	}
	models.OkResult(c)
}

// @Summary [系统内部]获取角色所有用户
// @Id GetRoleUsersById
// @Tags [系统内部]角色
// @version 1.0
// @Accept application/x-json-stream
// @Param	id		path 	string	true		"角色ID"
// @Success 200 object models.Resp 返回列表
// @Failure 400 object models.Resp 查询失败
// @Security ApiKeyAuth
// @Router /api/v1/role/users/get/{id} [get]
func GetRoleUsersById(c *gin.Context) {
	var role sys.SysRole
	id := cast.ToInt(c.Param("id"))
	err := global.Mysql.Where("id = ?", id).First(&role).Error
	if err != nil {
		models.FailWithDetailed(err, models.CustomError[models.NotOk], c)
	} else {
		roleUsers, _ := global.CasbinACLEnforcer.GetUsersForRole("group_" + role.Name)
		models.OkWithData(roleUsers, c)
	}
}

// @Summary [系统内部]创建角色用户
// @Id CreateRoleUser
// @Tags [系统内部]角色
// @version 1.0
// @Accept application/x-json-stream
// @Param	id		path 	string	true "角色ID"
// @Param	body	body 	[]string	true "用户"
// @Success 201 object models.Resp 返回创建
// @Failure 500 object models.Resp 创建失败
// @Security ApiKeyAuth
// @Router /api/v1/role/users/create/{id} [post]
func CreateRoleUser(c *gin.Context) {
	var role sys.SysRole
	var role_users []string
	id := cast.ToInt(c.Param("id"))
	query := global.Mysql.Where("id = ?", id).First(&role)
	if query.Error != nil {
		models.FailWithDetailed(query.Error, models.CustomError[models.NotOk], c)
		return
	}
	// 绑定参数
	err := c.ShouldBindJSON(&role_users)
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
	for _, user := range role_users {
		if _, err := global.CasbinACLEnforcer.AddRoleForUser(user, "group_"+role.Name); err != nil {
			models.FailWithDetailed("授权错误:"+role.Name+" "+user+" 错误："+err.Error(), models.CustomError[models.NotOk], c)
			return
		}
	}
	models.OkResult(c)
}

// @Summary [系统内部]删除角色用户
// @Id DeleteRoleUserById
// @Tags [系统内部]角色
// @version 1.0
// @Accept application/x-json-stream
// @Param	id		path 	string	true		"角色ID"
// @Param	body	body 	[]string	 true "用户"
// @Success 204 object models.Resp 返回创建
// @Failure 500 object models.Resp 创建失败
// @Security ApiKeyAuth
// @Router /api/v1/role/users/delete/{id} [delete]
func DeleteRoleUserById(c *gin.Context) {
	id := cast.ToInt(c.Param("id"))
	var role sys.SysRole
	var role_users []string
	query := global.Mysql.Where("id = ?", id).First(&role)
	if query.Error != nil {
		models.FailWithDetailed("记录不存在", models.CustomError[models.NotOk], c)
		return
	}
	// 绑定参数
	err := c.ShouldBindJSON(&role_users)
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
	for _, user := range role_users {
		if _, err := global.CasbinACLEnforcer.DeleteRoleForUser(user, "group_"+role.Name); err != nil {
			models.FailWithDetailed("删除API系统授权错误:"+role.Name+" "+user+" 错误："+err.Error(), models.CustomError[models.NotOk], c)
			return
		}
	}
	models.OkResult(c)
}
