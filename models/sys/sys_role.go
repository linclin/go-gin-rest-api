package sys

import (
	"fmt"
	"go-gin-rest-api/pkg/global"

	"gorm.io/gorm"
)

// 系统角色表
type SysRole struct {
	gorm.Model
	Name     string `gorm:"column:Name;uniqueIndex;comment:角色名称" json:"Name" binding:"required"   rql:"filter,sort,column=Name"`                   // 角色名称
	Keyword  string `gorm:"column:Keyword;comment:角色关键词" json:"Keyword" rql:"filter,sort,column=Keyword"`                                          // 角色关键词
	Desc     string `gorm:"column:Desc;comment:角色说明" json:"Desc" rql:"filter,sort,column=Desc"`                                                    // 角色说明
	Status   *bool  `gorm:"column:Status;index;type:tinyint(1);default:1;comment:角色状态(正常/禁用, 默认正常)" json:"Status" rql:"filter,sort,column=Status"` // 角色状态(正常/禁用, 默认正常)
	Operator string `gorm:"column:Operator;comment:操作人" json:"Operator"`
}
type RolePermission struct {
	Obj    string
	Obj1   string
	Obj2   string
	Action string
}

func InitSysRole() {
	roles := []SysRole{
		{
			Model: gorm.Model{
				ID: 1,
			},
			Name:    "admin",
			Keyword: "admin",
			Desc:    "超级管理员",
		},
		{
			Model: gorm.Model{
				ID: 2,
			},
			Name:    "operator",
			Keyword: "operator",
			Desc:    "管理员",
		},
		{
			Model: gorm.Model{
				ID: 3,
			},
			Name:    "dev",
			Keyword: "dev",
			Desc:    "开发用户",
		},
		{
			Model: gorm.Model{
				ID: 4,
			},
			Name:    "user",
			Keyword: "user",
			Desc:    "普通用户",
		},
	}
	for _, role := range roles {
		err := global.Mysql.Where(&role).FirstOrCreate(&role).Error
		if err != nil {
			global.Log.Error(fmt.Sprint("InitSysRole 数据初始化失败", err.Error()))
		}
	}
}
