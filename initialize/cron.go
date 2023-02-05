package initialize

import (
	"go-gin-rest-api/cronjob"
	"go-gin-rest-api/pkg/global"
	"time"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

// 初始化定时任务
func Cron() {
	nyc, _ := time.LoadLocation("Asia/Shanghai")
	zapLogger, _ := zap.NewStdLogAt(global.Logger, zap.InfoLevel)
	c := cron.New(cron.WithLocation(nyc), cron.WithSeconds(), cron.WithLogger(cron.VerbosePrintfLogger(zapLogger)))
	//清理超过一周的日志表数据
	c.AddJob("@every 1d", cron.NewChain(cron.Recover(cron.VerbosePrintfLogger(zapLogger)), cron.SkipIfStillRunning(cron.VerbosePrintfLogger(zapLogger))).Then(&cronjob.CleanLog{}))
	c.Start()
	global.Log.Debug("初始化定时任务完成")
}
