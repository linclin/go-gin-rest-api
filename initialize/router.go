package initialize

import (
	"fmt"
	"go-gin-rest-api/api"
	"go-gin-rest-api/middleware"
	"go-gin-rest-api/models/sys"
	"go-gin-rest-api/pkg/global"
	sysRouter "go-gin-rest-api/router/sys"

	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/pprof"
	"github.com/gin-contrib/requestid"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	ginprometheus "github.com/zsais/go-gin-prometheus"
)

// 初始化总路由
func Routers() *gin.Engine {
	// 初始化路由接口到数据库表中
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		//global.Log.Debug("router %v %v %v %v\n", httpMethod, absolutePath, handlerName, nuHandlers)
		go func(httpMethod, absolutePath, handlerName string) {
			sys_router := sys.SysRouter{
				HttpMethod:   httpMethod,
				AbsolutePath: absolutePath,
				HandlerName:  handlerName,
			}
			err := global.Mysql.Where(&sys_router).FirstOrCreate(&sys_router).Error
			if err != nil {
				global.Log.Error("SysRouter 数据初始化失败", err.Error())
			}
		}(httpMethod, absolutePath, handlerName)
	}
	if global.Conf.System.RunMode == "prd" {
		gin.SetMode(gin.ReleaseMode)
	}
	// 创建带有默认中间件的路由:
	// 日志与恢复中间件
	// r := gin.Default()
	// 创建不带中间件的路由:
	r := gin.New()
	// 初始化Trace中间件
	r.Use(requestid.New())
	// 添加访问记录
	r.Use(middleware.AccessLog)
	// 添加全局异常处理中间件
	r.Use(middleware.Exception)
	// GZip压缩插件
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	// 添加跨域中间件, 让请求支持跨域
	r.Use(cors.Default())
	// zap日志记录插件
	r.Use(ginzap.Ginzap(global.Logger, time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(global.Logger, true))
	// 添加速率访问中间件
	r.Use(middleware.RateLimiter())
	// 初始化pprof
	pprof.Register(r)
	// 初始化jwt auth中间件
	authMiddleware, err := middleware.InitAuth()
	if err != nil {
		panic(fmt.Sprintf("初始化JWT auth中间件失败: %v", err))
	}
	global.Log.Debug("初始化JWT auth中间件完成")
	// 初始化Prometheus中间件
	prome := ginprometheus.NewPrometheus("gin")
	prome.Use(r)
	global.Log.Debug("初始化Prometheus中间件完成")
	if global.Conf.System.RunMode != "prd" {
		// 初始化Swagger
		url := ginSwagger.URL(global.Conf.System.BaseApi + "/swagger/doc.json")
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
		global.Log.Debug("初始化Swagger完成")
	}
	// 初始化健康检查接口
	r.GET("/heatch_check", api.HeathCheck)
	global.Log.Debug("初始化健康检查接口完成")
	// 初始化API路由
	apiGroup := r.Group(global.Conf.System.UrlPathPrefix)
	// 方便统一添加路由前缀
	v1Group := apiGroup.Group("v1")
	sysRouter.InitPublicRouter(v1Group)                     // 注册公共路由
	sysRouter.InitBaseRouter(v1Group, authMiddleware)       // 注册基础路由, 不会鉴权
	sysRouter.InitRoleRouter(v1Group, authMiddleware)       // 注册角色路由
	sysRouter.InitSystemRouter(v1Group, authMiddleware)     // 注册系统路由
	sysRouter.InitRouterRouter(v1Group, authMiddleware)     // 注册系统路由路由
	sysRouter.InitApiLogRouter(v1Group, authMiddleware)     // 注册服务接口日志路由
	sysRouter.InitReqApiLogRouter(v1Group, authMiddleware)  // 注册请求接口日志路由
	sysRouter.InitCronjobLogRouter(v1Group, authMiddleware) // 注册任务日志路由
	global.Log.Debug("初始化基础路由完成")
	return r
}
