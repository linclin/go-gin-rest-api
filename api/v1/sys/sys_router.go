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
		Model: sys.SysRouter{},
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

	list := make([]sys.SysRouter, 0)
	query := global.Mysql
	query = query.Where(rqlParams.FilterExp, rqlParams.FilterArgs...)
	count := int64(0)
	err = query.Limit(rqlParams.Limit).Offset(rqlParams.Offset).Order(rqlParams.Sort).Find(&list).Count(&count).Error
	if err != nil {
		models.FailWithDetailed(err, models.CustomError[models.NotOk], c)
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
