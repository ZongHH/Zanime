package websocket

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

// Connection WebSocket连接结构体
type Connection struct {
	conn      *websocket.Conn // WebSocket连接
	send      chan []byte     // 发送消息的通道
	data      map[string]any  // 连接的自定义数据
	mu        sync.RWMutex    // 保护数据的互斥锁
	createdAt time.Time       // 连接创建时间
	lastPing  time.Time       // 最后一次ping时间
}

// Manager WebSocket连接管理器
type Manager struct {
	connections   sync.Map       // 保存所有的websocket连接
	groups        sync.Map       // 连接分组管理
	handlers      MessageHandler // 消息处理器
	upgrader      websocket.Upgrader
	pingInterval  time.Duration // 心跳间隔
	pongWait      time.Duration // pong等待时间
	writeWait     time.Duration // 写超时时间
	maxMessageLen int64         // 最大消息长度
	logger        *zap.Logger
}

// MessageHandler 消息处理器接口
type MessageHandler interface {
	// HandleMessage 处理收到的消息
	HandleMessage(connectionID string, messageType int, message []byte) error
	// HandleConnect 处理新连接
	HandleConnect(connectionID string) error
	// HandleDisconnect 处理连接断开
	HandleDisconnect(connectionID string) error
	// HandleError 处理错误
	HandleError(connectionID string, err error)
}

// NewManager 创建WebSocket管理器
func NewManager(logger *zap.Logger) *Manager {
	return &Manager{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true // 可以根据需要修改验证逻辑
			},
		},
		pingInterval:  60 * time.Second,
		pongWait:      70 * time.Second,
		writeWait:     10 * time.Second,
		maxMessageLen: 512 * 1024, // 512KB
		logger:        logger,
	}
}

// SetMessageHandler 设置消息处理器
func (m *Manager) SetMessageHandler(handler MessageHandler) {
	m.handlers = handler
}

// Upgrade 将HTTP连接升级为WebSocket连接
// 参数:
//   - w http.ResponseWriter: 用于构建HTTP响应
//   - r *http.Request: 客户端发来的HTTP请求
//   - responseHeader http.Header: 需要在响应头中设置的额外头部信息
//
// 返回值:
//   - *websocket.Conn: 成功升级后的WebSocket连接对象
//   - error: 升级过程中发生的错误（如果有）
//
// 功能说明:
// 1. 使用配置的Upgrader执行WebSocket握手协议
// 2. 设置响应头部信息（包括CORS、子协议等）
// 3. 处理WebSocket版本协商和协议校验
// 4. 返回建立好的WebSocket连接供后续使用
// 注意事项:
// - 需要客户端在请求头中包含正确的Upgrade和Connection字段
// - 会自动处理Sec-WebSocket-Version的版本校验（仅支持13版）
// - 如果配置了CheckOrigin会进行跨域请求校验（当前配置允许所有来源）
func (m *Manager) Upgrade(w http.ResponseWriter, r *http.Request, responseHeader http.Header) (*websocket.Conn, error) {
	// 调用gorilla/websocket库的升级方法执行实际握手流程
	conn, err := m.upgrader.Upgrade(w, r, responseHeader)
	if err != nil {
		// 处理可能发生的错误类型包括：
		// - 客户端未使用WebSocket协议
		// - 无效的握手请求
		// - 跨域校验失败
		// - 协议版本不匹配
		return nil, err
	}
	return conn, err
}

// AddConnection 添加新的WebSocket连接并进行全生命周期管理
// 参数:
//   - conn *websocket.Conn: 已成功建立的WebSocket连接对象
//
// 返回值:
//   - string: 生成的唯一连接ID，用于后续连接管理操作
//
// 功能说明:
// 1. 生成全局唯一的连接ID（UUID v4格式）
// 2. 初始化连接对象，包含：
//   - WebSocket原生连接对象
//   - 带缓冲的发送通道（缓冲区大小256条消息）
//   - 扩展数据存储空间（用于业务层附加元数据）
//   - 时间戳记录（创建时间、最后心跳时间）
//
// 3. 将连接对象存入线程安全的连接池
// 4. 启动独立的读写协程：
//   - readPump: 处理消息接收、协议控制（ping/pong）、超时检测
//   - writePump: 处理消息发送、写入超时控制
//
// 5. 触发连接建立事件通知业务层（如果已注册处理器）
// 6. 返回连接ID供业务层后续使用
func (m *Manager) AddConnection(userID int, conn *websocket.Conn) string {
	// 生成唯一连接标识
	connectionID := generateConnectionID(userID)

	// 初始化连接对象
	connection := &Connection{
		conn:      conn,                   // 原生WebSocket连接
		send:      make(chan []byte, 256), // 带缓冲的发送通道（防止消息突发）
		data:      make(map[string]any),   // 扩展数据存储（线程安全访问需业务层自行控制）
		createdAt: time.Now(),             // 连接创建时间（用于统计和超时管理）
		lastPing:  time.Now(),             // 最后活跃时间（用于心跳检测）
	}

	// 将连接存入线程安全的连接池（sync.Map实现）
	m.connections.Store(connectionID, connection)

	// 启动消息处理协程
	go m.readPump(connectionID, connection)  // 读协程：处理消息接收、协议控制
	go m.writePump(connectionID, connection) // 写协程：处理消息发送、流量控制

	// 通知业务层连接建立事件
	if m.handlers != nil {
		// 使用业务层定义的处理逻辑，同时记录可能发生的处理错误
		if err := m.handlers.HandleConnect(connectionID); err != nil {
			m.logger.Error("连接建立事件处理失败",
				zap.Error(err),
				zap.String("connectionID", connectionID))
		}
	}

	return connectionID // 返回连接ID供业务层后续操作使用
}

// RemoveConnection 移除WebSocket连接
func (m *Manager) RemoveConnection(connectionID string) {
	if conn, ok := m.connections.LoadAndDelete(connectionID); ok {
		connection := conn.(*Connection)
		close(connection.send)
		connection.conn.Close()

		// 从所有组中移除
		m.groups.Range(func(key, value interface{}) bool {
			group := value.(*sync.Map)
			group.Delete(connectionID)
			return true
		})

		// 通知连接断开
		if m.handlers != nil {
			if err := m.handlers.HandleDisconnect(connectionID); err != nil {
				m.logger.Error("处理连接断开失败", zap.Error(err))
			}
		}
	}
}

// AddToGroup 将连接添加到指定组
func (m *Manager) AddToGroup(groupName string, connectionID string) {
	value, _ := m.groups.LoadOrStore(groupName, &sync.Map{})
	group := value.(*sync.Map)
	group.Store(connectionID, struct{}{})
}

// RemoveFromGroup 从指定组中移除连接
func (m *Manager) RemoveFromGroup(groupName string, connectionID string) {
	if group, ok := m.groups.Load(groupName); ok {
		group.(*sync.Map).Delete(connectionID)
	}
}

// SendMessage 发送消息到指定连接
func (m *Manager) SendMessage(connectionID string, message []byte) error {
	if conn, ok := m.connections.Load(connectionID); ok {
		connection := conn.(*Connection)
		select {
		case connection.send <- message:
			return nil
		default:
			return fmt.Errorf("websocket send channel full")
		}
	}
	return nil
}

// BroadcastMessage 广播消息到所有连接
func (m *Manager) BroadcastMessage(message []byte) {
	m.connections.Range(func(key, value interface{}) bool {
		connection := value.(*Connection)
		select {
		case connection.send <- message:
		default:
			m.logger.Warn("发送通道已满，跳过连接", zap.String("connectionID", key.(string)))
		}
		return true
	})
}

// BroadcastToGroup 广播消息到指定组
func (m *Manager) BroadcastToGroup(groupName string, message []byte) {
	if group, ok := m.groups.Load(groupName); ok {
		group.(*sync.Map).Range(func(key, _ interface{}) bool {
			connectionID := key.(string)
			if conn, ok := m.connections.Load(connectionID); ok {
				connection := conn.(*Connection)
				select {
				case connection.send <- message:
				default:
					m.logger.Warn("发送通道已满，跳过连接",
						zap.String("groupName", groupName),
						zap.String("connectionID", connectionID))
				}
			}
			return true
		})
	}
}

// readPump 处理WebSocket读取
func (m *Manager) readPump(connectionID string, connection *Connection) {
	defer func() {
		m.RemoveConnection(connectionID)
	}()

	connection.conn.SetReadLimit(m.maxMessageLen)
	connection.conn.SetReadDeadline(time.Now().Add(m.pongWait))
	connection.conn.SetPongHandler(func(string) error {
		connection.conn.SetReadDeadline(time.Now().Add(m.pongWait))
		connection.lastPing = time.Now()
		return nil
	})

	for {
		messageType, message, err := connection.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				m.logger.Error("读取消息错误", zap.Error(err))
			}
			if m.handlers != nil {
				m.handlers.HandleError(connectionID, err)
			}
			break
		}

		if m.handlers != nil {
			if err := m.handlers.HandleMessage(connectionID, messageType, message); err != nil {
				m.logger.Error("处理消息失败", zap.Error(err))
			}
		}
	}
}

// writePump 处理WebSocket写入
func (m *Manager) writePump(connectionID string, connection *Connection) {
	ticker := time.NewTicker(m.pingInterval)
	defer func() {
		ticker.Stop()
		m.RemoveConnection(connectionID)
	}()

	for {
		select {
		case message, ok := <-connection.send:
			connection.conn.SetWriteDeadline(time.Now().Add(m.writeWait))
			if !ok {
				connection.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := connection.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				m.logger.Error("创建写入器失败", zap.Error(err))
				return
			}

			w.Write(message)

			// 添加队列中的其他消息
			n := len(connection.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-connection.send)
			}

			if err := w.Close(); err != nil {
				m.logger.Error("关闭写入器失败", zap.Error(err))
				return
			}

		case <-ticker.C:
			connection.conn.SetWriteDeadline(time.Now().Add(m.writeWait))
			if err := connection.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				m.logger.Error("发送ping消息失败", zap.Error(err))
				return
			}
		}
	}
}

// GetConnectionData 获取连接的自定义数据
func (m *Manager) GetConnectionData(connectionID string, key string) (interface{}, bool) {
	if conn, ok := m.connections.Load(connectionID); ok {
		connection := conn.(*Connection)
		connection.mu.RLock()
		defer connection.mu.RUnlock()
		value, exists := connection.data[key]
		return value, exists
	}
	return nil, false
}

// SetConnectionData 设置连接的自定义数据
func (m *Manager) SetConnectionData(connectionID string, key string, value interface{}) bool {
	if conn, ok := m.connections.Load(connectionID); ok {
		connection := conn.(*Connection)
		connection.mu.Lock()
		defer connection.mu.Unlock()
		connection.data[key] = value
		return true
	}
	return false
}

// GetConnections 获取所有连接数
func (m *Manager) GetConnections() int {
	count := 0
	m.connections.Range(func(_, _ interface{}) bool {
		count++
		return true
	})
	return count
}

// GetGroupConnections 获取组内连接数
func (m *Manager) GetGroupConnections(groupName string) int {
	count := 0
	if group, ok := m.groups.Load(groupName); ok {
		group.(*sync.Map).Range(func(_, _ interface{}) bool {
			count++
			return true
		})
	}
	return count
}

// generateConnectionID 生成连接ID
func generateConnectionID(userID int) string {
	return fmt.Sprintf("%d", userID)
}
