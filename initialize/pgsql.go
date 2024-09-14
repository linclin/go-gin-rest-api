package initialize

import (
	"context"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"go-gin-rest-api/models/sys"
	"go-gin-rest-api/pkg/global"
)

// 初始化pgsql数据库
func Pgsql() {
	logLevel := logger.Info
	if global.Conf.System.RunMode == "prd" {
		logLevel = logger.Warn
	}
	if global.Conf.Pgsql == nil {
		return
	}
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		global.Conf.Pgsql.Host,
		global.Conf.Pgsql.Port,
		global.Conf.Pgsql.Username,
		global.Conf.Pgsql.Password,
		global.Conf.Pgsql.Database,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logLevel)})
	if err != nil {
		global.Log.Error(fmt.Sprintf("数据库连接错误: %v", err))
		panic(fmt.Sprintf("数据库连接错误: %v", err))
	}

	sqlDB, err := db.DB()
	if err != nil {
		global.Log.Error(fmt.Sprintf("获取数据库连接错误: %v", err))
		panic(fmt.Sprintf("获取数据库连接错误: %v", err))
	}

	// SetMaxIdleCons 设置连接池中的最大闲置连接数。
	sqlDB.SetMaxIdleConns(10)
	// SetMaxOpenCons 设置数据库的最大连接数量。
	sqlDB.SetMaxOpenConns(500)
	// SetConnMaxLifetime 设置连接的最大可复用时间。
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 检查数据库是否存在，如果不存在则创建
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	pgdb, err := sqlDB.Conn(ctx)
	if err != nil {
		global.Log.Error(fmt.Sprintf("连接数据库错误: %v", err))
		panic(fmt.Sprintf("连接数据库错误: %v", err))
	}
	defer pgdb.Close()
	err = createDatabaseIfNotExists(ctx, db, global.Conf.Pgsql.Database)

	global.DB = db
	// 自动迁移表结构
	global.DB.AutoMigrate(&sys.SysSystem{})
	global.DB.AutoMigrate(&sys.SysRouter{})
	global.DB.AutoMigrate(&sys.SysRole{})
	global.DB.AutoMigrate(&sys.SysApiLog{})
	global.DB.AutoMigrate(&sys.SysReqApiLog{})
	global.DB.AutoMigrate(&sys.SysCronjobLog{})
	global.DB.AutoMigrate(&sys.SysLock{})
	global.Log.Info("初始化pgsql完成")
}

func createDatabaseIfNotExists(ctx context.Context, db *gorm.DB, dbName string) error {
	// 检查数据库是否存在
	var count int64
	err := db.Raw("SELECT COUNT(*) FROM pg_database WHERE datname = ?", dbName).Count(&count).Error
	if err != nil {
		return err
	}

	// 如果数据库不存在，则创建数据库
	if count == 0 {
		_, err = db.Statement.ConnPool.ExecContext(ctx, "CREATE DATABASE  "+dbName)
		if err != nil {
			return err
		}
	}

	return nil
}
