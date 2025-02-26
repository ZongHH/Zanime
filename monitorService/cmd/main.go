package main

import (
	"log"
	"monitorService/controller"
	"monitorService/dao/mysql"
	"monitorService/pkg/config"
	"monitorService/service"
	"os"
	"os/signal"
	"syscall"
)

func Init() {
	config.ConfigInit()

	// 初始化数据库连接
	dbConfig := &mysql.Config{
		Username: "root",
		Password: "123456",
		Host:     "localhost",
		Port:     3306,
		Database: "anime",
	}

	if err := mysql.InitDB(dbConfig); err != nil {
		log.Fatal(err)
	}

	service.Init()
}

func Close() {
	service.Close()
	controller.Close()
	mysql.Close()
}

func main() {
	Init()

	// 启动 HTTP 服务
	controller.Init()

	// 创建信号接收通道
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 等待中断信号
	<-quit

	Close()
}
