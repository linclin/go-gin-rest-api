package models

import (
	"net/http"

	"github.com/spf13/cast"

	"github.com/gin-gonic/gin"
)

// http请求响应体
type Resp struct {
	RequestId string      `json:"request_id"` // 请求ID
	Code      int         `json:"code"`       // 错误代码
	Data      interface{} `json:"data"`       // 数据内容
	Msg       string      `json:"msg"`        // 消息提示
}

// 自定义错误码与错误信息

const (
	Ok                  = 201
	NotOk               = 405
	Unauthorized        = 401
	Forbidden           = 403
	InternalServerError = 500
)

const (
	OkMsg                  = "操作成功"
	NotOkMsg               = "操作失败"
	UnauthorizedMsg        = "登录过期, 需要重新登录"
	LoginCheckErrorMsg     = "用户名或密码错误"
	ForbiddenMsg           = "无权访问该资源, 请联系网站管理员授权"
	InternalServerErrorMsg = "服务器内部错误"
)

var CustomError = map[int]string{
	Ok:                  OkMsg,
	NotOk:               NotOkMsg,
	Unauthorized:        UnauthorizedMsg,
	Forbidden:           ForbiddenMsg,
	InternalServerError: InternalServerErrorMsg,
}

const (
	ERROR   = 1
	SUCCESS = 0
)

var EmptyArray = []interface{}{}

func Result(code int, data interface{}, msg string, c *gin.Context) {
	requestId, _ := c.Get("RequestId")
	c.JSON(http.StatusOK, Resp{
		cast.ToString(requestId),
		code,
		data,
		msg,
	})
}

func OkResult(c *gin.Context) {
	Result(SUCCESS, map[string]interface{}{}, "操作成功", c)
}

func OkWithMessage(message string, c *gin.Context) {
	Result(SUCCESS, map[string]interface{}{}, message, c)
}

func OkWithData(data interface{}, c *gin.Context) {
	Result(SUCCESS, data, "操作成功", c)
}

func OkWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(SUCCESS, data, message, c)
}

func FailResult(c *gin.Context) {
	Result(ERROR, map[string]interface{}{}, "操作失败", c)
}

func FailWithMessage(message string, c *gin.Context) {
	Result(ERROR, map[string]interface{}{}, message, c)
}

func FailWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(ERROR, data, message, c)
}

// 分页封装
type PageInfo struct {
	Total  int64  `json:"total"`                // 数据总条数
	Offset int    `json:"offset" form:"offset"` // 当前页码
	Limit  int    `json:"limit" form:"limit"`   // 每页显示条数
	SortBy string `json:"sortby"`               // 排序字段

}

// 带分页数据封装
type PageData struct {
	PageInfo
	List interface{} `json:"list"` // 数据列表
}
