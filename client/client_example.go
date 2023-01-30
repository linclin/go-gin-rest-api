package main

import (
	"fmt"
	"log"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/patrickmn/go-cache"
)

type ReqToken struct {
	AppId     string `json:"AppId"`     // AppId
	AppSecret string `json:"AppSecret"` // AppSecret
}
type TokenResp struct {
	Token   string    `json:"token"`   // token
	Expires time.Time `json:"expires"` // 过期时间
}

func main() {
	token := ""
	client := resty.New()
	gocache := cache.New(5*time.Minute, 10*time.Minute)
	tokencache, found := gocache.Get("token")
	if found {
		fmt.Println("Cache Token:", tokencache.(TokenResp).Token)
		token = tokencache.(TokenResp).Token
	} else {
		// 1.请求token
		tokenresp := &TokenResp{}
		resp, err := client.R().
			SetHeader("accept", "application/json").
			SetHeader("Content-Type", "application/json").
			SetBody(ReqToken{AppId: "2023012801", AppSecret: "fa2e25cb060c8d748fd16ac5210581f41"}).
			SetResult(tokenresp).
			Post("http://127.0.0.1:8080/api/v1/base/auth")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("token Response Info:")
		fmt.Println("Status Code:", resp.StatusCode())
		fmt.Println("Token:", tokenresp.Token)
		fmt.Println("Token Expires:", tokenresp.Expires)
		token = tokenresp.Token
		//缓存时间小于120分钟
		gocache.Set("token", tokenresp, 100*time.Minute)
	}
	if token != "" {
		// 2.请求token
		resp, err := client.R().
			SetHeader("accept", "application/json").
			SetHeader("Content-Type", "application/json").
			SetAuthToken(token).
			SetBody(`{
				"filter": {
				  "RequestId": "84692443-987c-4df1-b91c-606fcff6b556"
				},
				"limit": 10,
				"offset": 0,
				"sort": [
				  "-StartTime"
				]
			  }`).
			Post("http://127.0.0.1:8080/api/v1/apilog/list")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("sysapilog Response Info:")
		fmt.Println("Status Code:", resp.StatusCode())
		fmt.Println("sysapilog Resp:", string(resp.Body()))
	}

}
