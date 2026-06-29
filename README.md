# go-alpha

### 技术栈
* gin框架、Gorm框架
* mysql数据持久化
* redis缓存数据，去重
* CORS跨域：前后端port不同

### pkg
```sh
go get -u "github.com/gin-gonic/gin"
go get -u gorm.io/gorm
go get -u gorm.io/driver/mysql
go get -u "github.com/redis/go-redis/v9"
go ger -u "github.com/gin-contrib/cors"

"github.com/gin-contrib/sessions"
"github.com/gin-contrib/sessions/cookie"
```

### 注意事项

用户信息存放在config.yaml文件中
```
mysql:
  host: localhost
  port: 3306
  user: root
  password: root123456
  dbname: go_alpha

redis:
  host: localhost
  port: 6379
```

docker容器中的mysql，redis改了映射，yaml文件也要改？

### mysql
```shell
mysql -u root -p
```
```sql
create database go_alpha;
```
更改字段要先删除再运行
```sql
DROP TABLE IF EXISTS user; 
```

### CORS
```go
r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
	}))
```

### linux
```shell
kill -9 $(lsof -t -i:8080)
```


### 统计访客人数
方案对比

┌────────────────────────────┬────────────────────────┬──────────────────────────────┬────────────────────┬──────────┐                                                    
│            方案            │          原理          │            准确性            │      实现成本      │ 隐私影响 │                                                    
├────────────────────────────┼────────────────────────┼──────────────────────────────┼────────────────────┼──────────┤                                                    
│ IP                         │ 记录访客 IP            │ 低（NAT/代理导致多人同IP）   │ 低                 │ 低       │                                                    
├────────────────────────────┼────────────────────────┼──────────────────────────────┼────────────────────┼──────────┤                                                    
│ IP + User-Agent            │ IP + UA 组合指纹       │ 中                           │ 低                 │ 低       │                                                    
├────────────────────────────┼────────────────────────┼──────────────────────────────┼────────────────────┼──────────┤                                                    
│ Cookie/ Session            │ 服务端分配唯一标识     │ 高（用户清cookie会重新计数） │ 低                 │ 中       │                                                    
├────────────────────────────┼────────────────────────┼──────────────────────────────┼────────────────────┼──────────┤                                                    
│ 浏览器指纹(Canvas/WebGL等) │ 采集浏览器特征生成哈希 │ 高                           │ 高（需前端JS）     │ 中       │                                                    
├────────────────────────────┼────────────────────────┼──────────────────────────────┼────────────────────┼──────────┤                                                    
│ 用户登录态                 │ 基于 UserID            │ 最高                         │ 低（已有用户体系） │ 低       │                                                    
└────────────────────────────┴────────────────────────┴──────────────────────────────┴────────────────────┴──────────┘

我的推荐：Cookie/Session  

### docker
1. 重新构建镜像并启动
docker compose up -d --build

--build 会重新执行 Dockerfile 里的 go build，把新代码编译进去，然后启动新容器。

2. 如果只想重建后端（不重启 MySQL/Redis）

docker compose up -d --build backend 