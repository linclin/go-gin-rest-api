package global

import (
	"github.com/casbin/casbin/v2"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10" 
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	// 系统配置
	Conf Configuration
	// zap日志
	Log    *zap.SugaredLogger
	Logger *zap.Logger
	// mysql实例
	Mysql *gorm.DB
	// Casbin实例
	CasbinACLEnforcer *casbin.SyncedEnforcer
	// validation.v10校验器
	Validate *validator.Validate
	// validation.v10相关翻译器
	Translator ut.Translator 
)
