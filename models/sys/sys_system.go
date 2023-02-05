package sys

import (
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
	return nil
}

// 权限
type RbacRolePerm struct {
	HttpMethod   string `validate:"required"` // HTTP方法
	AbsolutePath string `validate:"required"` // 路由地址
}

func InitSysSystem() {
	systems := []SysSystem{
		{
			Model: gorm.Model{
				ID: 1,
			},
			AppId:      "2023012801",
			AppSecret:  "fa2e25cb060c8d748fd16ac5210581f41",
			SystemName: "api",
			IP:         "",
			Operator:   "lc",
		},
		{
			Model: gorm.Model{
				ID: 2,
			},
			AppId:      "2023012802",
			AppSecret:  "61c94399f47c485334b48f8f340bc07b2",
			SystemName: "UI",
			IP:         "",
			Operator:   "lc",
		},
	}
	for _, system := range systems {
		err := global.Mysql.Where(&system).FirstOrCreate(&system).Error
		if err != nil {
			global.Log.Error("InitSysSystem 数据初始化失败", err.Error())
		}
	}
}
