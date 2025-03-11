package connections

import "sync"

type ConnectionManager struct {
	connections map[string]Connection
	mu          sync.RWMutex
}

// NewConnectionManager creates a new ConnectionManager with an initialized connections map
func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		connections: make(map[string]Connection),
	}
}

func (cm *ConnectionManager) AddConnection(conn Connection) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.connections[conn.GetID()] = conn
}

func (cm *ConnectionManager) RemoveConnection(id string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	delete(cm.connections, id)
}
