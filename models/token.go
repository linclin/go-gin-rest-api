package models

import "time"

type ReqToken struct {
	AppId     string `json:"AppId" binding:"required"`     // AppId
	AppSecret string `json:"AppSecret" binding:"required"` // AppSecret
}

// token返回结构体
type Token struct {
	Success bool      `json:"success"` // 请求是否成功
	Token   string    `json:"token"`   // token
	Expires time.Time `json:"expires"` // 过期时间
}
