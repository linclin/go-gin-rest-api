package cronjob

import "go-gin-rest-api/pkg/global"

// 用户同步任务
type UserSync struct {
}

func (u UserSync) Run() {
	global.Log.Debug("UserSync tick every 1 min")
}
