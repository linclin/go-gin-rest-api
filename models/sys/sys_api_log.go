package sys

import (
	"time"

	"gorm.io/gorm"
)

// 系统API日志表
type SysApiLog struct {
	gorm.Model
	RequestId     string    `gorm:"column:RequestId;unique;index;comment:请求ID" json:"RequestId" rql:"filter,sort,column=RequestId"`               // 请求ID
	RequestMethod string    `gorm:"column:RequestMethod;comment:请求方法" json:"RequestMethod" rql:"filter,sort,column=RequestMethod"`                // HTTP方法
	RequestURI    string    `gorm:"column:RequestURI;index;comment:请求路径" json:"RequestURI" rql:"filter,sort,column=RequestURI"`                   // 路由地址
	RequestBody   string    `gorm:"column:RequestBody;type:longText;comment:请求体" json:"RequestBody" rql:"filter,sort,column=RequestBody"`         // 请求体
	StatusCode    int       `gorm:"column:StatusCode;index;comment:状态码" json:"StatusCode" rql:"filter,sort,column=StatusCode"`                    // 状态码
	RespBody      string    `gorm:"column:RespBody;type:longText;comment:返回体" json:"RespBody" rql:"filter,sort,column=RespBody"`                  // 返回体
	ClientIP      string    `gorm:"column:ClientIP;comment:访问IP" json:"ClientIP" rql:"filter,sort,column=ClientIP"`                               // 访问IP
	StartTime     time.Time `gorm:"column:StartTime;comment:访问时间" json:"StartTime" rql:"filter,sort,column=StartTime,layout=2006-01-02 15:04:05"` // 访问时间
	ExecTime      string    `gorm:"column:ExecTime;comment:结束时间" json:"ExecTime" rql:"filter,sort,column=ExecTime,layout=2006-01-02 15:04:05"`    // 执行时间
}
