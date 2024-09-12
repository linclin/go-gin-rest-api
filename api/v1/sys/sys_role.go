package sys

import (
	"go-gin-rest-api/models"
	"go-gin-rest-api/models/sys"
	"go-gin-rest-api/pkg/global"
	"strings"

	"github.com/a8m/rql"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	loggable "github.com/linclin/gorm2-loggable"
	"github.com/samber/lo"
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
		Model:        sys.SysRole{},
		DefaultLimit: -1,
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
	if rqlParams.Sort == "" {
		rqlParams.Sort = "id desc"
	}
	list := make([]sys.SysRole, 0)
	query := global.DB
	count := int64(0)
	err = query.Model(sys.SysRole{}).Where(rqlParams.FilterExp, rqlParams.FilterArgs...).Count(&count).Error
	if err != nil {
		models.FailWithDetailed("", err.Error(), c)
		return
	}
	err = query.Where(rqlParams.FilterExp, rqlParams.FilterArgs...).Limit(rqlParams.Limit).Offset(rqlParams.Offset).Order(rqlParams.Sort).Find(&list).Error
	if err != nil {
		models.FailWithDetailed("", err.Error(), c)
		return
	}
	models.OkWithDataList(list, count, c)
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
	err := global.DB.Where("id = ?", id).First(&sysrole).Error
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
	appId, _ := c.Get("AppId")
	requestId, _ := c.Get("RequestId")
	userName := c.GetHeader("User")
	err = global.DB.Set(loggable.LoggableUserTag, &loggable.User{Name: userName, ID: cast.ToString(requestId), Class: cast.ToString(appId)}).Create(&role).Error
	if err != nil {
		models.FailWithDetailed("", err.Error(), c)
	} else {
		if !strings.HasPrefix(role.Name, "group_") {
			role.Name = "group_" + role.Name
		}
		global.CasbinACLEnforcer.AddPolicy(role.Name, "/*", "*", "*", "GET", "allow")
		models.OkWithData(role, c)
	}
}

// @Summary [系统内部]更新角色
// @Id 4 UpdateRoleById
// @Tags [系统内部]角色
// @version 1.0
// @Accept application/x-json-stream
// @Param	id		path 	string	true		"角色ID"
// @Param	body	body 	sys.SysRole	true		"角色"
// @Success 200 object models.Resp 返回创建
// @Failure 500 object models.Resp 创建失败
// @Security ApiKeyAuth
// @Router /api/v1/role/update/{id} [patch]
func UpdateRoleById(c *gin.Context) {
	id := cast.ToInt(c.Param("id"))
	var role sys.SysRole
	query := global.DB.Where("id = ?", id).First(&role)
	if query.Error != nil {
		models.FailWithDetailed("记录不存在", models.CustomError[models.NotOk], c)
		return
	}
	oldvalue := role
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
	role.Name = oldvalue.Name
	appId, _ := c.Get("AppId")
	requestId, _ := c.Get("RequestId")
	userName := c.GetHeader("User")
	err = global.DB.Set(loggable.LoggablePrevVersion, &oldvalue).
		Set(loggable.LoggableUserTag, &loggable.User{Name: userName, ID: cast.ToString(requestId), Class: cast.ToString(appId)}).
		Model(&role).Omit("ID", "CreatedAt", "UpdatedAt", "DeletedAt").
		Where("id = ?", id).Updates(&role).Error
	if err != nil {
		models.FailWithDetailed("", err.Error(), c)
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
	err := global.DB.Where("id = ?", id).First(&sysrole).Error
	if err != nil {
		models.FailWithMessage(err.Error(), c)
		return
	}
	if lo.Contains([]int{1, 2, 3, 4}, id) {
		models.FailWithDetailed(cast.ToString(id)+"角色不允许删除", models.CustomError[models.NotOk], c)
		return
	}
	err = global.DB.Where("id = ?", id).Delete(&sys.SysRole{}).Error
	if err != nil {
		models.FailWithDetailed("", err.Error(), c)
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
	err := global.DB.Where("id = ?", id).First(&role).Error
	if err != nil {
		models.FailWithDetailed("", err.Error(), c)
	} else {
		filteredNamedPolicy, err := global.CasbinACLEnforcer.GetPermissionsForUser("group_" + role.Name)
		if err != nil {
			models.FailWithDetailed("", err.Error(), c)
		}
		var role_perms []sys.RolePermission
		for key, perm := range filteredNamedPolicy {
			role_perms = append(role_perms, sys.RolePermission{
				ID:     key,
				Role:   perm[0],
				Obj:    perm[1],
				Obj1:   perm[2],
				Obj2:   perm[3],
				Action: perm[4],
				Eft:    perm[5],
			})
		}
		models.OkWithDataList(role_perms, cast.ToInt64(len(role_perms)), c)
	}
}

// @Summary [系统内部]创建角色权限
// @Id CreateRolePerm
// @Tags [系统内部]角色
// @version 1.0
// @Accept application/x-json-stream
// @Param	id		path 	string	true "角色ID"
// @Param	body	body 	sys.RolePermission	true "角色权限"
// @Success 200 object models.Resp 返回创建
// @Failure 500 object models.Resp 创建失败
// @Security ApiKeyAuth
// @Router /api/v1/role/perm/create/{id} [post]
func CreateRolePerm(c *gin.Context) {
	var role sys.SysRole
	var role_perms sys.RolePermission
	id := cast.ToInt(c.Param("id"))
	query := global.DB.Where("id = ?", id).First(&role)
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
	if _, err := global.CasbinACLEnforcer.AddPolicy("group_"+role.Name, role_perms.Obj, role_perms.Obj1, role_perms.Obj2, role_perms.Action, role_perms.Eft); err != nil {
		models.FailWithDetailed("授权错误:"+"group_"+role.Name+" "+role_perms.Obj+" "+" "+role_perms.Obj1+" "+" "+role_perms.Obj2+" "+role_perms.Action+" 错误："+err.Error(), models.CustomError[models.NotOk], c)
		return
	}
	models.OkResult(c)
}

// @Summary [系统内部]删除角色权限
// @Id DeleteRolePermById
// @Tags [系统内部]角色
// @version 1.0
// @Accept application/x-json-stream
// @Param	id		path 	string	true		"角色ID"
// @Param	body	body 	sys.RolePermission	 true "角色权限"
// @Success 204 object models.Resp 返回创建
// @Failure 500 object models.Resp 创建失败
// @Security ApiKeyAuth
// @Router /api/v1/role/perm/delete/{id} [delete]
func DeleteRolePermById(c *gin.Context) {
	id := cast.ToInt(c.Param("id"))
	var role sys.SysRole
	var role_perms sys.RolePermission
	query := global.DB.Where("id = ?", id).First(&role)
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
	if lo.Contains([]int{1, 2, 3, 4}, id) {
		models.FailWithDetailed("", cast.ToString(id)+"角色权限不允许删除", c)
		return
	}
	if _, err := global.CasbinACLEnforcer.RemoveNamedPolicy("p", "group_"+role.Name, role_perms.Obj, role_perms.Obj1, role_perms.Obj2, role_perms.Action, "allow"); err != nil {
		models.FailWithDetailed("删除API系统授权错误:"+"group_"+role.Name+role.Name+" "+role_perms.Obj+" "+" "+role_perms.Obj1+" "+" "+role_perms.Obj2+" "+role_perms.Action+" 错误："+err.Error(), models.CustomError[models.NotOk], c)
		return
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
	err := global.DB.Where("id = ?", id).First(&role).Error
	if err != nil {
		models.FailWithDetailed(err, models.CustomError[models.NotOk], c)
	} else {
		roleUsers, err := global.CasbinACLEnforcer.GetUsersForRole("group_" + role.Name)
		if err != nil {
			models.FailWithDetailed([]string{}, err.Error(), c)
			return
		} else {
			if len(roleUsers) == 0 {
				roleUsers, err = global.CasbinACLEnforcer.GetRoleManager().GetUsers("group_" + role.Name)
				if err != nil {
					models.FailWithDetailed([]string{}, err.Error(), c)
					return
				}
			}
		}
		models.OkWithData(roleUsers, c)
	}
}

// @Summary [系统内部]创建角色用户
// @Id CreateRoleUser
// @Tags [系统内部]角色
// @version 1.0
// @Accept application/x-json-stream
// @Param	id		path 	string	true "角色ID"
// @Param	body	body 	[]string  true "用户"
// @Success 201 object models.Resp 返回创建
// @Failure 500 object models.Resp 创建失败
// @Security ApiKeyAuth
// @Router /api/v1/role/users/create/{id} [post]
func CreateRoleUser(c *gin.Context) {
	var role sys.SysRole
	var role_users []string
	id := cast.ToInt(c.Param("id"))
	query := global.DB.Where("id = ?", id).First(&role)
	if query.Error != nil {
		models.FailWithDetailed("", query.Error.Error(), c)
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
		if err := global.CasbinACLEnforcer.GetRoleManager().AddLink(user, "group_"+role.Name); err != nil {
			models.FailWithDetailed("", "授权错误:"+role.Name+" "+user+" 错误："+err.Error(), c)
			return
		}
		if _, err := global.CasbinACLEnforcer.AddRoleForUser(user, "group_"+role.Name, "2024-09-01 00:00:00", "2054-09-01 00:00:00"); err != nil {
			models.FailWithDetailed("", "授权错误:"+role.Name+" "+user+" 错误："+err.Error(), c)
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
	query := global.DB.Where("id = ?", id).First(&role)
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
