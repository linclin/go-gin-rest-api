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
	zapLogger, _ := zap.NewStdLogAt(global.Logger, zap.DebugLevel)
	c := cron.New(cron.WithLocation(nyc), cron.WithSeconds(), cron.WithLogger(cron.VerbosePrintfLogger(zapLogger)))
	c.AddJob("@every 1m", cron.NewChain(cron.Recover(cron.DefaultLogger), cron.SkipIfStillRunning(cron.DefaultLogger)).Then(&cronjob.UserSync{}))
	c.Start()
	global.Log.Debug("初始化定时任务完成")
}
