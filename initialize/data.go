package initialize

import (
	"go-gin-rest-api/models/sys"
	"go-gin-rest-api/pkg/global"
)

// 初始化数据
func InitData() {
	go sys.InitSysSystem()
	go sys.InitSysRole()
	global.Log.Info("初始化表数据完成")
}
