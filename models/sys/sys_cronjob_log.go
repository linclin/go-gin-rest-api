package sys

import (
	"fmt"
	"go-gin-rest-api/pkg/global"
	"time"

	"gorm.io/gorm"
)

// 任务日志表
type SysCronjobLog struct {
	gorm.Model
	CronMethod string    `gorm:"column:CronMethod;comment:任务名称" json:"CronMethod" rql:"filter,sort,column=CronMethod"` // 任务名称
	CronParam  string    `gorm:"column:CronParam;comment:任务参数" json:"CronParam"`                                       // 任务参数
	StartTime  time.Time `gorm:"column:StartTime;comment:开始时间" json:"StartTime"  rql:"filter,sort,column=StartTime"`   // 开始时间
	EndTime    time.Time `gorm:"column:EndTime;comment:结束时间" json:"EndTime" rql:"filter,sort,column=EndTime"`          // 结束时间
	ExecTime   float64   `gorm:"column:ExecTime;comment:执行时间(秒)" json:"ExecTime"`                                      // 执行时间(秒)
	Status     string    `gorm:"column:Status;comment:执行状态" json:"Status" rql:"filter,sort,column=StartTime"`          // 执行状态
	ErrMsg     string    `gorm:"column:ErrMsg;comment:错误信息" json:"ErrMsg"`                                             // 错误信息
}

func AddSysCronjobLog(cronMethod, cronParam, status, errMsg string, startTime, endTime time.Time, execTime float64) error {
	var cronjoblog = SysCronjobLog{
		CronMethod: cronMethod,
		CronParam:  cronParam,
		StartTime:  startTime,
		EndTime:    endTime,
		ExecTime:   execTime,
		Status:     status,
		ErrMsg:     errMsg,
	}
	err := global.DB.Create(&cronjoblog).Error
	if err != nil {
		global.Log.Error(fmt.Sprint("AddSysCronjobLog写入定时任务日志表失败", err.Error()))
		return err
	}
	return nil
}
