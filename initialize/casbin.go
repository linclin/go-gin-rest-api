package initialize

import (
	"go-gin-rest-api/pkg/global"
	"time"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
)

// 获取casbin策略管理器
func InitCasbin() {
	// 初始化数据库适配器
	casbinAdapter, err := gormadapter.NewAdapterByDB(global.Mysql)
	if err != nil {
		global.Log.Error("Casbin初始化错误", err.Error())
	}
	// 读取配置文件
	CasbinACLEnforcer, err := casbin.NewSyncedEnforcer("conf/"+global.Conf.Casbin.ModelPath, casbinAdapter, true)
	if err != nil {
		global.Log.Error("Casbin初始化错误", err.Error())
	}
	global.CasbinACLEnforcer = CasbinACLEnforcer
	// 加载策略
	CasbinACLEnforcer.StartAutoLoadPolicy(time.Minute * time.Duration(1))
	CasbinACLEnforcer.EnableLog(false)
	// Load the policy from DB.
	//CasbinACLEnforcer.LoadPolicy()
	//CasbinACLEnforcer.EnableEnforce(false)
	// 添加API系统权限
	CasbinACLEnforcer.AddPolicy("2023012701", "/*", "(GET)|(POST)|(PUT)|(DELETE)|(OPTIONS)|(PATCH)")
	CasbinACLEnforcer.AddPolicy("2023012702", "/*", "(GET)|(POST)|(PUT)|(DELETE)|(OPTIONS)|(PATCH)")
	// 添加前台普通用户组权限
	CasbinACLEnforcer.AddPolicy("group_user", "/*", "GET")
	// 添加前台开发用户组权限
	CasbinACLEnforcer.AddPolicy("group_dev", "/*", "GET")
	// 添加前台运维操作用户组权限
	CasbinACLEnforcer.AddPolicy("group_operator", "/*", "GET")
	// 添加前台管理员组权限
	CasbinACLEnforcer.AddPolicy("group_admin", "/*", "(GET)|(POST)|(PUT)|(DELETE)|(OPTIONS)|(PATCH)")
	CasbinACLEnforcer.AddRoleForUser("lc", "group_admin")
}
