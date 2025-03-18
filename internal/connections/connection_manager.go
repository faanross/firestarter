package connections

import (
	"firestarter/internal/interfaces"
	"firestarter/internal/websocket"
	"fmt"
	"sync"
	"time"
)

// ConnectionManager implements interfaces.ConnectionManager
type ConnectionManager struct {
	connections       map[string]interfaces.Connection
	connectionHistory map[string][]string  // Maps agent UUID to a list of connection IDs
	connectionTimes   map[string]time.Time // Maps connection ID to creation time
	mu                sync.RWMutex
	wsServer          *websocket.WebSocketServer // Allows us to broadcast connections to UI
}

// NewConnectionManager creates a new ConnectionManager with an initialized connections map
func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		connections:       make(map[string]interfaces.Connection),
		connectionHistory: make(map[string][]string),
		connectionTimes:   make(map[string]time.Time),
	}
}

func (cm *ConnectionManager) AddConnection(conn interfaces.Connection) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	id := conn.GetID()
	cm.connections[id] = conn

	// Record connection creation time
	cm.connectionTimes[id] = conn.GetCreatedAt()

	// Track connection history by UUID if available
	agentUUID := conn.GetAgentUUID()
	if agentUUID != "" {
		// Add this connection to the agent's history
		cm.connectionHistory[agentUUID] = append(cm.connectionHistory[agentUUID], id)

		// Check if this is a reconnection
		if len(cm.connectionHistory[agentUUID]) > 1 {
			fmt.Printf("Agent %s reconnected with connection %s\n", agentUUID, id)
		}
	}

	// Log the addition
	fmt.Printf("Connection added: %s (Protocol: %v, UUID: %s, Total active: %d)\n",
		id, conn.GetProtocol(), agentUUID, len(cm.connections))

	if agentUUID != "" {
		fmt.Printf("[UUID-TRACK-DEBUG] Connection manager: Connection %s has UUID %s on addition\n",
			id, agentUUID)
	}

	// Broadcast connection created event to WebSocket clients
	if cm.wsServer != nil {
		// Convert the connection to UI-friendly format
		connInfo := websocket.ConvertConnection(conn)

		// Broadcast the event
		cm.wsServer.Broadcast(websocket.Message{
			Type:    websocket.ConnectionCreated,
			Payload: connInfo,
		})

		fmt.Printf("[WS-DEBUG] Broadcasted ConnectionCreated event for %s\n", id)
	}
}

func (cm *ConnectionManager) RemoveConnection(id string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if conn, exists := cm.connections[id]; exists {
		// The connection still exists in memory, so we can get its UUID
		agentUUID := conn.GetAgentUUID()

		// Keep a copy of the connection info before removing it
		connInfo := websocket.ConvertConnection(conn)

		// Remove from active connections
		delete(cm.connections, id)

		// Note: We intentionally keep the connection in history
		// This preserves the connection history for future reference

		fmt.Printf("Connection removed: %s (UUID: %s, Total remaining: %d)\n",
			id, agentUUID, len(cm.connections))

		// Broadcast connection stopped event to WebSocket clients
		if cm.wsServer != nil {
			cm.wsServer.Broadcast(websocket.Message{
				Type:    websocket.ConnectionStopped,
				Payload: connInfo,
			})

			fmt.Printf("[WS-DEBUG] Broadcasted ConnectionStopped event for %s\n", id)
		}
	}
}

// GetAllConnections returns a slice of all active connections
func (cm *ConnectionManager) GetAllConnections() []interfaces.Connection {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	connections := make([]interfaces.Connection, 0, len(cm.connections))
	for _, conn := range cm.connections {
		connections = append(connections, conn)
	}

	return connections
}

// Count returns the total number of active connections
func (cm *ConnectionManager) Count() int {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	return len(cm.connections)
}

// GetConnection retrieves a specific connection by ID
func (cm *ConnectionManager) GetConnection(id string) (interfaces.Connection, bool) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	conn, exists := cm.connections[id]

	if exists && conn.GetAgentUUID() != "" {
		fmt.Printf("[UUID-de🪲] -> Connection manager: Retrieved connection %s with UUID %s\n",
			id, conn.GetAgentUUID())
	}

	return conn, exists
}

// SetWebSocketServer sets the WebSocket server reference
func (cm *ConnectionManager) SetWebSocketServer(server *websocket.WebSocketServer) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.wsServer = server
	fmt.Println("[🔗] -> Connection manager linked to WebSocket server")
}
