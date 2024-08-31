package sys

import (
	"fmt"
	"go-gin-rest-api/pkg/global"
	"strings"
	"time"

	"gorm.io/gorm"
)

// 任务锁
type SysLock struct {
	gorm.Model
	LockMethod string `gorm:"column:LockMethod;uniqueIndex;comment:任务名称" json:"LockMethod"` // 任务名称
	ExpireTime int64  `gorm:"column:ExpireTime;comment:过期时间" json:"ExpireTime"`             // 过期时间
}

func NewLock(lockMethod string, expireTime int64) *SysLock {
	return &SysLock{
		LockMethod: lockMethod,
		ExpireTime: expireTime,
	}
}

func (lock *SysLock) TryLock() bool {
	if err := lock.deleteExpiredLock(); err != nil {
		global.Log.Error(fmt.Sprint("清理过期任务锁失败", lock.LockMethod, err.Error()))
		return false
	}
	var newlock = SysLock{
		LockMethod: lock.LockMethod,
		ExpireTime: time.Now().Unix() + lock.ExpireTime,
	}
	newlock.DeletedAt = gorm.DeletedAt{
		Time:  time.Now(),
		Valid: true,
	}
	err := global.Mysql.Create(&newlock).Error
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return false
		}
		global.Log.Error(fmt.Sprint("获取任务锁失败", err.Error()))
		return false
	}
	return true
}

func (lock *SysLock) deleteExpiredLock() error {
	var now = time.Now().Unix()
	return global.Mysql.Where("LockMethod = ? AND ExpireTime < ?", lock.LockMethod, now).Unscoped().Delete(SysLock{}).Error
}

func (lock *SysLock) DeleteLock() error {
	return global.Mysql.Where("LockMethod = ? ", lock.LockMethod).Unscoped().Delete(SysLock{}).Error
}
