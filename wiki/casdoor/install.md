## 安装casdoor
- 1.Docker安装
``` shell
docker run --name casdoor -d  -e driverName=mysql -e dataSourceName='root:mysql@tcp(172.31.116.212:3306)/' -v /etc/localtime:/etc/localtime -p 8000:8000 registry.cn-shenzhen.aliyuncs.com/dev-ops/casdoor:v1.696.0
```