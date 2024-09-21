package initialize

import (
	"go-gin-rest-api/pkg/global"
	"os"

	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
)

// 初始化casdoor客户端
func InitCasdoor() {
	certificate, err := os.ReadFile("./conf/rsa/" + global.Conf.Casdoor.CertificatePath)
	if err != nil {
		return
		//panic(fmt.Sprintf("初始化casdoor客户端失败: %v", err))
	}
	global.Conf.Casdoor.Certificate = string(certificate)
	casdoorsdk.InitConfig(
		global.Conf.Casdoor.Endpoint,
		global.Conf.Casdoor.ClientID,
		global.Conf.Casdoor.ClientSecret,
		global.Conf.Casdoor.Certificate,
		global.Conf.Casdoor.Organization,
		global.Conf.Casdoor.Application,
	)
	global.Log.Info("初始化casdoor客户端完成")
}
