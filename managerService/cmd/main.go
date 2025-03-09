package main

import (
	"managerService/internal/bootstrap"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	container := bootstrap.LoadContainer("../configs/config.yaml")
	bootstrap := bootstrap.NewBootstrap(container)
	bootstrap.Run()

	// 创建一个信号通道，用于接收操作系统的终止信号
	quit := make(chan os.Signal, 1)
	// 注册多种系统信号，包括中断信号(SIGINT)、终止信号(SIGTERM)、
	// 挂起信号(SIGHUP)、退出信号(SIGQUIT)和杀死信号(SIGKILL)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGKILL)
	// 阻塞等待信号，当接收到上述任一信号时继续执行
	<-quit
	// 收到信号后，程序将继续执行并调用 bootstrap.Shutdown() 进行优雅关闭

	bootstrap.Shutdown()
}
