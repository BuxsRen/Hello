package database

import (
	"Hello/app/libs/utils"
	"Hello/bootstrap/config"
	"fmt"
	goruntime "runtime"
	"strings"
	"time"
)

// LogLevel 日志级别
type LogLevel uint

const (
	// LOG_SQL ...
	LOG_SQL LogLevel = iota
	// LOG_SLOW ...
	LOG_SLOW
	// LOG_ERROR ...
	LOG_ERROR
)

// String ...
func (l LogLevel) String() string {
	switch l {
	case LOG_SQL:
		return "SQL"
	case LOG_SLOW:
		return "SLOW"
	case LOG_ERROR:
		return "ERROR"
	}
	return ""
}

// LogOption ...
type LogOption struct {
	FilePath       string  // 日志保存路径
	EnableSqlLog   bool    // sql日志
	EnableSlowLog  float64 // 是否记录慢查询, 默认0s, 不记录, 设置记录的时间阀值, 比如 1, 则表示超过1s的都记录
	EnableErrorLog bool    // 错误日志
}

// Logger ...
type Logger struct {
	filePath string
	sqlLog   bool
	slowLog  float64
	errLog   bool
}

// Engin ...
type Engin struct {
	logger ILogger
}

// ILogger ...
type ILogger interface {
	Sql(sqlStr string, runtime time.Duration)
	Slow(sqlStr string, runtime time.Duration)
	Error(msg string)
	EnableSqlLog() bool
	EnableErrorLog() bool
	EnableSlowLog() float64
}

var _ ILogger = (*Logger)(nil)

//var onceLogger sync.Once
var logger *Logger

// 初始化
func NewLogger(o *LogOption) *Logger {
	//onceLogger.Do(func() {
	logger = &Logger{filePath: "./"}
	if o.FilePath != "" {
		logger.filePath = o.FilePath
	}
	logger.sqlLog = o.EnableSqlLog
	logger.slowLog = o.EnableSlowLog
	logger.errLog = o.EnableErrorLog
	//})
	return logger
}

// 日志是否开启
func (l *Logger) EnableSqlLog() bool {
	return l.sqlLog
}

// 错误日志是否开启
func (l *Logger) EnableErrorLog() bool {
	return l.errLog
}

// 慢日志是否开启
func (l *Logger) EnableSlowLog() float64 {
	return l.slowLog
}

// 记录慢日志
func (l *Logger) Slow(sqlStr string, runtime time.Duration) {
	if l.EnableSlowLog() > 0 && runtime.Seconds() > l.EnableSlowLog() {
		logger.write(LOG_SLOW, "slow", sqlStr, runtime.String())
	}
}

// 日志日志
func (l *Logger) Sql(sqlStr string, runtime time.Duration) {
	if l.EnableSqlLog() {
		logger.write(LOG_SQL, "sql", sqlStr, runtime.String())
	}
}

// 记录错误日志
func (l *Logger) Error(msg string) {
	if l.EnableErrorLog() {
		logger.write(LOG_ERROR, "error", msg, "0")
	}
}

func (l *Logger) write(ll LogLevel, filename string, msg string, runtime string) {
	now := utils.GetNow()
	var content = ""
	var skip int
	if config.App.Server.Debug || config.App.Mysql.SaveLog {
		for {
			skip++
			_, file, line, ok := goruntime.Caller(skip)
			if ok {
				if strings.Index(file, "/app/models/") >= 0 {
					content += fmt.Sprintf("\n[FILE] %s:%d", file, line)
				}
			} else {
				break
			}
		}
	}
	content += fmt.Sprintf("\n[%v] [%v] [%v] (%v)\n\n", ll.String(), now.Format("2006-01-02 15:04:05"), runtime, msg)
	if config.App.Server.Debug {
		fmt.Printf("\x1b[%dm %s \x1b[0m", 35, content)
	}
	if config.App.Mysql.SaveLog {
		_ = LOG.Output(1, content)
	}
}
