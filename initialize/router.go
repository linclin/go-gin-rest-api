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
	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
	sentinelPlugin "github.com/sentinel-group/sentinel-go-adapters/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	ginprometheus "github.com/zsais/go-gin-prometheus"
)

// 初始化总路由
func Routers() *gin.Engine {
	// 初始化路由接口到数据库表中
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		//global.Log.Debug(fmt.Sprint("router %v %v %v %v\n", httpMethod, absolutePath, handlerName, nuHandlers))
		go func(httpMethod, absolutePath, handlerName string) {
			sys_router := sys.SysRouter{
				HttpMethod:   httpMethod,
				AbsolutePath: absolutePath,
				HandlerName:  handlerName,
			}
			err := global.Mysql.Where(&sys_router).FirstOrCreate(&sys_router).Error
			if err != nil {
				global.Log.Error(fmt.Sprint("SysRouter 数据初始化失败", err.Error()))
			}
		}(httpMethod, absolutePath, handlerName)
	}
	if global.Conf.System.RunMode == "prd" {
		gin.SetMode(gin.ReleaseMode)
	}
	gin.ForceConsoleColor()
	// 创建带有默认中间件的路由:
	// 日志与恢复中间件
	// r := gin.Default()
	// 创建不带中间件的路由:
	r := gin.New()
	// 初始化Trace中间件
	r.Use(requestid.New())
	// slog日志
	r.Use(sloggin.New(global.Log))
	// 添加访问记录
	r.Use(middleware.AccessLog)
	// 添加全局异常处理中间件
	r.Use(middleware.Exception)
	// GZip压缩插件
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	// 添加跨域中间件, 让请求支持跨域
	// 定义跨域配置
	crosConfig := cors.Config{
		AllowOrigins:     []string{"*", "http://localhost:3000", "http://127.0.0.1:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Access-Control-Allow-Origin", "X-Requested-With", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	// 注册跨域中间件
	r.Use(cors.New(crosConfig))
	// 添加sentinel中间件
	r.Use(
		sentinelPlugin.SentinelMiddleware(
			// customize resource extractor if required
			// method_path by default
			// sentinelPlugin.WithResourceExtractor(func(ctx *gin.Context) string {
			// 	return ctx.GetHeader("X-Real-IP")
			// }),
			// customize block fallback if required
			// abort with status 429 by default
			sentinelPlugin.WithBlockFallback(func(ctx *gin.Context) {
				ctx.AbortWithStatusJSON(429, map[string]interface{}{
					"code": 429,
					"data": "",
					"msg":  "too many request; the quota used up",
				})
			}),
		),
	)
	// 初始化pprof
	pprof.Register(r)
	// 初始化jwt auth中间件
	authMiddleware, err := middleware.InitAuth()
	if err != nil {
		panic(fmt.Sprintf("初始化JWT auth中间件失败: %v", err))
	}
	global.Log.Info("初始化JWT auth中间件完成")
	// 初始化Prometheus中间件
	prome := ginprometheus.NewPrometheus("gin")
	prome.Use(r)
	global.Log.Info("初始化Prometheus中间件完成")
	if global.Conf.System.RunMode != "prd" {
		// 初始化Swagger
		url := ginSwagger.URL(global.Conf.System.BaseApi + "/swagger/doc.json")
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
		global.Log.Info("初始化Swagger完成")
	}
	// 初始化健康检查接口
	r.GET("/heatch_check", api.HeathCheck)
	global.Log.Info("初始化健康检查接口完成")
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
	sysRouter.InitCronjobLogRouter(v1Group, authMiddleware) // 注册定时任务日志路由
	sysRouter.InitChangeLogRouter(v1Group, authMiddleware)  // 注册数据审计日志路由
	sysRouter.InitUserRouter(v1Group, authMiddleware)       // 注册用户权限路由
	sysRouter.InitDataRouter(v1Group, authMiddleware)
	global.Log.Info("初始化基础路由完成")
	return r
}
