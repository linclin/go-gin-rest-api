package sys

import (
	"go-gin-rest-api/models"
	"go-gin-rest-api/models/sys"
	"go-gin-rest-api/pkg/global"

	"github.com/a8m/rql"
	"github.com/gin-gonic/gin"
)

// @Summary [系统内部]获取服务接口日志
// @Id GetApiLog
// @Tags [系统内部]日志
// @version 1.0
// @Accept application/x-json-stream
// @Param	body body  models.Req true "RQL查询json"
// @Success 200 object models.Resp 返回列表
// @Failure 400 object models.Resp 查询失败
// @Security ApiKeyAuth
// @Router /api/v1/apilog/list [post]
func GetApiLog(c *gin.Context) {
	rqlQueryParser, err := rql.NewParser(rql.Config{
		Model: sys.SysApiLog{},
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

	list := make([]sys.SysApiLog, 0)
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

// @Summary [系统内部]根据ID获取接口日志
// @Id GetApiLogById
// @Tags [系统内部]日志
// @version 1.0
// @Accept application/x-json-stream
// @Param	requestid path string true "RequestId"
// @Success 200 object models.Resp 返回列表
// @Failure 400 object models.Resp 查询失败
// @Security ApiKeyAuth
// @Router /api/v1/apilog/get/{requestid} [get]
func GetApiLogById(c *gin.Context) {
	var apilog sys.SysApiLog
	requestid := c.Param("requestid")
	err := global.Mysql.Where("RequestId = ?", requestid).First(&apilog).Error
	if err != nil {
		models.FailWithDetailed(err, models.CustomError[models.NotOk], c)
	} else {
		models.OkWithData(apilog, c)
	}
}
