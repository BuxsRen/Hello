package database

import (
	"Hello/bootstrap/config"
	"Hello/bootstrap/helper"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gohouse/gorose/v2" // https://github.com/gohouse/gorose & https://www.kancloud.cn/fizz/gorose/769179
	"log"
	"os"
)

var MySQL *gorose.Engin
var LOG *log.Logger // 日志

func init() {
	MySQL = connect()
}

type con struct {
	user     string
	pass     string
	host     string
	port     string
	database string
	prefix   string
}

// 连接Mysql，返回实例
func connect() *gorose.Engin {
	var h = helper.Helper{}
	var con = getConfig()

	db, e := gorose.Open(&gorose.Config{
		Driver: "mysql",
		Dsn:    fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=true", con.user, con.pass, con.host, con.port, con.database),
		Prefix: con.prefix,
	})
	if e != nil {
		fmt.Println("➦ " + e.Error())
		h.Exit("✘ Mysql Connection Failed !", 3)
	}
	if config.App.Mysql.Log { // 数据库日志
		logs(h)
		db.Use(func(eg *gorose.Engin) {
			eg.SetLogger(NewLogger(&LogOption{
				EnableSqlLog:   true, // sql日志
				EnableSlowLog:  5,
				EnableErrorLog: true, // 错误日志
			}))
		})
	}

	fmt.Printf("🐬 Mysql -> @tcp(%v:%v)/%v\n", con.host, con.port, con.database)
	return db
}

func getConfig() *con {
	return &con{
		user:     config.App.Mysql.UserName,
		pass:     config.App.Mysql.PassWord,
		host:     config.App.Mysql.Host,
		port:     config.App.Mysql.Port,
		database: config.App.Mysql.Database,
		prefix:   config.App.Mysql.Prefix,
	}
}

// 创建数据库日志保存文件
func logs(h helper.Helper) {
	if config.App.Mysql.SaveLog {
		logFile, e := os.Create(config.App.Mysql.LogPath)
		if e != nil {
			fmt.Println("➦ " + e.Error())
			h.Exit("✘ Mysql Log File Create Failed !", 3)
		}
		LOG = log.New(logFile, "", log.Ldate|log.Ltime|log.Lshortfile)
	}
}