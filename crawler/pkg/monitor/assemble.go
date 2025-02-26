package monitor

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
	"time"
)

var lp *log_producer

func checkConfig(config *logConfig) {
	if config.ServiceName == "" {
		config.ServiceName = "Crawler"
	}
	if config.TopicName == "" {
		config.TopicName = "ServiceLog"
	}
	if config.PoolSize == 0 {
		config.PoolSize = 2
	}
}

func parseMsg(level string, format string, v ...any) []byte {
	// 获取调用者的信息
	_, file, line, ok := runtime.Caller(3)
	msg := fmt.Sprintf(format, v...)

	// 获取当前时间
	currentTime := time.Now().Format("2006-01-02 15:04:05") // 格式化时间

	LogInfo := log_info{
		Level:     level,
		LevelType: levelType(level),
		Service:   lp.ServiceName,
		TimeStamp: currentTime,
	}

	if ok {
		// 提取当前文件名和上一级目录的文件名
		baseFile := path.Base(file)          // 当前文件名
		dirFile := path.Base(path.Dir(file)) // 上一级目录的文件名
		LogInfo.Message = fmt.Sprintf("%s/%s:%d: %s", dirFile, baseFile, line, msg)
	} else {
		LogInfo.Message = msg
	}

	LogInfoJson, _ := json.Marshal(&LogInfo)

	log.Printf("%s", LogInfo.Message)
	return LogInfoJson
}

func levelType(level string) string {
	if level == "Info" {
		return "info"
	} else if level == "Warning" {
		return "warning"
	} else if level == "Error" {
		return "danger"
	} else if level == "Fatal" {
		return "danger"
	}
	return ""
}

// info_ 记录信息级别日志
func info_(format string, v ...any) {
	_ = parseMsg("Info", format, v...)
	// lp.publish(msg)
}

// warning_ 记录警告级别日志
func warning_(format string, v ...any) {
	_ = parseMsg("Warning", format, v...)
	// lp.publish(msg)
}

// error_ 记录错误级别日志
func error_(format string, v ...any) {
	_ = parseMsg("Error", format, v...)
	// lp.publish(msg)
}

// fatal_ 记录致命错误级别日志
func fatal_(format string, v ...any) {
	_ = parseMsg("Fatal", format, v...)
	// lp.publish(msg)
	lp.close()
	os.Exit(1)
}
