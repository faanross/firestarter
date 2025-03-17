package context

import (
	"net"
	"net/http"
	"sync"
)

// ConnectionRegistry maps between HTTP requests and their underlying TCP connections
type ConnectionRegistry struct {
	// Map from TCP connection remote address to connection ID
	connMap map[string]string
	// Map from connection ID to agent UUID
	uuidMap map[string]string
	mutex   sync.RWMutex
}

// NewConnectionRegistry creates a new connection registry
func NewConnectionRegistry() *ConnectionRegistry {
	return &ConnectionRegistry{
		connMap: make(map[string]string),
		uuidMap: make(map[string]string),
	}
}

// RegisterConnection associates a TCP connection with a connection ID
func (cr *ConnectionRegistry) RegisterConnection(conn net.Conn, connID string) {
	cr.mutex.Lock()
	defer cr.mutex.Unlock()

	remoteAddr := conn.RemoteAddr().String()
	cr.connMap[remoteAddr] = connID
}

// RegisterUUID associates a connection ID with an agent UUID from an HTTP request
func (cr *ConnectionRegistry) RegisterUUID(req *http.Request, agentUUID string) {
	cr.mutex.Lock()
	defer cr.mutex.Unlock()

	remoteAddr := req.RemoteAddr
	connID, exists := cr.connMap[remoteAddr]
	if exists {
		cr.uuidMap[connID] = agentUUID
	}
}

// GetUUID retrieves the UUID associated with a connection ID
func (cr *ConnectionRegistry) GetUUID(connID string) string {
	cr.mutex.RLock()
	defer cr.mutex.RUnlock()

	return cr.uuidMap[connID]
}
