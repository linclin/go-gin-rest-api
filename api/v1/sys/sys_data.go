package sys

import (
	"go-gin-rest-api/models"
	"go-gin-rest-api/models/sys"
	"go-gin-rest-api/pkg/global"
	"time"

	"github.com/gin-gonic/gin"
)

// @Summary [系统内部]获取系统运营数据
// @Id GetSysData
// @Tags [系统内部]日志
// @version 1.0
// @Accept application/x-json-stream
// @Success 200 object models.Resp 返回列表
// @Failure 400 object models.Resp 查询失败
// @Security ApiKeyAuth
// @Router /api/v1/data/list [get]
func GetSysData(c *gin.Context) {
	data := sys.SysData{}
	query := global.Mysql
	err := query.Model(sys.SysSystem{}).Count(&data.SystemCount).Error
	if err != nil {
		models.FailWithDetailed("", err.Error(), c)
		return
	}
	err = query.Model(sys.SysRouter{}).Count(&data.RouterCount).Error
	if err != nil {
		models.FailWithDetailed("", err.Error(), c)
		return
	}
	err = query.Model(sys.SysApiLog{}).Count(&data.AllApiCount).Error
	if err != nil {
		models.FailWithDetailed("", err.Error(), c)
		return
	}
	err = query.Model(sys.SysApiLog{}).Where("StartTime >= ?", time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Now().Location())).Count(&data.ApiCount).Error
	if err != nil {
		models.FailWithDetailed("", err.Error(), c)
		return
	}
	err = query.Model(sys.SysApiLog{}).Select("DATE_FORMAT(StartTime, '%Y-%m-%d') as date, COUNT(*) as count").Where("StartTime >= ?", time.Now().AddDate(0, 0, -7)).Group("date").Scan(&data.WeekApiCount).Error
	if err != nil {
		models.FailWithDetailed("", err.Error(), c)
		return
	}
	err = query.Model(sys.SysApiLog{}).Select("ClientIP , COUNT(*) as count").Where("StartTime >= ?", time.Now().AddDate(0, 0, -7)).Group("ClientIP").Scan(&data.WeekClientApiCount).Error
	if err != nil {
		models.FailWithDetailed("", err.Error(), c)
		return
	}
	err = query.Model(sys.SysReqApiLog{}).Count(&data.AllReqApiCount).Error
	if err != nil {
		models.FailWithDetailed("", err.Error(), c)
		return
	}
	err = query.Model(sys.SysReqApiLog{}).Where("StartTime >= ?", time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Now().Location())).Count(&data.ReqApiCount).Error
	if err != nil {
		models.FailWithDetailed("", err.Error(), c)
		return
	}
	models.OkWithDataList(data, 0, c)
}
