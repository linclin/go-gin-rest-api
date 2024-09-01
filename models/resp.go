package models

import (
	"net/http"

	"github.com/spf13/cast"

	"github.com/gin-gonic/gin"
)

// http请求响应体
type Resp struct {
	RequestId string      `json:"request_id"` // 请求ID
	Success   bool        `json:"success"`    // 请求是否成功
	Data      interface{} `json:"data"`       // 数据内容
	Msg       string      `json:"msg"`        // 消息提示
	Total     int64       `json:"total"`      // 数据总条数
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
	ERROR   = false
	SUCCESS = true
)

var EmptyArray = []interface{}{}

func Result(success bool, data interface{}, msg string, total int64, c *gin.Context) {
	requestId, _ := c.Get("RequestId")
	c.JSON(http.StatusOK, Resp{
		cast.ToString(requestId),
		success,
		data,
		msg,
		total,
	})
}

func OkResult(c *gin.Context) {
	Result(SUCCESS, map[string]interface{}{}, "操作成功", 0, c)
}

func OkWithMessage(message string, c *gin.Context) {
	Result(SUCCESS, map[string]interface{}{}, message, 0, c)
}

func OkWithData(data interface{}, c *gin.Context) {
	Result(SUCCESS, data, "操作成功", 0, c)
}
func OkWithDataList(data interface{}, total int64, c *gin.Context) {
	Result(SUCCESS, data, "操作成功", total, c)
}
func OkWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(SUCCESS, data, message, 0, c)
}

func FailResult(c *gin.Context) {
	Result(ERROR, map[string]interface{}{}, "操作失败", 0, c)
}

func FailWithMessage(message string, c *gin.Context) {
	Result(ERROR, map[string]interface{}{}, message, 0, c)
}

func FailWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(ERROR, data, message, 0, c)
}
