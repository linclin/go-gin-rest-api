package initialize

import (
	"database/sql"
	"fmt"
	"go-gin-rest-api/models/sys"
	"go-gin-rest-api/pkg/global"
	"time"

	sqlDriver "github.com/go-sql-driver/mysql" // mysql驱动
	loggable "github.com/linclin/gorm2-loggable"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// 初始化mysql数据库
func Mysql() {
	logLevel := gormlogger.Info
	if global.Conf.System.RunMode == "prd" {
		logLevel = gormlogger.Warn
	}
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?%s",
			global.Conf.Mysql.Username,
			global.Conf.Mysql.Password,
			global.Conf.Mysql.Host,
			global.Conf.Mysql.Port,
			global.Conf.Mysql.Database,
			global.Conf.Mysql.Query,
		), // DSN data source name
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{Logger: gormlogger.Default.LogMode(logLevel)})
	if err != nil {
		switch e := err.(type) {
		case *sqlDriver.MySQLError:
			// MySQL error unkonw database;
			// refer https://dev.mysql.com/doc/refman/5.6/en/error-messages-server.html
			if e.Number == 1049 {
				createsql := fmt.Sprintf("CREATE DATABASE `%s` CHARSET utf8mb4 COLLATE utf8mb4_general_ci;", global.Conf.Mysql.Database)

				dbForCreateDatabase, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=utf8", global.Conf.Mysql.Username, global.Conf.Mysql.Password, global.Conf.Mysql.Host, global.Conf.Mysql.Port)+"&loc=Asia%2FShanghai")
				if err != nil {
					global.Log.Error(fmt.Sprintf("数据库连接错误:%v", err))
					panic(fmt.Sprintf("数据库连接错误:%v", err))
				}
				defer dbForCreateDatabase.Close()
				info, err := dbForCreateDatabase.Exec(createsql)
				if err != nil {
					global.Log.Error(fmt.Sprintf("数据库创建错误:%v, %v", info, err))
					panic(fmt.Sprintf("数据库创建错误:%v, %v", info, err))
				} else {
					global.Log.Info("数据库" + global.Conf.Mysql.Database + "创建成功")
				}
			} else {
				global.Log.Error(fmt.Sprintf("数据库连接错误: %v", err))
				panic(fmt.Sprintf("数据库连接错误: %v", err))
			}
		default:
			global.Log.Error(fmt.Sprintf("初始化mysql异常: %v", err))
			panic(fmt.Sprintf("初始化mysql异常: %v", err))
		}
	}
	sqlDB, _ := db.DB()
	// SetMaxIdleCons 设置连接池中的最大闲置连接数。
	sqlDB.SetMaxIdleConns(10)
	// SetMaxOpenCons 设置数据库的最大连接数量。
	sqlDB.SetMaxOpenConns(500)
	// SetConnMaxLifetiment 设置连接的最大可复用时间。
	sqlDB.SetConnMaxLifetime(time.Hour)
	global.Mysql = db
	// 自动迁移表结构
	global.Mysql.AutoMigrate(&sys.SysSystem{})
	global.Mysql.AutoMigrate(&sys.SysRouter{})
	global.Mysql.AutoMigrate(&sys.SysRole{})
	global.Mysql.AutoMigrate(&sys.SysApiLog{})
	global.Mysql.AutoMigrate(&sys.SysReqApiLog{})
	global.Mysql.AutoMigrate(&sys.SysCronjobLog{})
	global.Mysql.AutoMigrate(&sys.SysLock{})
	global.Mysql.Table("sys_change_logs").AutoMigrate(&loggable.ChangeLog{})
	global.Log.Debug("初始化mysql完成")
	//初始化数据变更记录插件
	_, err = loggable.Register(db, "sys_change_logs", loggable.ComputeDiff())
	if err != nil {
		panic(err)
	}
}
