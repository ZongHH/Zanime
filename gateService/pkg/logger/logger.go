package logger

import (
	"fmt"
	"os"
	"path"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Log *zap.Logger

// InitLogger 初始化日志
func InitLogger(env string) {
	// 设置日志输出目录
	logPath := "logs"
	if err := os.MkdirAll(logPath, 0777); err != nil {
		panic(fmt.Sprintf("create log directory failed: %v", err))
	}

	// 设置日志轮转配置
	hook := lumberjack.Logger{
		Filename:   path.Join(logPath, "app.log"), // 日志文件路径
		MaxSize:    128,                           // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: 30,                            // 日志文件最多保存多少个备份
		MaxAge:     7,                             // 文件最多保存多少天
		Compress:   true,                          // 是否压缩
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 设置日志级别
	var level zapcore.Level
	if env == "production" {
		level = zap.InfoLevel
	} else {
		level = zap.DebugLevel
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)),
		level,
	)

	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()
	// 开启文件及行号
	development := zap.Development()
	// 设置初始化字段
	filed := zap.Fields(zap.String("service", "gate-service"))
	// 构造日志
	Log = zap.New(core, caller, development, filed)
}
