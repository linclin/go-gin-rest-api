package sys

import (
	"go-gin-rest-api/models"
	"go-gin-rest-api/models/sys"
	"go-gin-rest-api/pkg/global"
	"net/http"
	"strings"

	"github.com/a8m/rql"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	loggable "github.com/linclin/gorm2-loggable"
	"github.com/spf13/cast"
)

// @Summary [系统内部]获取系统列表
// @Id GetSystems
// @Tags [系统内部]系统
// @version 1.0
// @Accept application/x-json-stream
// @Param	body		body 	models.Req	true		"RQL查询json"
// @Success 200 object models.Resp 返回列表
// @Failure 400 object models.Resp 查询失败
// @Security ApiKeyAuth
// @Router /api/v1/system/list [post]
func GetSystems(c *gin.Context) {
	rqlQueryParser, err := rql.NewParser(rql.Config{
		Model:        sys.SysSystem{},
		DefaultLimit: -1,
	})
	if err != nil {
		models.FailWithDetailed(err, models.CustomError[models.NotOk], c)
		return
	}
	// 绑定参数
	var rqlQuery rql.Query
	err = c.ShouldBindJSON(&rqlQuery)
	if err != nil {
		models.FailWithDetailed(err, models.CustomError[models.NotOk], c)
		return
	}
	rqlParams, err := rqlQueryParser.ParseQuery(&rqlQuery)
	if err != nil {
		models.FailWithDetailed(err, models.CustomError[models.NotOk], c)
		return
	}
	if rqlParams.Sort == "" {
		rqlParams.Sort = "id desc"
	}
	list := make([]sys.SysSystem, 0)
	query := global.Mysql
	count := int64(0)
	err = query.Model(sys.SysSystem{}).Where(rqlParams.FilterExp, rqlParams.FilterArgs...).Count(&count).Error
	if err != nil {
		models.FailWithDetailed(err, models.CustomError[models.NotOk], c)
		return
	}
	err = query.Where(rqlParams.FilterExp, rqlParams.FilterArgs...).Limit(rqlParams.Limit).Offset(rqlParams.Offset).Order(rqlParams.Sort).Find(&list).Error
	if err != nil {
		models.FailWithDetailed(err, models.CustomError[models.NotOk], c)
		return
	}
	models.OkWithDataList(list, count, c)
}

// @Summary [系统内部]根据ID获取系统
// @Id GetSystemById
// @Tags [系统内部]系统
// @version 1.0
// @Accept application/x-json-stream
// @Param	id		path 	string	true		"系统ID"
// @Success 200 object models.Resp 返回列表
// @Failure 400 object models.Resp 查询失败
// @Security ApiKeyAuth
// @Router /api/v1/system/get/{id} [get]
func GetSystemById(c *gin.Context) {
	var system sys.SysSystem
	id := c.Param("id")
	err := global.Mysql.Where("id = ?", id).First(&system).Error
	if err != nil {
		models.FailWithDetailed(err, models.CustomError[models.NotOk], c)
	} else {
		models.OkWithData(system, c)
	}
}

// @Summary [系统内部]创建系统
// @Id CreateSystem
// @Tags [系统内部]系统
// @version 1.0
// @Accept application/x-json-stream
// @Param	body		body 	sys.SysSystem	true		"系统"
// @Success 201 object models.Resp 返回创建
// @Failure 500 object models.Resp 创建失败
// @Security ApiKeyAuth
// @Router /api/v1/system/create [post]
func CreateSystem(c *gin.Context) {
	// 绑定参数
	var system sys.SysSystem
	err := c.ShouldBindJSON(&system)
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
	if !strings.HasPrefix(system.AppId, "api-") {
		system.AppId = "api-" + system.AppId
	}
	appId, _ := c.Get("AppId")
	requestId, _ := c.Get("RequestId")
	userName := c.GetHeader("User")
	err = global.Mysql.Set(loggable.LoggableUserTag, &loggable.User{Name: userName, ID: cast.ToString(requestId), Class: cast.ToString(appId)}).Create(&system).Error
	if err != nil {
		models.FailWithDetailed(err, models.CustomError[models.NotOk], c)
	} else {
		global.CasbinACLEnforcer.AddPolicy(system.AppId, "/*", "*", "*", "(GET)", "allow")
		global.CasbinACLEnforcer.AddPolicy(system.AppId, "/api/v1/system/*", "*", "*", "(GET)", "deny")
		global.CasbinACLEnforcer.AddPolicy(system.AppId, "/api/v1/role/*", "*", "*", "(GET)", "deny")
		global.CasbinACLEnforcer.AddPolicy(system.AppId, "/api/v1/permission/*", "*", "*", "(GET)", "deny")
		models.OkWithData(system, c)
	}
}

// @Summary [系统内部]更新系统
// @Id UpdateSystemById
// @Tags [系统内部]系统
// @version 1.0
// @Accept application/x-json-stream
// @Param	id		path 	string	true		"系统ID"
// @Param	body		body 	sys.SysSystem	true		"系统"
// @Success 200 object models.Resp 返回创建
// @Failure 500 object models.Resp 创建失败
// @Security ApiKeyAuth
// @Router /api/v1/system/update/{id} [patch]
func UpdateSystemById(c *gin.Context) {
	id := cast.ToInt(c.Param("id"))
	var system sys.SysSystem
	query := global.Mysql.Where("id = ?", id).First(&system)
	if query.Error != nil {
		models.FailWithDetailed("记录不存在", models.CustomError[models.NotOk], c)
		return
	}
	oldvalue := system
	// 绑定参数
	err := c.ShouldBindJSON(&system)
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
	err = global.Mysql.Set(loggable.LoggablePrevVersion, &oldvalue).
		Set(loggable.LoggableUserTag, &loggable.User{Name: userName, ID: cast.ToString(requestId), Class: cast.ToString(appId)}).
		Model(&system).Omit("ID", "CreatedAt", "UpdatedAt", "DeletedAt").
		Where("id = ?", id).Updates(&system).Error
	if err != nil {
		models.FailWithDetailed(err, models.CustomError[models.NotOk], c)
	} else {
		models.OkWithData(system, c)
	}
}

// @Summary [系统内部]删除系统
// @Id DeleteSystemById
// @Tags [系统内部]系统
// @version 1.0
// @Accept application/x-json-stream
// @Param	id		path 	string	true		"系统ID"
// @Success 204 object models.Resp 返回创建
// @Failure 500 object models.Resp 创建失败
// @Security ApiKeyAuth
// @Router /api/v1/system/delete/{id} [delete]
func DeleteSystemById(c *gin.Context) {
	id := cast.ToInt(c.Param("id"))
	var system sys.SysSystem
	appId, _ := c.Get("AppId")
	requestId, _ := c.Get("RequestId")
	userName := c.GetHeader("User")
	query := global.Mysql.Set(loggable.LoggableUserTag, &loggable.User{Name: userName, ID: cast.ToString(requestId), Class: cast.ToString(appId)}).
		Where("id = ?", id).First(&system)
	if query.Error != nil {
		models.FailWithDetailed("记录不存在", models.CustomError[models.NotOk], c)
		return
	}
	err := query.Delete(&sys.SysSystem{}).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Resp{
			Success: models.ERROR,
			Data:    err.Error(),
			Msg:     models.CustomError[models.NotOk],
		})
		models.FailWithDetailed(err, models.CustomError[models.NotOk], c)
	} else {
		filteredNamedPolicy, err := global.CasbinACLEnforcer.GetFilteredNamedPolicy("p", 0, system.AppId)
		if err != nil {
			models.FailWithDetailed(err, models.CustomError[models.NotOk], c)
		}
		if len(filteredNamedPolicy) > 0 {
			for _, p := range filteredNamedPolicy {
				if _, err := global.CasbinACLEnforcer.RemoveNamedPolicy("p", system.AppId, p[1], p[2]); err != nil {
					models.FailWithDetailed("删除API系统授权错误:"+system.AppId+" "+p[1]+" "+p[2]+" 错误："+err.Error(), models.CustomError[models.NotOk], c)
					return
				}
			}
		}
		models.OkResult(c)
	}
}

// @Summary [系统内部]获取系统所有权限
// @Id GetSystemPermById
// @Tags [系统内部]系统
// @version 1.0
// @Accept application/x-json-stream
// @Param	id		path 	string	true		"系统ID"
// @Success 200 object models.Resp 返回列表
// @Failure 400 object models.Resp 查询失败
// @Security ApiKeyAuth
// @Router /api/v1/system/perm/get/{id} [get]
func GetSystemPermById(c *gin.Context) {
	var system sys.SysSystem
	id := cast.ToInt(c.Param("id"))
	err := global.Mysql.Where("id = ?", id).First(&system).Error
	if err != nil {
		models.FailWithDetailed(err, models.CustomError[models.NotOk], c)
	} else {
		filteredNamedPolicy, err := global.CasbinACLEnforcer.GetFilteredNamedPolicy("p", 0, system.AppId)
		if err != nil {
			models.FailWithDetailed(err, models.CustomError[models.NotOk], c)
		}
		var system_perms []sys.SystemPermission
		for key, perm := range filteredNamedPolicy {
			system_perms = append(system_perms, sys.SystemPermission{
				ID:            key,
				AppId:         perm[0],
				AbsolutePath:  perm[1],
				AbsolutePath1: perm[2],
				AbsolutePath2: perm[3],
				HttpMethod:    perm[4],
				Eft:           perm[5],
			})
		}
		models.OkWithDataList(system_perms, cast.ToInt64(len(filteredNamedPolicy)), c)
	}
}

// @Summary [系统内部]创建系统权限
// @Id CreateSystemPerm
// @Tags [系统内部]系统
// @version 1.0
// @Accept application/x-json-stream
// @Param	id		path 	string	true "系统ID"
// @Param	body	body 	sys.SystemPermission	true "系统权限"
// @Success 201 object models.Resp 返回创建
// @Failure 500 object models.Resp 创建失败
// @Security ApiKeyAuth
// @Router /api/v1/system/perm/create/{id} [post]
func CreateSystemPerm(c *gin.Context) {
	var system sys.SysSystem
	var system_perms sys.SystemPermission
	id := cast.ToInt(c.Param("id"))
	err := global.Mysql.Where("id = ?", id).First(&system).Error
	if err != nil {
		models.FailWithDetailed(err, models.CustomError[models.NotOk], c)
		return
	}
	// 绑定参数
	err = c.ShouldBindJSON(&system_perms)
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
	if _, err := global.CasbinACLEnforcer.AddPolicy(system.AppId, system_perms.AbsolutePath, system_perms.AbsolutePath1, system_perms.AbsolutePath2, system_perms.HttpMethod, system_perms.Eft); err != nil {
		models.FailWithDetailed("授权错误:"+system.AppId+" "+system_perms.AbsolutePath+" "+system_perms.AbsolutePath1+" "+system_perms.AbsolutePath2+" "+system_perms.HttpMethod+" 错误："+err.Error(), models.CustomError[models.NotOk], c)
		return
	}
	models.OkResult(c)
}

// @Summary [系统内部]删除系统权限
// @Id DeleteSystemPermById
// @Tags [系统内部]系统
// @version 1.0
// @Accept application/x-json-stream
// @Param	id		path 	string	true		"系统ID"
// @Param	body	body 	sys.SystemPermission	 true "系统权限"
// @Success 204 object models.Resp 返回创建
// @Failure 500 object models.Resp 创建失败
// @Security ApiKeyAuth
// @Router /api/v1/system/perm/delete/{id} [delete]
func DeleteSystemPermById(c *gin.Context) {
	id := cast.ToInt(c.Param("id"))
	var system sys.SysSystem
	var system_perms sys.SystemPermission
	query := global.Mysql.Where("id = ?", id).First(&system)
	if query.Error != nil {
		models.FailWithDetailed("记录不存在", models.CustomError[models.NotOk], c)
		return
	}
	// 绑定参数
	err := c.ShouldBindJSON(&system_perms)
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
	if _, err := global.CasbinACLEnforcer.RemoveNamedPolicy("p", system.AppId, system_perms.AbsolutePath, system_perms.AbsolutePath1, system_perms.AbsolutePath2, system_perms.HttpMethod, system_perms.Eft); err != nil {
		models.FailWithDetailed("删除API系统授权错误:"+system.AppId+" "+system_perms.AbsolutePath1+" "+system_perms.AbsolutePath2+" "+system_perms.HttpMethod+" 错误："+err.Error(), models.CustomError[models.NotOk], c)
		return
	}
	models.OkResult(c)
}
