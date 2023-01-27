package api

import (
	"fmt"
	"go-gin-rest-api/models"
	"go-gin-rest-api/pkg/global"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
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
	db, _ := global.Mysql.DB()
	err := db.Ping()
	if err != nil {
		errStr += "健康检查失败 数据库连接错误：" + err.Error() + "\r\n"
	}
	if errStr != "" {
		c.JSON(http.StatusInternalServerError, models.Resp{
			Code: http.StatusInternalServerError,
			Data: errStr,
			Msg:  models.CustomError[models.NotOk],
		})
		return
	}
	c.JSON(http.StatusOK, models.Resp{
		Code: http.StatusOK,
		Data: "健康检查完成",
		Msg:  models.CustomError[models.Ok],
	})
}

func NetworkCheck(c *gin.Context) {
	client := resty.New()
	resp, err := client.R().Get("https://baidu.com")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Response Info:")
	fmt.Println("Status Code:", resp.StatusCode())
	fmt.Println("Status:", resp.Status())
	fmt.Println("Proto:", resp.Proto())
	fmt.Println("Time:", resp.Time())
	fmt.Println("Received At:", resp.ReceivedAt())
	fmt.Println("Size:", resp.Size())
	fmt.Println("Headers:")
	for key, value := range resp.Header() {
		fmt.Println(key, "=", value)
	}
	fmt.Println("Cookies:")
	for i, cookie := range resp.Cookies() {
		fmt.Printf("cookie%d: name:%s value:%s\n", i, cookie.Name, cookie.Value)
	}
}
