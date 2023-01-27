<h1 align="center">go-gin-rest-api</h1>

<div align="center">
gin + viper + gorm + jwt + casbin实现的golang后台API开发脚手架
</div>

## 使用开源类库，选择github star最多的库
- [gin](https://github.com/gin-gonic/gin) 一款高效的golang web框架 [教程](https://gin-gonic.com/zh-cn/docs/)
- [gorm](https://gorm.io/gorm) 数据库ORM管理框架, 可自行扩展多种数据库类型 [教程](https://gorm.io/zh_CN/)
- [rql](https://github.com/a8m/rql) 用于REST的资源查询语言，作为http请求和orm之间的连接器，gin和gorm的连接器
- [viper](https://github.com/spf13/viper)  配置管理工具, 支持多种配置文件类型.[教程](https://darjun.github.io/2020/01/18/godailylib/viper/)
- [jwt](https://github.com/appleboy/gin-jwt) JWT token中间件
- [casbin](https://github.com/casbin/casbin) 基于角色的访问控制模型(RBAC) [教程](https://darjun.github.io/2020/06/12/godailylib/casbin/)
- [validator](https://github.com/go-playground/validator) 请求参数校验, 版本V10.  [教程](https://darjun.github.io/2020/04/04/godailylib/validator/)
- [zap](https://github.com/uber-go/zap): uber开源的日志库zap，对性能和内存分配做了极致的优化.  [教程](https://darjun.github.io/2020/04/23/godailylib/zap/)
- [lumberjack](https://github.com/natefinch/lumberjack) 日志切割工具, 高效分离大日志文件, 按日期保存文件
- [cast](https://github.com/spf13/cast) 一个小巧、实用的类型转换库，用于将一个类型转为另一个类型 [教程](https://darjun.github.io/2020/01/20/godailylib/cast/)
- [go-cache](https://github.com/patrickmn/go-cache)  缓存库
- [resty](https://github.com/go-resty/resty) Go的简单HTTP和REST请求客户端   
- [cron](https://github.com/robfig/cron) 实现了 cron 规范解析器和任务运行器，简单来讲就是包含了定时任务所需的功能  [教程](https://darjun.github.io/2020/06/25/godailylib/cron) 
- [tunny](https://github.com/Jeffail/tunny) 协程池，支持同步执行汇总结果,支持超时、取消 [教程](https://https://darjun.github.io/2021/06/10/godailylib/tunny/) 

感谢[Go 每日一库](https://github.com/darjun/go-daily-lib)提供的详细教程 


## gin中间件
- `RateLimiter` 访问速率限制中间件 -- 限制访问流量
- `Exception` 全局异常处理中间件 -- 使用golang recover特性, 捕获所有异常, 保存到日志, 方便追溯
- `Transaction` 全局事务处理中间件 -- 每次请求无异常自动提交, 有异常自动回滚事务, 无需每个service单独调用(GET/OPTIONS跳过)
- `AccessLog` 请求日志中间件 -- 每次请求的路由、IP自动写入日志
- `Cors`  跨域中间件-- 所有请求均可跨域访问
- `JwtAuth` 权限认证中间件 -- 处理登录、登出、无状态token校验
- `CasbinMiddleware` 权限访问中间件 -- 基于Cabin RBAC, 对不同角色访问不同API进行校验

## 项目结构概览

```
├── api
│   └── v1 # v1版本接口目录, 如果有新版本可以继续添加v2/v3
├── cronjob # cron定时任务
├── conf # 配置文件目录(包含测试/预发布/生产环境配置参数及casbin模型配置)
├── initialize # 数据初始化目录
├── logs # 日志文件默认目录(运行代码时生成)
├── middleware # 中间件目录
├── models # 存储层模型定义目录
├── pkg # 公共模块目录
│   ├── global # 全局变量目录 
├── router # 路由目录 
├── vendor # 第三方依赖库 
``` 
## MySQL数据库准备
```  
#启动mysql linux
docker run -d --name mysql -h mysql   --network=host --restart=always  -v /data/mysql:/var/lib/mysql -v /etc/localtime:/etc/localtime -v /etc/resolv.conf:/etc/resolv.conf -e MYSQL_ROOT_PASSWORD=mysql   --restart always   mysql:8.0.32 --character-set-server=utf8mb4 --collation-server=utf8mb4_general_ci --default-authentication-plugin=mysql_native_password
#启动mysql windows
docker run -d --name mysql -h mysql -p 3306:3306 --restart=always -v D:\MySQL:/var/lib/mysql   -e MYSQL_ROOT_PASSWORD=mysql  --restart always   mysql:8.0.32 --character-set-server=utf8mb4 --collation-server=utf8mb4_general_ci --default-authentication-plugin=mysql_native_password

```

## 快速开始开发

``` 
golang版本>1.19
# 设置常用环境变量
go env -w GOROOT=C:\Go 
go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.cn,direct
go env -w GOSUMDB=off 
# 运行
go run main.go

# 热编译运行 gowatch
go install github.com/silenceper/gowatch@latest
gowatch
```

## Prometheus exportor集成

```  
访问地址:http://127.0.0.1:8080/metrics
```

## Swagger文档自动生成

``` 
go install github.com/swaggo/swag/cmd/swag@latest 
swag init --parseDependency --parseVendor --parseInternal --parseDepth 1

访问地址:http://127.0.0.1:8080/swagger/index.html
```

## 更新依赖包
```  
go install github.com/oligot/go-mod-upgrade@latest
go-mod-upgrade
```

## 编译
``` 
##编译时环境变量
#git分支
GitBranch=$(git name-rev --name-only HEAD)
GitRevision=$(git rev-parse --short HEAD) 
#git commit
GitCommitLog=`git log --pretty=oneline -n 1` 
GitCommitLog=${GitCommitLog//\'/\"}
#构建时间
BuildTime=`date +'%Y.%m.%d.%H%M%S'`
#构建go版本
BuildGoVersion=`go version`
LDFlags="-w -X 'main.GitBranch=${GitBranch}' -X 'main.GitRevision=${GitRevision}' -X 'main.GitCommitLog=${GitCommitLog}' -X 'main.BuildTime=${BuildTime}' -X 'main.BuildGoVersion=${BuildGoVersion}'"
go build -ldflags="$LDFlags" -o  ./go-gin-rest-api
```
## Docker 镜像制作
``` shell
# 使用multi-stage(多阶段构建)需要docker 17.05+版本支持
DOCKER_BUILDKIT=1 docker build  --network=host --no-cache --force-rm  -t   go-gin-rest-api:1.0.0 .
docker push  go-gin-rest-api
docker save -o go-gin-rest-api-1.0.0.tar  go-gin-rest-api:1.0.0
docker load -i go-gin-rest-api-1.0.0.tar
```


