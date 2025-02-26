package websocket

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	// 允许所有CORS跨域请求
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WsHandler(c *gin.Context) {
	// 升级 HTTP 连接为 WebSocket 连接
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("升级为WebSocket连接失败: %v", err)
		return
	}

	client := &Client{
		ID:     uuid.New().String(),
		Socket: ws,
		Send:   make(chan []byte, 256),
	}

	WSManager.Register <- client

	// 启动读写协程
	go client.writePump()
	go client.readPump()
}

func (c *Client) writePump() {
	defer func() {
		c.Socket.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			err := c.Socket.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Printf("写入消息失败: %v", err)
				return
			}
		}
	}
}

func (c *Client) readPump() {
	defer func() {
		WSManager.Unregister <- c
		c.Socket.Close()
	}()

	for {
		_, message, err := c.Socket.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("读取消息错误: %v", err)
			}
			break
		}

		// 处理接收到的消息
		log.Printf("收到客户端 %s 的消息: %s", c.ID, string(message))
	}
}
