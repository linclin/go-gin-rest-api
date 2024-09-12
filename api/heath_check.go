package api

import (
	"go-gin-rest-api/models"
	"go-gin-rest-api/pkg/global"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary [系统内部]健康检查接口
// @Id HeathCheck
// @Tags [系统内部]路由
// @version 1.0
// @Accept application/x-json-stream
// @Success 200 object models.Resp 返回列表
// @Failure 500 object models.Resp 查询失败
// @Router /heatch_check [get]
func HeathCheck(c *gin.Context) {
	errStr := ""
	// MySQL连接检查
	db, _ := global.DB.DB()
	err := db.Ping()
	if err != nil {
		errStr += "健康检查失败 数据库连接错误：" + err.Error() + "\r\n"
	}
	if errStr != "" {
		c.JSON(http.StatusInternalServerError, models.Resp{
			Success: models.SUCCESS,
			Data:    errStr,
			Msg:     models.CustomError[models.NotOk],
		})
		return
	}
	c.JSON(http.StatusOK, models.Resp{
		Success: models.SUCCESS,
		Data:    "健康检查完成",
		Msg:     models.CustomError[models.Ok],
	})
}
