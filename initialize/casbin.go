package initialize

import (
	"fmt"
	"go-gin-rest-api/pkg/global"
	"go-gin-rest-api/pkg/utils"
	"time"

	"github.com/casbin/casbin/v2"
	defaultrolemanager "github.com/casbin/casbin/v2/rbac/default-role-manager"
	gormadapter "github.com/casbin/gorm-adapter/v3"
)

// 获取casbin策略管理器
func InitCasbin() {
	// 初始化数据库适配器
	casbinAdapter, err := gormadapter.NewAdapterByDBUseTableName(global.DB, "sys", "casbin_rule")
	if err != nil {
		global.Log.Error(fmt.Sprint("Casbin初始化错误", err.Error()))
	}
	// 读取配置文件
	CasbinACLEnforcer, err := casbin.NewSyncedEnforcer("conf/"+global.Conf.Casbin.ModelPath, casbinAdapter, true)
	if err != nil {
		global.Log.Error(fmt.Sprint("Casbin初始化错误", err.Error()))
		panic(err)
	}
	global.CasbinACLEnforcer = CasbinACLEnforcer
	global.CasbinACLEnforcer.SetRoleManager(defaultrolemanager.NewConditionalRoleManager(10))
	global.CasbinACLEnforcer.BuildRoleLinks()
	global.CasbinACLEnforcer.LoadPolicy()
	global.CasbinACLEnforcer.EnableAutoBuildRoleLinks(true)
	// 加载策略
	global.CasbinACLEnforcer.StartAutoLoadPolicy(time.Minute * time.Duration(1))
	//global.CasbinACLEnforcer.EnableEnforce(false)
	// 添加API系统权限
	global.CasbinACLEnforcer.AddPolicy("api-00000001", "/*", "*", "*", "(GET)|(POST)|(PUT)|(DELETE)|(OPTIONS)|(PATCH)", "allow")
	global.CasbinACLEnforcer.AddPolicy("api-00000002", "/*", "*", "*", "(GET)|(POST)|(PUT)|(DELETE)|(OPTIONS)|(PATCH)", "allow")
	// 添加前台普通用户组权限
	global.CasbinACLEnforcer.AddPolicy("group_user", "/*", "*", "*", "GET", "allow")
	// 添加前台操作用户组权限
	global.CasbinACLEnforcer.AddPolicy("group_operator", "/*", "*", "*", "GET", "allow")
	// 添加前台管理员组权限
	global.CasbinACLEnforcer.AddPolicy("group_admin", "/*", "*", "*", "(GET)|(POST)|(PUT)|(DELETE)|(OPTIONS)|(PATCH)", "allow")
	// global.CasbinACLEnforcer.AddRoleForUser("lc", "group_admin", time.Now().Format("2006-01-02 15:04:05"), time.Now().AddDate(100, 0, 0).Format("2006-01-02 15:04:05"))
	global.CasbinACLEnforcer.AddRoleForUser("lc", "group_admin", "2024-09-01 00:00:00", "2054-09-01 00:00:00")
	global.CasbinACLEnforcer.GetRoleManager().AddLink("lc", "group_admin")
	global.CasbinACLEnforcer.AddNamedLinkConditionFunc("g", "lc", "group_admin", utils.TimeMatchFunc)
}
