package websocket

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID     string
	Socket *websocket.Conn
	Send   chan []byte
}

type Manager struct {
	Clients    map[string]*Client
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan []byte
	mutex      sync.RWMutex
}

var WSManager = NewManager()

func NewManager() *Manager {
	return &Manager{
		Clients:    make(map[string]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan []byte),
	}
}

func (m *Manager) Start() {
	for {
		select {
		case client := <-m.Register:
			m.mutex.Lock()
			m.Clients[client.ID] = client
			m.mutex.Unlock()
			log.Printf("新客户端连接: %s", client.ID)

		case client := <-m.Unregister:
			m.mutex.Lock()
			if _, ok := m.Clients[client.ID]; ok {
				close(client.Send)
				delete(m.Clients, client.ID)
			}
			m.mutex.Unlock()
			log.Printf("客户端断开连接: %s", client.ID)

		case message := <-m.Broadcast:
			m.mutex.RLock()
			for _, client := range m.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(m.Clients, client.ID)
				}
			}
			m.mutex.RUnlock()
		}
	}
}

func (m *Manager) SendToAll(message []byte) {
	m.Broadcast <- message
}
