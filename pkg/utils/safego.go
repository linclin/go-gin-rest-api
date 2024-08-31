package utils

import (
	"fmt"
	"go-gin-rest-api/pkg/global"
	"runtime/debug"
)

func SafeGo(f func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				global.Log.Error(fmt.Sprintf("运行panic异常: %v\n堆栈信息: %v", err, string(debug.Stack())))
			}
		}()
		f()
	}()
}
