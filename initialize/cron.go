package initialize

import (
	"go-gin-rest-api/cronjob"
	"go-gin-rest-api/pkg/global"
	"time"

	"github.com/robfig/cron/v3"
)

// 初始化定时任务
func Cron() {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	c := cron.New(cron.WithLocation(loc), cron.WithSeconds(), cron.WithLogger(cron.VerbosePrintfLogger(global.Logger)))
	//清理超过一周的日志表数据
	c.AddJob("0 0 1 * * *", cron.NewChain(cron.Recover(cron.VerbosePrintfLogger(global.Logger)), cron.SkipIfStillRunning(cron.VerbosePrintfLogger(global.Logger))).Then(&cronjob.CleanLog{}))
	c.Start()
	global.Log.Info("初始化定时任务完成")
}
