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
	global.Log.Info("AccessLogWriter WriteString", string(p))
	if n, err := w.body.Write(p); err != nil {
		global.Log.Info("AccessLogWriter Write", err.Error())
		return n, err
	}
	return w.ResponseWriter.Write(p)
}

func (w AccessLogWriter) WriteString(p string) (int, error) {
	global.Log.Info("AccessLogWriter WriteString", p)
	if n, err := w.body.WriteString(p); err != nil {
		global.Log.Info("AccessLogWriter WriteString", err.Error())
		return n, err
	}
	return w.ResponseWriter.WriteString(p)
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
	requestId := requestid.Get(c)
	c.Set("RequestId", requestId)
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
	respBody := bodyWriter.body.String()
	// 请求IP
	clientIP := c.ClientIP()
	if reqMethod != "OPTIONS" {
		sysApiLog := sys.SysApiLog{
			RequestId:     requestId,
			RequestMethod: reqMethod,
			RequestURI:    reqUri,
			RequestBody:   reqBody,
			StatusCode:    statusCode,
			RespBody:      respBody,
			ClientIP:      clientIP,
			StartTime:     startTime,
			ExecTime:      execTime.String(),
		}
		err = global.Mysql.Create(&sysApiLog).Error
		if err != nil {
			global.Log.Info("接口日志存库错误", err, requestId, reqMethod, reqUri, reqBody, respBody, statusCode, execTime.String(), clientIP)
		}
	}
}
