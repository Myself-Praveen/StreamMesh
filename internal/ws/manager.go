package ws

import (
	"sync"
)

// Manager keeps track of all active WebSocket connections
type Manager struct {
	connections map[string]*Connection
	mu          sync.RWMutex
}

// NewManager creates a new connection manager
func NewManager() *Manager {
	return &Manager{
		connections: make(map[string]*Connection),
	}
}

// Add adds a new connection to the manager
func (m *Manager) Add(conn *Connection) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.connections[conn.ID] = conn
}

// Remove removes a connection by its ID
func (m *Manager) Remove(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.connections, id)
}

// Get retrieves a connection by its ID
func (m *Manager) Get(id string) (*Connection, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	conn, ok := m.connections[id]
	return conn, ok
}

// Count returns the total number of active connections
func (m *Manager) Count() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.connections)
}

// GetAll returns a list of all active connections (snapshot)
func (m *Manager) GetAll() []*Connection {
	m.mu.RLock()
	defer m.mu.RUnlock()
	conns := make([]*Connection, 0, len(m.connections))
	for _, conn := range m.connections {
		conns = append(conns, conn)
	}
	return conns
}
