<h1 align="center">go-gin-rest-api</h1>

<div align="center">
gin + viper + gorm + jwt + casbin实现的golang后台API开发脚手架
<p align="center">
<img src="https://img.shields.io/badge/Golang-1.23-brightgreen" alt="Go version"/>
<img src="https://img.shields.io/badge/Gin-1.10-brightgreen" alt="Gin version"/>
<img src="https://img.shields.io/badge/Gorm-1.25-brightgreen" alt="Gorm version"/> 
</p>
</div>

## 配套前端[ice-antd](https://github.com/linclin/ice-antd)
- [ice.js](https://github.com/alibaba/ice) 阿里开源飞冰React框架 [文档](https://v3.ice.work/docs/guide/start)
- [Ant Design Pro](https://github.com/ant-design/pro-components) 蚂蚁开源Ant Design Pro前端UI组件 [文档](https://procomponents.ant.design/components)

## 主要功能
- 基于`Gin`框架开发的REST API
- API使用`JWT`的Token认证和`Casbin`的接口ACL、RBAC权限控制
- 使用`rql`资源查询语言,支持功能更丰富的数据查询接口
- MySQL数据库使用`Gorm`支持，使用`gorm2-loggable`插件实现数据变更记录表
- 配置文件管理使用`Viper`进行，配置文件变更热加载，无需重启应用(框架基础http服务、日志、数据库不变，仅针对业务配置生效) 
- 定时任务使用`robfig/cron`运行,使用基于MySQL唯一索引的分布式锁避免多副本重复执行，并记录任务执行状态到数据表
- 日志输出使用Golang 1.21版本日志库`slog`，搭配`lumberjack`日志文件切割
- 全局panic日志输出设置使用Golang 1.23版本`debug.SetCrashOutput`，最低运行版本为1.23

## 使用golang大众开源类库
- [gin](https://github.com/gin-gonic/gin) 一款高效的golang web框架 [教程](https://gin-gonic.com/zh-cn/docs/)
- [gorm](https://gorm.io/gorm) 数据库ORM管理框架, 可自行扩展多种数据库类型 [教程](https://gorm.io/zh_CN/)
- [rql](https://github.com/a8m/rql) 用于REST的资源查询语言，作为http请求和orm之间的连接器，gin和gorm的连接器
- [viper](https://github.com/spf13/viper)  配置管理工具, 支持多种配置文件类型.[教程](https://darjun.github.io/2020/01/18/godailylib/viper/)
- [jwt](https://github.com/appleboy/gin-jwt) JWT token中间件
- [casbin](https://github.com/casbin/casbin) 基于角色的访问控制模型(RBAC) [教程](https://darjun.github.io/2020/06/12/godailylib/casbin/)
- [validator](https://github.com/go-playground/validator) 请求参数校验, 版本V10.  [教程](https://darjun.github.io/2020/04/04/godailylib/validator/) 
- [lumberjack](https://github.com/natefinch/lumberjack) 日志切割工具, 高效分离大日志文件, 按日期保存文件
- [cast](https://github.com/spf13/cast) 一个小巧、实用的类型转换库，用于将一个类型转为另一个类型 [教程](https://darjun.github.io/2020/01/20/godailylib/cast/)
- [go-cache](https://github.com/patrickmn/go-cache)  缓存库 [教程](https://cloud.tencent.com/developer/article/2176204)
- [resty](https://github.com/go-resty/resty) Go的简单HTTP和REST请求客户端  [教程](https://darjun.github.io/2021/06/26/godailylib/resty/) 
- [cron](https://github.com/robfig/cron) 实现了 cron 规范解析器和任务运行器，简单来讲就是包含了定时任务所需的功能  [教程](https://darjun.github.io/2020/06/25/godailylib/cron) 
- [lo](https://github.com/samber/lo) 基于泛型实现的Golang工具库

感谢[Go 每日一库](https://github.com/darjun/go-daily-lib)提供的详细教程 


## gin中间件
- [appleboy/gin-jwt](https://github.com/appleboy/gin-jwt)  权限认证中间件 -- 处理登录、登出、无状态token校验
- [casbin](https://github.com/casbin/casbin/) 和[casbin/gorm-adapter](https://github.com/casbin/gorm-adapter)权限访问中间件 -- 基于Cabin RBAC, 对不同角色访问不同API进行校验
- [sentinel](https://github.com/sentinel-group/sentinel-go-adapters/gin) sentinel流量控制中间件 -- api流量控制和熔断降级
- [gin-contrib/requestid](https://github.com/gin-contrib/requestid)  requestid中间件 -- 接口trace id，每次api请求均会生成唯一id返回客户端，并保存数据库接口日志表
- [gin-contrib/cors](https://github.com/gin-contrib/cors)  cors跨域中间件-- 所有请求均可跨域访问  
- [gin-contrib/gzip](https://github.com/gin-contrib/gzip)  gzip中间件-- 所有API返回均进行压缩
- [gin-contrib/pprof](https://github.com/gin-contrib/pprof)  pprof中间件 
- [zsais/go-gin-prometheus](https://github.com/zsais/go-gin-prometheus)  prometheus中间件  
- [swaggo/gin-swagger](https://github.com/swaggo/gin-swagger)  swagger中间件 
- `Exception` 全局异常处理中间件 -- 使用golang recover特性, 捕获所有异常, 保存到日志, 方便追溯  
- `AccessLog` 请求日志中间件 -- 每次请求的路由、IP自动写入日志

## 接口使用指南  
- golang客户端示例代码参见client目录，使用`resty`和`go-cache`编写

- 1.获取JWT生成的token(2小时超时)，AppId和AppSecret存放SysSystem模型
``` shell
curl -X 'POST' \
  'http://127.0.0.1:8080/api/v1/base/auth' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/x-json-stream' \
  -d '{
  "AppId": "2023012801",
  "AppSecret": "fa2e25cb060c8d748fd16ac5210581f41"
}'
```
接口返回
``` shell
{
  "token": "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJBcHBJZCI6IjIwMjMwMTI4MDEiLCJleHAiOjE2NzQ5MjQ3MTcsImlkZW50aXR5IjoiMjAyMzAxMjgwMSIsIm9yaWdfaWF0IjoxNjc0OTE3NTE3fQ.tL0dbrBMOwHJLAvpdHc88OGGxqUiSQHDR_fAupZ4PKsIrDV3CmYgoM4PPrE-3UW2qjV0cE_f_zQ25aRW8-4SNATAoUU1FGd500Ts-SP5WqK30SqEzih0nhhfGS8nRvDoJVzyD7xKip08EwfrmCyVgajYEi1DNC0y02g7jZ47qlManu53xArIFMAzu7bxQxYZPq5DcF0QF8LQipfKNx_LiFdv5ddv_3qf2J3o9uWWGgR_VjZ5p5u4qFFxGla4mSEUX-t9ZtD335D0YJPUayilGCOw7HyyxbdxVRIq1V6R-S17rJFaB48n8pCvpDY_nfs3tbggAMcoJpdBPwvRZlYwKg",
  "expires": "2023-01-29T00:51:57.1720961+08:00"
}
```
- 2.使用token请求API，如下也是rql的使用示例
``` shell
curl -X 'POST' \
  'http://127.0.0.1:8080/api/v1/apilog/list' \
  -H 'accept: application/json' \
  -H 'Authorization: Bearer eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJBcHBJZCI6IjIwMjMwMTI4MDEiLCJleHAiOjE2NzQ5MjQ3MTcsImlkZW50aXR5IjoiMjAyMzAxMjgwMSIsIm9yaWdfaWF0IjoxNjc0OTE3NTE3fQ.tL0dbrBMOwHJLAvpdHc88OGGxqUiSQHDR_fAupZ4PKsIrDV3CmYgoM4PPrE-3UW2qjV0cE_f_zQ25aRW8-4SNATAoUU1FGd500Ts-SP5WqK30SqEzih0nhhfGS8nRvDoJVzyD7xKip08EwfrmCyVgajYEi1DNC0y02g7jZ47qlManu53xArIFMAzu7bxQxYZPq5DcF0QF8LQipfKNx_LiFdv5ddv_3qf2J3o9uWWGgR_VjZ5p5u4qFFxGla4mSEUX-t9ZtD335D0YJPUayilGCOw7HyyxbdxVRIq1V6R-S17rJFaB48n8pCvpDY_nfs3tbggAMcoJpdBPwvRZlYwKg' \
  -H 'Content-Type: application/x-json-stream' \
  -d '{
  "filter": {
    "RequestId": "84692443-987c-4df1-b91c-606fcff6b556"
  },
  "limit": 10,
  "offset": 0,
  "sort": [
    "-StartTime"
  ]
}'
```
接口返回
``` shell
{
  "request_id": "cd7548de-9622-46be-8543-478e70ceb793",
  "success": true,
  "data": {
      {
        "ID": 33,
        "CreatedAt": "2023-01-28T22:51:57+08:00",
        "UpdatedAt": "2023-01-28T22:51:57+08:00",
        "DeletedAt": null,
        "RequestId": "84692443-987c-4df1-b91c-606fcff6b556",
        "RequestMethod": "POST",
        "RequestURI": "/api/v1/base/auth",
        "RequestBody": "{\n  \"AppId\": \"2023012801\",\n  \"AppSecret\": \"fa2e25cb060c8d748fd16ac5210581f41\"\n}",
        "StatusCode": 200,
        "RespBody": "{\"token\":\"eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJBcHBJZCI6IjIwMjMwMTI4MDEiLCJleHAiOjE2NzQ5MjQ3MTcsImlkZW50aXR5IjoiMjAyMzAxMjgwMSIsIm9yaWdfaWF0IjoxNjc0OTE3NTE3fQ.tL0dbrBMOwHJLAvpdHc88OGGxqUiSQHDR_fAupZ4PKsIrDV3CmYgoM4PPrE-3UW2qjV0cE_f_zQ25aRW8-4SNATAoUU1FGd500Ts-SP5WqK30SqEzih0nhhfGS8nRvDoJVzyD7xKip08EwfrmCyVgajYEi1DNC0y02g7jZ47qlManu53xArIFMAzu7bxQxYZPq5DcF0QF8LQipfKNx_LiFdv5ddv_3qf2J3o9uWWGgR_VjZ5p5u4qFFxGla4mSEUX-t9ZtD335D0YJPUayilGCOw7HyyxbdxVRIq1V6R-S17rJFaB48n8pCvpDY_nfs3tbggAMcoJpdBPwvRZlYwKg\",\"expires\":\"2023-01-29T00:51:57.1720961+08:00\"}",
        "ClientIP": "127.0.0.1",
        "StartTime": "2023-01-28T22:51:57+08:00",
        "ExecTime": "4.6134ms"
      }
  },
  "total": 1, 
  "msg": "操作成功"
}
```

## 项目结构概览

```
├── api
│   └── v1 # v1版本接口目录, 如果有新版本可以继续添加v2/v3
├── cronjob # cron定时任务
├── conf # 配置文件目录(包含测试/预发布/生产环境配置参数及casbin模型配置)
├── initialize # 数据初始化目录
├── internal  # 内部方法
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
docker run -d --name mysql -h mysql   --network=host --restart=always  -v /data/mysql:/var/lib/mysql -v /etc/localtime:/etc/localtime -v /etc/resolv.conf:/etc/resolv.conf -e MYSQL_ROOT_PASSWORD=mysql   --restart always   registry.cn-shenzhen.aliyuncs.com/dev-ops/mysql:8.4.2 --character-set-server=utf8mb4 --collation-server=utf8mb4_general_ci
#启动mysql windows
docker run -d --name mysql -h mysql -p 3306:3306 --restart=always -v D:\MySQL:/var/lib/mysql   -e MYSQL_ROOT_PASSWORD=mysql  --restart always   registry.cn-shenzhen.aliyuncs.com/dev-ops/mysql:8.4.2 --character-set-server=utf8mb4 --collation-server=utf8mb4_general_ci 

```

## 快速开始开发
```  
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
## pprof
```  
go tool pprof http://localhost:8080/debug/pprof/heap
go tool pprof http://localhost:8080/debug/pprof/profile
go tool pprof http://localhost:8080/debug/pprof/block
wget http://localhost:8080/debug/pprof/trace?seconds=5
```
## 更新依赖包
```  
go install github.com/oligot/go-mod-upgrade@latest
go-mod-upgrade
或
go get -u && go mod tidy
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
## 镜像构建-Docker构建 
``` shell
# 使用multi-stage(多阶段构建)需要docker 17.05+版本支持
DOCKER_BUILDKIT=1 docker build  --network=host --no-cache --force-rm -t registry.cn-shenzhen.aliyuncs.com/dev-ops/go-gin-rest-api:1.0.0 .
docker push  registry.cn-shenzhen.aliyuncs.com/dev-ops/go-gin-rest-api:1.0.0
docker save -o go-gin-rest-api-1.0.0.tar  registry.cn-shenzhen.aliyuncs.com/dev-ops/go-gin-rest-api:1.0.0
docker load -i go-gin-rest-api-1.0.0.tar
```
## 镜像构建-Buildkit构建
``` shell
# 使用K8S job方式
kubectl apply -f ./wiki/buildkit-build-job.yaml
# 使用Argo Workflow方式
kubectl apply -f ./wiki/buildkit-argo-workflow-template.yaml
```
## 容器运行-Docker运行
``` shell
docker run -d --name go-gin-rest-api -e RunMode=se --network=host -v /etc/localtime:/etc/localtime --restart always registry.cn-shenzhen.aliyuncs.com/dev-ops/go-gin-rest-api:1.0.0
docker run -d --name go-gin-rest-api -e RunMode=prd -e MYSQL_HOST=xx.xx.xx.xx -e MYSQL_PORT=3306 -e MYSQL_USERNAME=root -e MYSQL_PASSWORD=xxxxxxx -e MYSQL_DATABASE=go-gin-rest-api -p 8080:8080 -v /etc/localtime:/etc/localtime --restart always registry.cn-shenzhen.aliyuncs.com/dev-ops/go-gin-rest-api:1.0.0
```
## 容器运行-K8S Deployment运行
``` shell
kubectl apply -f .\k8s-deploy.yaml
```

