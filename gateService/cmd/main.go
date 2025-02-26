package main

import (
	"gateService/internal/bootstrap"
	"gateService/pkg/errors"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gateService/pkg/logger"
)

func main() {
	// 初始化日志
	logger.InitLogger("development") // 或 "production"

	// 设置错误处理配置
	errors.SetErrorConfig(errors.ErrorConfig{
		Env:            "development",
		AlertThreshold: 10,
		AlertInterval:  time.Minute * 5,
	})

	bootstrap := bootstrap.NewBootstrap("../configs/config.yaml")
	bootstrap.Start()
	defer bootstrap.Stop()

	done := make(chan struct{})

	// 创建一个信号通道，用于接收系统信号
	sigChan := make(chan os.Signal, 1)
	// 注册要处理的信号
	signal.Notify(sigChan,
		syscall.SIGTERM, // 终止信号
		syscall.SIGINT,  // 中断信号 (Ctrl+C)
		syscall.SIGQUIT, // 退出信号
		syscall.SIGHUP,  // 终端断开信号
	)

	go func() {
		sig := <-sigChan
		log.Printf("收到系统信号: %v，开始优雅退出", sig)
		done <- struct{}{}
	}()

	<-done
}
