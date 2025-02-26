package service

import (
	"monitorService/service/collect"
	"monitorService/service/websocket"
)

func Init() {
	collect.Init()
	// 启动WebSocket管理器
	go websocket.WSManager.Start()
}

func Close() {
	collect.Close()
}
