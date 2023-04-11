package main

import (
	"context"
	"fmt"
	_ "go-gin-rest-api/docs"
	"go-gin-rest-api/initialize"
	"go-gin-rest-api/pkg/global"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"runtime/debug"
	"syscall"
	"time"
)

var (
	// 初始化为 unknown，如果编译时没有传入这些值，则为 unknown
	GitBranch      = "unknown"
	GitRevision    = "unknown"
	GitCommitLog   = "unknown"
	BuildTime      = "unknown"
	BuildGoVersion = "unknown"
)

func init() {
	//输出程序分支 commit golang版本  构建时间
	fmt.Fprint(os.Stdout, buildInfo())
	// 初始化配置
	initialize.InitConfig()
	// 初始化日志
	initialize.Logger()
	// 初始化数据库
	initialize.Mysql()
	// 初始化Sentinel流控规则
	initialize.InitSentinel()
	// 初始校验器
	initialize.Validate("zh")
	// 初始化jwt的rsa证书
	initialize.InitRSA()
	if global.Conf.System.InitData {
		// 初始化数据
		initialize.InitData()
	}
	// 初始化casbin策略管理器
	initialize.InitCasbin()
	// 初始化定时任务
	initialize.Cron()
}

// @title go-gin-rest-api
// @version 1.0.0
// @description go-gin-rest-api Golang后台api开发脚手架
// @termsOfService https://github.com/linclin
// @contact.name LC
// @contact.url https://github.com/linclin
// @contact.email 13579443@qq.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	defer func() {
		if err := recover(); err != nil {
			// 将异常写入日志
			global.Log.Error(fmt.Sprintf("项目启动失败: %v\n堆栈信息: %v", err, string(debug.Stack())))
		}
	}()
	// 初始化路由
	r := initialize.Routers()
	host := "0.0.0.0"
	port := global.Conf.System.Port
	// 服务启动及优雅关闭
	// 参考地址https://github.com/gin-gonic/examples/blob/master/graceful-shutdown/graceful-shutdown/server.go
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", host, port),
		Handler: r,
	}
	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			global.Log.Error("listen error: ", err)
		}
	}()
	global.Log.Info(fmt.Sprintf("Server is running at %s:%d/%s", host, port, global.Conf.System.UrlPathPrefix))
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	global.Log.Info("Shutting down server...")
	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		global.Log.Error("Server forced to shutdown: ", err)
	}
	global.Log.Info("Server exiting")
}

// 返回构建信息 多行格式
func buildInfo() string {
	return fmt.Sprintf("GitBranch=%s\nGitRevision=%s\nGitCommitLog=%s\nBuildTime=%s\nGoVersion=%s\nruntime=%s/%s\n",
		GitBranch, GitRevision, GitCommitLog, BuildTime, BuildGoVersion, runtime.GOOS, runtime.GOARCH)
}
