package sys

import (
	"gorm.io/gorm"
)

// 系统路由表  鉴权用
type SysRouter struct {
	gorm.Model
	Name         string `gorm:"column:Name;comment:接口名称" json:"Name" rql:"filter,sort,column=Name"`                                                                     // 接口名称
	Group        string `gorm:"column:Group;comment:路由分组" json:"Group" rql:"filter,sort,column=Group"`                                                                  // 路由分组
	HttpMethod   string `gorm:"column:HttpMethod;index:idx_sysrouter_httpmethod_absolutepath;comment:HTTP方法" json:"HttpMethod" rql:"filter,sort,column=HttpMethod"`     // HTTP方法
	AbsolutePath string `gorm:"column:AbsolutePath;index:idx_sysrouter_httpmethod_absolutepath;comment:路由地址" json:"AbsolutePath" rql:"filter,sort,column=AbsolutePath"` // 路由地址
	HandlerName  string `gorm:"column:HandlerName;comment:控制器" json:"HandlerName" rql:"filter,sort,column=HandlerName"`                                                 // 控制器
}
