package handler

import (
	"fmt"
	"go-gin-rest-api/models"
	"go-gin-rest-api/pkg/global"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// APIHandler 通用API处理器
type APIHandler struct {
	db  *gorm.DB
	log interface{} // 使用 interface{} 以兼容不同的日志类型
}

// NewAPIHandler 创建新的API处理器
func NewAPIHandler() *APIHandler {
	return &APIHandler{
		db:  global.DB,
		log: global.Log,
	}
}

// HandleListQuery 通用列表查询处理函数
func (h *APIHandler) HandleListQuery(c *gin.Context, model interface{}, tableName ...string) {
	// 此函数需要根据具体类型实现，因为GORM需要具体类型的切片
	// 为不同模型类型创建特定的处理函数会更合适
	models.FailWithDetailed("", "通用查询函数需要针对特定模型实现", c)
	return
}

// HandleGetByID 通用根据ID获取记录处理函数
func (h *APIHandler) HandleGetByID(c *gin.Context, model interface{}, idField string) {
	id := c.Param("id")
	err := global.DB.Where(fmt.Sprintf("%s = ?", idField), id).First(model).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, models.Resp{
				Success: false,
				Msg:     "记录不存在",
				Data:    "",
			})
		} else {
			models.FailWithDetailed("", err.Error(), c)
		}
	} else {
		models.OkWithData(model, c)
	}
}
