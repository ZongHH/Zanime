package websocket

// type Connection struct {
// 	conn    *websocket.Conn
// 	send    chan []byte
// 	closeCh chan struct{}
// }

// type Manager struct {
// 	connections sync.Map // key: connection ID, value: *Connection
// 	upgrader    websocket.Upgrader
// 	mu          sync.Mutex
// }

// func NewManager() *Manager {
// 	return &Manager{
// 		upgrader: websocket.Upgrader{
// 			ReadBufferSize:  1024,
// 			WriteBufferSize: 1024,
// 		},
// 	}
// }

// func (m *Manager) AddConnection(conn *websocket.Conn) string {
// 	c := &Connection{
// 		conn:    conn,
// 		send:    make(chan []byte, 256),
// 		closeCh: make(chan struct{}),
// 	}

// 	id := generateID()
// 	m.connections.Store(id, c)

// 	go m.readPump(c, id)
// 	go m.writePump(c)

// 	return id
// }

// func (m *Manager) RemoveConnection(id string) {
// 	if conn, ok := m.connections.LoadAndDelete(id); ok {
// 		c := conn.(*Connection)
// 		close(c.send)
// 		c.conn.Close()
// 		close(c.closeCh)
// 	}
// }

// func (m *Manager) GetConnection(id string) (*Connection, bool) {
// 	conn, ok := m.connections.Load(id)
// 	if !ok {
// 		return nil, false
// 	}
// 	return conn.(*Connection), true
// }

// func (m *Manager) readPump(c *Connection, id string) {
// 	defer m.RemoveConnection(id)

// 	c.conn.SetReadLimit(512)
// 	for {
// 		_, message, err := c.conn.ReadMessage()
// 		if err != nil {
// 			if websocket.IsUnexpectedCloseError(err) {
// 				log.Printf("error: %v", err)
// 			}
// 			break
// 		}
// 		log.Printf("Received message: %s", message)
// 	}
// }

// func (m *Manager) writePump(c *Connection) {
// 	defer c.conn.Close()

// 	for {
// 		select {
// 		case message, ok := <-c.send:
// 			if !ok {
// 				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
// 				return
// 			}

// 			w, err := c.conn.NextWriter(websocket.TextMessage)
// 			if err != nil {
// 				return
// 			}
// 			w.Write(message)

// 			if err := w.Close(); err != nil {
// 				return
// 			}
// 		case <-c.closeCh:
// 			return
// 		}
// 	}
// }

// func (c *Connection) SendMessage(message []byte) {
// 	defer func() {
// 		if r := recover(); r != nil {
// 			log.Println("SendMessage panic:", r)
// 		}
// 	}()

// 	select {
// 	case c.send <- message:
// 	default:
// 		close(c.send)
// 	}
// }

// func (m *Manager) Broadcast(message []byte) {
// 	m.connections.Range(func(key, value interface{}) bool {
// 		conn := value.(*Connection)
// 		conn.SendMessage(message)
// 		return true
// 	})
// }

// func (m *Manager) CloseAll() {
// 	m.connections.Range(func(key, value interface{}) bool {
// 		m.RemoveConnection(key.(string))
// 		return true
// 	})
// }

// func generateID() string {
// 	// 实现自己的ID生成逻辑
// 	return "unique_id"
// }
