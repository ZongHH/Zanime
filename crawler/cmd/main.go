package main

import (
	"crawler/internal/bootstrap"
	"crawler/pkg/monitor"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	bootstrap := bootstrap.NewBootstrap("../configs/config.yaml")

	bootstrap.Start()

	// 创建一个信号通道，用于接收系统信号
	sigChan := make(chan os.Signal, 1)
	// 注册要接收的信号
	signal.Notify(sigChan,
		syscall.SIGTERM, // 终止信号
		syscall.SIGINT,  // 中断信号 (Ctrl+C)
		syscall.SIGQUIT, // 退出信号
		syscall.SIGHUP,  // 终端断开信号
	)
	// 等待信号到来
	sig := <-sigChan
	monitor.Info("Crawler服务收到信号 %v, 开始优雅关闭", sig)

	bootstrap.Stop()
	time.Sleep(2 * time.Second)
}
