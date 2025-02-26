package controller

import "monitorService/service/websocket"

func (c *Controller) setAPIRouters() {
	// 添加WebSocket路由
	c.engine.GET("/ws", websocket.WsHandler)
}
