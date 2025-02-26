package monitor

import (
	"log"
)

func Fatal(format string, v ...any) {
	fatal_(format, v...)
}

func Error(format string, v ...any) {
	error_(format, v...)
}

func Warning(format string, v ...any) {
	warning_(format, v...)
}

func Info(format string, v ...any) {
	info_(format, v...)
}

func Init(config *logConfig) {
	checkConfig(config)

	var err error
	lp, err = newLogProducer(config)
	if err != nil {
		log.Fatal("newLogProducer ", err)
	}
}

func Close() {
	lp.close()
}

func NewLogConfig() *logConfig {
	return &logConfig{}
}
