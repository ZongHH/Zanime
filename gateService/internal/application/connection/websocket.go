package connection

import (
	"context"
	"fmt"
	"gateService/internal/infrastructure/middleware/websocket"
	"gateService/internal/interfaces/dto"

	"github.com/gin-gonic/gin"
)

type WebSocketServiceImpl struct {
	websocketManager *websocket.Manager
}

func NewWebSocketServiceImpl(webSocketManager *websocket.Manager) *WebSocketServiceImpl {
	return &WebSocketServiceImpl{
		websocketManager: webSocketManager,
	}
}

// EstablishConnection 建立WebSocket连接并关联用户信息
// 参数:
//
//	ctx - 上下文对象，需要包含*gin.Context用于HTTP升级
//	request - 包含用户ID等身份信息的请求参数
//
// 返回值:
//
//	*dto.EstablishWebSocketResponse - 包含状态码的响应对象
//	error - 错误信息，包含具体的失败原因
func (w *WebSocketServiceImpl) EstablishConnection(ctx context.Context, request *dto.EstablishWebSocketRequest) (*dto.EstablishWebSocketResponse, error) {
	// 类型断言获取gin.Context，用于WebSocket协议升级
	// 需要确保传入的context是gin框架的Context对象
	c, ok := ctx.(*gin.Context)
	if !ok {
		return nil, fmt.Errorf("*gin.Context转换失败: 上下文类型不匹配，请确保使用gin框架的Context")
	}

	// 升级HTTP连接到WebSocket协议
	// 使用websocket管理器进行协议升级，需要访问ResponseWriter和Request对象
	conn, err := w.websocketManager.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return nil, fmt.Errorf("升级WebSocket连接失败: %v（可能原因：协议不匹配/头信息错误）", err)
	}

	// 生成唯一连接ID并注册到连接管理器
	// 每个WebSocket连接都会分配一个唯一标识符用于后续管理
	connectionID := w.websocketManager.AddConnection(request.UserID, conn)

	// 在连接元数据中存储用户ID
	// 将业务相关的用户标识与连接ID绑定，方便后续消息路由
	w.websocketManager.SetConnectionData(connectionID, "user_id", request.UserID)

	// 返回成功响应（状态码200表示连接已成功建立）
	// 注意：实际WebSocket通信会在连接建立后异步进行
	return &dto.EstablishWebSocketResponse{
		Code: 200, // HTTP状态码，表示协议升级成功
	}, nil
}
