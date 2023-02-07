package cronjob

import (
	"fmt"
	"go-gin-rest-api/models/sys"
	"go-gin-rest-api/pkg/global"
	"runtime/debug"
	"time"

	loggable "github.com/linclin/gorm2-loggable"
)

// 清理超过一周的日志表数据
type CleanLog struct {
}

func (u CleanLog) Run() {
	startTime := time.Now()
	global.Log.Debug("cronjob定时任务:CleanLog开始执行")
	defer func() {
		if panicErr := recover(); panicErr != nil {
			global.Log.Error(fmt.Sprintf("cronjob定时任务:CleanLog执行失败: %v\n堆栈信息: %v", panicErr, string(debug.Stack())))
		}
		lock := sys.NewLock("CleanLog", 600)
		lock.DeleteLock()
	}()
	//获取任务锁
	lock := sys.NewLock("CleanLog", 600)
	if !lock.TryLock() {
		global.Log.Error("cronjob定时任务:CleanLog获取任务锁失败")
		return
	}
	defer lock.DeleteLock()
	//删除日志
	if err := global.Mysql.Where("StartTime < ? ", time.Now().AddDate(0, 0, -7)).Unscoped().Delete(sys.SysApiLog{}).Error; err != nil {
		global.Log.Error("cronjob定时任务:CleanLog删除SysApiLog失败")
	}
	if err := global.Mysql.Where("StartTime < ? ", time.Now().AddDate(0, 0, -7)).Unscoped().Delete(sys.SysReqApiLog{}).Error; err != nil {
		global.Log.Error("cronjob定时任务:CleanLog删除SysReqApiLog失败")
	}
	if err := global.Mysql.Where("StartTime < ? ", time.Now().AddDate(0, 0, -7)).Unscoped().Delete(sys.SysCronjobLog{}).Error; err != nil {
		global.Log.Error("cronjob定时任务:CleanLog删除SysCronjobLog失败")
	}
	if err := global.Mysql.Table("sys_change_logs").Where("created_at < ? ", time.Now().AddDate(0, 0, -7).Unix()).Unscoped().Delete(loggable.ChangeLog{}).Error; err != nil {
		global.Log.Error("cronjob定时任务:CleanLog删除ChangeLog失败")
	}
	//记录任务日志表
	endTime := time.Now()
	execTime := endTime.Sub(startTime).Seconds()
	go sys.AddSysCronjobLog("CleanLog", "", "success", "", startTime, endTime, execTime)
	return
}
