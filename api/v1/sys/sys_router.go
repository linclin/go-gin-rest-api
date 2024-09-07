package sys

import (
	"go-gin-rest-api/models"
	"go-gin-rest-api/models/sys"
	"go-gin-rest-api/pkg/global"

	"github.com/a8m/rql"
	"github.com/gin-gonic/gin"
)

// @Summary [系统内部]获取系统路由
// @Id GetRouters
// @Tags [系统内部]路由
// @version 1.0
// @Accept application/x-json-stream
// @Param	body		body 	models.Req	true		"RQL查询json"
// @Success 200 object models.Resp 返回列表
// @Failure 400 object models.Resp 查询失败
// @Security ApiKeyAuth
// @Router /api/v1/router/list [post]
func GetRouters(c *gin.Context) {
	rqlQueryParser, err := rql.NewParser(rql.Config{
		Model:        sys.SysRouter{},
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
	list := make([]sys.SysRouter, 0)
	query := global.Mysql
	count := int64(0)
	err = query.Model(sys.SysRouter{}).Where(rqlParams.FilterExp, rqlParams.FilterArgs...).Count(&count).Error
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
