package service

import (
	"context"
	"gateService/internal/interfaces/dto"
)

// WebSocketService 定义WebSocket服务接口
// 该接口负责处理WebSocket连接的建立和管理，包括：
// - 新连接的认证和初始化
// - 连接状态的维护
// - 消息协议的协商
type WebSocketService interface {
	// EstablishConnection 建立新的WebSocket连接
	// 参数:
	//   ctx context.Context      - 上下文，用于控制请求生命周期和取消操作
	//   request *dto.EstablishWebSocketRequest - 包含连接请求参数，如：
	//     - 用户身份信息
	//     - 客户端设备信息
	//     - 认证令牌
	// 返回值:
	//   *dto.EstablishWebSocketResponse - 包含连接响应信息，如：
	//     - 连接状态码
	//     - 连接唯一标识符
	//     - 心跳超时时间
	//   error - 错误信息，可能包含：
	//     - 参数校验失败错误
	//     - 身份认证失败错误
	//     - 服务器内部错误
	// 注意事项:
	//   - 需要处理并发连接请求
	//   - 需要记录连接日志
	//   - 需要实现连接超时机制
	EstablishConnection(ctx context.Context, request *dto.EstablishWebSocketRequest) (*dto.EstablishWebSocketResponse, error)
}
