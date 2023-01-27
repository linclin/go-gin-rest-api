package middleware

import (
	"bytes"
	"go-gin-rest-api/models/sys"
	"go-gin-rest-api/pkg/global"
	"io/ioutil"
	"time"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
)

type AccessLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w AccessLogWriter) Write(p []byte) (int, error) {
	if n, err := w.body.Write(p); err != nil {
		return n, err
	}
	return w.ResponseWriter.Write(p)
}

// 访问日志
func AccessLog(c *gin.Context) {
	bodyWriter := &AccessLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
	c.Writer = bodyWriter
	var bodyBytes []byte // 我们需要的body内容
	// 从原有Request.Body读取
	bodyBytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		global.Log.Info("请求体获取错误", err)
	}
	// 新建缓冲区并替换原有Request.body
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	// 开始时间
	startTime := time.Now()
	// 处理请求
	c.Next()
	// 结束时间
	endTime := time.Now()
	// 执行时间
	execTime := endTime.Sub(startTime)
	// 请求方式
	reqMethod := c.Request.Method
	// 请求路由
	reqUri := c.Request.RequestURI
	// 请求体
	reqBody := string(bodyBytes)
	// 状态码
	statusCode := c.Writer.Status()
	// 返回体
	RespBody := bodyWriter.body.String()
	// 请求IP
	clientIP := c.ClientIP()
	if reqMethod != "OPTIONS" {
		sysApiLog := sys.SysApiLog{
			RequestId:     requestid.Get(c),
			RequestMethod: reqMethod,
			RequestURI:    reqUri,
			RequestBody:   reqBody,
			StatusCode:    statusCode,
			RespBody:      RespBody,
			ClientIP:      clientIP,
			StartTime:     startTime,
			ExecTime:      execTime.String(),
		}
		err = global.Mysql.Create(&sysApiLog).Error
		if err != nil {
			global.Log.Info("接口日志存库错误", err, requestid.Get(c), reqMethod, reqUri, reqBody, RespBody, statusCode, execTime.String(), clientIP)
		}
	}
}
