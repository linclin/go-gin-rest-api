package sys

import (
	"go-gin-rest-api/pkg/global"
	"time"

	"github.com/spf13/cast"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 系统请求API日志表
type SysReqApiLog struct {
	gorm.Model
	RequestId     string    `gorm:"column:RequestId;comment:请求ID" json:"RequestId" rql:"filter,sort,column=RequestId"`                            // 请求ID
	RequestMethod string    `gorm:"column:RequestMethod;comment:请求方法" json:"RequestMethod" rql:"filter,sort,column=RequestMethod"`                // HTTP方法
	RequestURI    string    `gorm:"column:RequestURI;index;comment:请求路径" json:"RequestURI" rql:"filter,sort,column=RequestURI"`                   // 路由地址
	RequestBody   string    `gorm:"column:RequestBody;type:longText;comment:请求体" json:"RequestBody" rql:"filter,sort,column=RequestBody"`         // 请求体
	StatusCode    int       `gorm:"column:StatusCode;index;comment:状态码" json:"StatusCode" rql:"filter,sort,column=StatusCode"`                    // 状态码
	RespBody      string    `gorm:"column:RespBody;type:longText;comment:返回体" json:"RespBody" rql:"filter,sort,column=RespBody"`                  // 返回体
	StartTime     time.Time `gorm:"column:StartTime;comment:访问时间" json:"StartTime" rql:"filter,sort,column=StartTime,layout=2006-01-02 15:04:05"` // 访问时间
	ExecTime      string    `gorm:"column:ExecTime;comment:结束时间" json:"ExecTime" rql:"filter,sort,column=ExecTime,layout=2006-01-02 15:04:05"`    // 执行时间
}

func AddReqApi(c *gin.Context, RequestMethod, RequestURI, RequestBody, RespBody, ExecTime string, StatusCode int, StartTime time.Time) (err error) {
	requestId, _ := c.Get("RequestId")
	reqapilog := SysReqApiLog{
		RequestId:     cast.ToString(requestId),
		RequestMethod: RequestMethod,
		RequestURI:    RequestURI,
		RequestBody:   RequestBody,
		StatusCode:    StatusCode,
		RespBody:      RespBody,
		StartTime:     StartTime,
		ExecTime:      ExecTime,
	}
	err = global.Mysql.Create(&reqapilog).Error
	if err != nil {
		global.Log.Error("AddReqApi 写入请求日志数据初始化失败", err.Error())
	}
	return
}
