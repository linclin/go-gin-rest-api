package initialize

import (
	"fmt"

	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/config"
	"github.com/alibaba/sentinel-golang/core/flow"

	"go-gin-rest-api/pkg/global"
)

func InitSentinel() {
	//初始化 sentinel
	conf := config.NewDefaultConfig()
	conf.Sentinel.App.Name = global.Conf.System.AppName
	conf.Sentinel.Log.Dir = "./logs/sentinel"
	err := sentinel.InitWithConfig(conf)
	if err != nil {
		global.Log.Error("初始化Sentinel配置出错", err)
	}
	if _, err := flow.LoadRules([]*flow.Rule{
		{
			Resource:               "POST:/api/v1/base/auth",
			Threshold:              10,
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
		},
	}); err != nil {
		global.Log.Error(fmt.Sprintf("初始化Sentinel流控规则出错: %+v", err))
		return
	}
}
