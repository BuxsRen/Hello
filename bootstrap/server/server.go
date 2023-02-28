package server

import (
	"Hello/app/exceptions"
	"Hello/app/libs/utils"
	"Hello/app/socket/udpsocket"
	"Hello/bootstrap/config"
	"Hello/bootstrap/helper"
	"Hello/bootstrap/routes"
	"fmt"
	"github.com/gin-gonic/gin" // https://gin-gonic.com/zh-cn/docs/ Gin开发文档
	"io"
	"log"
	"math/rand"
	"os"
	"runtime"
	"sync"
)

func init() {
	rand.Seed(utils.GetNow().UnixNano()) // 初始化随机数种子
}

type server struct {
	server    *gin.Engine
	h         *helper.Helper
	host      string
	port      string
	debug     bool
	logAccess string
	LogError  string
	template  bool
}

// 启动应用
func Run() {
	app := &server{
		host:      config.App.Server.Host,
		port:      config.App.Server.Port,
		debug:     config.App.Server.Debug,
		logAccess: config.App.Server.LogAccess,
		LogError:  config.App.Server.LogError,
		template:  config.App.Server.Template,
	}
	app.print()
	app.h = &helper.Helper{}
	app.isDebug()
	app.udp()
	app.server = gin.Default()
	_ = app.server.SetTrustedProxies(nil)
	app.loadTemplate()
	app.server.Use(exceptions.Handle)            // 异常处理
	(&routes.Route{Router: app.server}).Handle() // 加载路由
	app.start()
}

// 输出信息
func (this *server) print() {
	fmt.Printf("⏰ Runtime: %s\n", utils.GetNow().Format("2006-01-02 15:04:05"))
	fmt.Printf("🚀 Server: 2021 - %d By Break\n", utils.GetNow().Year())
	fmt.Printf("💻 System：%s (%s)\n", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("🌏 http://127.0.0.1:%s\n", this.port)
	fmt.Printf("🔗 Listen Web Server -> %s:%s\n", this.host, this.port)
	fmt.Printf("► OK! Start Web Service ...\n")
}

// 启动udp服务
func (this *server) udp() {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go udpsocket.Run(wg)
	wg.Wait()
}

// 是否开启debug
func (this *server) isDebug() {
	if !this.debug { //未开启debug 异常错误日志写入文件中,开启debug 异常打印到终端并返回给浏览器
		gin.SetMode(gin.ReleaseMode)
		access, _ := os.Create(this.logAccess)
		gin.DefaultWriter = io.MultiWriter(access) // 访问日志
		this.logs()
	}
}

// 是否加载模板
func (this *server) loadTemplate() {
	if this.template { // 加载模板
		this.server.LoadHTMLGlob("./resources/views/**/*")
	}
}

// 创建输出日志文件
func (this *server) logs() {
	logFile, e := os.Create(config.App.Server.LogError)
	if e != nil {
		fmt.Println("➦ " + e.Error())
		this.h.Exit("✘ Error Log File Create Failed !", 3)
	}
	helper.LOG = log.New(logFile, "", log.Ldate|log.Ltime|log.Lshortfile)
}

// 开始web服务
func (this *server) start() {
	err := this.server.Run(this.host + ":" + this.port) // 启动服务
	if err != nil {
		fmt.Println("➦ " + err.Error())
		this.h.Exit("x This Port Is Already In Use!", 0)
	}
}
