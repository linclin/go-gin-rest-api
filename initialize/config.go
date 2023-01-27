package initialize

import (
	"fmt"
	"go-gin-rest-api/pkg/global"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

const (
	configType  = "yml"
	configPath  = "./conf"
	devConfig   = "config.se.yml"
	stageConfig = "config.st.yml"
	prodConfig  = "config.prd.yml"
)

// 初始化配置文件
func InitConfig() {

	// 获取实例(可创建多实例读取多个配置文件, 这里不做演示)
	v := viper.New()
	// 读取当前go运行环境变量
	env := os.Getenv("RunMode")
	configName := devConfig
	if env == "st" {
		configName = stageConfig
	} else if env == "prd" {
		configName = prodConfig
	}
	v.SetConfigName(configName)
	v.SetConfigType(configType)
	v.AddConfigPath(configPath)
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Sprintf("初始化配置文件失败: %v", err))
	}
	// 转换为结构体
	if err := v.Unmarshal(&global.Conf); err != nil {
		panic(fmt.Sprintf("初始化配置文件失败: %v", err))
	}
	// 绑定环境变量
	v.AutomaticEnv()
	// 监听文件修改，热加载配置。因此不需要重启服务器，就能让配置生效。
	v.WatchConfig()
	// 监听文件修改回调函数
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Printf("配置文件:%s 发生变更:%s\n", e.Name, e.Op)
		// 转换为结构体
		if err := v.Unmarshal(&global.Conf); err != nil {
			panic(fmt.Sprintf("初始化配置文件失败: %v", err))
		}
	})
}
