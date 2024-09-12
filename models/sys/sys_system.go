package sys

import (
	"fmt"
	"go-gin-rest-api/pkg/global"

	loggable "github.com/linclin/gorm2-loggable"
	"gorm.io/gorm"
)

// 系统表
type SysSystem struct {
	gorm.Model
	AppId      string `gorm:"column:AppId;unique;comment:AppId" json:"AppId" binding:"required"  rql:"filter,sort,column=AppId"` // AppId
	AppSecret  string `gorm:"column:AppSecret;comment:AppSecret" json:"AppSecret" binding:"required" `                           // AppSecret
	SystemName string `gorm:"column:SystemName;comment:系统名称" json:"SystemName" rql:"filter,sort,column=SystemName"`              // 系统名称
	IP         string `gorm:"column:IP;comment:系统来源IP" json:"IP" rql:"filter,sort,column=IP"`                                    // 系统来源IP
	Operator   string `gorm:"column:Operator;comment:操作人" json:"Operator" rql:"filter,sort,column=Operator"`                     // 操作人
	loggable.LoggableModel
}

func (system SysSystem) Meta() interface{} {
	return struct {
		CreatedBy string
	}{
		CreatedBy: system.Operator,
	}
}

// 权限
type SystemPermission struct {
	ID            int
	AppId         string `validate:"required"` // AppId
	AbsolutePath  string `validate:"required"` // 路由地址
	AbsolutePath1 string `validate:"required"` // 路由地址
	AbsolutePath2 string `validate:"required"` // 路由地址
	HttpMethod    string `validate:"required"` // HTTP方法
	Eft           string `validate:"required"` // 动作
}

func InitSysSystem() {
	systems := []SysSystem{
		{
			Model: gorm.Model{
				ID: 1,
			},
			AppId:      "api-00000001",
			AppSecret:  "fa2e25cb060c8d748fd16ac5210581f41",
			SystemName: "api",
			IP:         "",
			Operator:   "lc",
		},
		{
			Model: gorm.Model{
				ID: 2,
			},
			AppId:      "api-00000002",
			AppSecret:  "61c94399f47c485334b48f8f340bc07b2",
			SystemName: "UI",
			IP:         "",
			Operator:   "lc",
		},
	}
	for _, system := range systems {
		err := global.DB.Where(&system).FirstOrCreate(&system).Error
		if err != nil {
			global.Log.Error(fmt.Sprint("InitSysSystem 数据初始化失败", err.Error()))
			continue
		}
	}
	return
}
