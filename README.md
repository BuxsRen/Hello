# Hello App 服务端

> 安卓客户端
>
> https://github.com/BuxsRen/HelloApp

#### 环境与包
| CentOS | Golang | MySQL | Nginx | Redis |  Gin  |  Gorose  |
|:------:|:------:| :---: | :---: | :---: | :---: | :------: |
|   7.9  |  1.18  |  5.7  | 1.18  |  6.2  |  1.7  |  2.1.5   |

#### 安装依赖

```shell
# 启用 go mod
go env -w GO111MODULE=on
#使用七牛云代理
go env -w GOPROXY=https://goproxy.cn,direct

# go mod init Hello
go mod tidy
```

#### 前置

- 复制 app.yaml.example 内容 到 app.yaml

#### 配置 /config/app.yaml 配置文件
- 根据需要进行配置，用不到的配置不用配
```yaml
server: # 服务配置
  url: http://127.0.0.1:9310 # 项目地址/域名
  host: 127.0.0.1 # 服务监听地址，推荐 127.0.0.1
  port: 9310 # 服务监听端口
  udp_port: 3000 # udp服务监听端口
  debug: false # 开启debug模式
  env: local # 运行环境 local(开发) production(线上)
  log_access: ./storage/logs/go_access.log # 访问日志保存路径
  log_error: ./storage/logs/go_error.log # 错误日志保存路径
  template: false # 加载模板 false 的时候 部署可以不需要resources目录

mysql: # mysql 配置
  host: 127.0.0.1  # 数据库地址
  port: 3306       # 数据库端口
  database: xxxxx  # 数据库名
  username: root   # 用户名
  password: 123456 # 密码
  prefix: h_     # 表前缀
  log: false        # 开启sql日志，打印sql执行日志到控制台(server.debug模式打开的时候才会输出到控制台)
  save_log: false   # 保存sql日志到文本，需要先打开 "开启sql日志"，server.debug模式关闭的时候照样可以写入到文件
  log_path: ./storage/logs/sql.log # sql 日志保存地址，需要先打开 "保存sql日志到文本"

redis: # redis 配置
  host: 127.0.0.1  # 地址
  port: 6379 # 端口`
  password: # 密码
  prefix: h_ # redis前缀

other: # 其他配置
  public_dir: /www/wwwroot/public/upload/ # 静态文件保存目录，后面一定要加上 / ,其中 /www/wwwroot 是nginx的静态资源目录
  public_prefix: /public/upload/ # 前端寻址前缀
  token_key: xxxxxxxxxxxx # 接口token签发密钥

push: # 推送配置
  use: false # 开启异常推送(env是production时可用)
  mode: bark # 推送方式：bark、dingTalk、dingTalkMarkDown、wechat
  bark_url: https://api.day.app/xxxxx/ # bark推送地址
  dingTalk_url: https://oapi.dingtalk.com/robot/send?access_token=xxxxxxx # 钉钉推送地址
  wechat_url: https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=xxxxxxxx # 企业微信推送地址

qiniuyun: # 七牛云配置
  bucket: xxx # 七牛云空间名称
  access_key: xxxxxxxxxxx # 七牛云AK
  secret_key: xxxxxxxxxxx # 七牛云SK

email: # 邮箱配置
  name: Break技术团队 # 发件人名称
  user: xxxx@163.com # 发件人邮箱
  pass: xxxxxxxxxxxx # 发件人密码
  host: smtp.163.com # 邮箱服务器
  port: 465 # 邮箱端口

alipay: # 支付宝网页&移动应用支付配置
  appid: 20220122224500 # AppID
  private_key: ./config/alipay.key # 应用私钥路径，后缀也可以是txt，反正是文本就行
  notify_url: https://www.xxx.com/xxxx # 回调地址
```

#### 运行&编译
```shell script
# 运行
go run main.go
# 编译
go build main.go

# Windows 环境下使用build.bat 一键生成打包文件
```

#### 部署清单
```shell script
- Hello         # 主程序
- config/app.yaml         # 配置文件
- resources/views/errors  # 错误页面模板目录
- storage/cache           # 缓存目录
- storage/logs            # 日志目录
```

#### 配置nginx反向代理
```shell script
location /api/go {
    proxy_set_header X-Forward-For $remote_addr;
    proxy_set_header X-real-ip $remote_addr;
    proxy_set_header Host $http_host;
    proxy_set_header SERVER_PORT $server_port;
    proxy_set_header REMOTE_ADDR $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_pass http://127.0.0.1:9310; # 注意修改服务监听端口号
}
```