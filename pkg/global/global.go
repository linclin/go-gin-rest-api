package global

import (
	"log"
	"log/slog"

	"github.com/casbin/casbin/v2"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

var (
	// 系统配置
	Conf Configuration
	// slog日志
	Logger *log.Logger
	Log    *slog.Logger
	// mysql实例
	Mysql *gorm.DB
	// Casbin实例
	CasbinACLEnforcer *casbin.SyncedEnforcer
	// validation.v10校验器
	Validate *validator.Validate
	// validation.v10相关翻译器
	Translator ut.Translator
)
