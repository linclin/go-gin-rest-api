package sys

import (
	"go-gin-rest-api/models"
	"go-gin-rest-api/models/sys"
	"go-gin-rest-api/pkg/global"

	"github.com/a8m/rql"
	"github.com/gin-gonic/gin"
)

// @Summary [系统内部]获取定时任务日志
// @Id GetCronjobLog
// @Tags [系统内部]日志
// @version 1.0
// @Accept application/x-json-stream
// @Param	body body  models.Req true "RQL查询json"
// @Success 200 object models.Resp 返回列表
// @Failure 400 object models.Resp 查询失败
// @Security ApiKeyAuth
// @Router /api/v1/cronjoblog/list [post]
func GetCronjobLog(c *gin.Context) {
	rqlQueryParser, err := rql.NewParser(rql.Config{
		Model:        sys.SysCronjobLog{},
		DefaultLimit: -1,
	})
	if err != nil {
		models.FailWithDetailed("", err.Error(), c)
		return
	}
	// 绑定参数
	var rqlQuery rql.Query
	err = c.ShouldBindJSON(&rqlQuery)
	if err != nil {
		models.FailWithDetailed("", err.Error(), c)
		return
	}
	rqlParams, err := rqlQueryParser.ParseQuery(&rqlQuery)
	if err != nil {
		models.FailWithDetailed("", err.Error(), c)
		return
	}
	if rqlParams.Sort == "" {
		rqlParams.Sort = "id desc"
	}
	list := make([]sys.SysReqApiLog, 0)
	query := global.DB
	count := int64(0)
	err = query.Model(sys.SysReqApiLog{}).Where(rqlParams.FilterExp, rqlParams.FilterArgs...).Count(&count).Error
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
