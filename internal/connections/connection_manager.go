package connections

import (
	"firestarter/internal/interfaces"
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
		fmt.Printf("[UUID-Track-DEBUG] Connection manager: Connection %s has UUID %s on addition\n",
			id, agentUUID)
	}
}

func (cm *ConnectionManager) RemoveConnection(id string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if conn, exists := cm.connections[id]; exists {
		// The connection still exists in memory, so we can get its UUID
		agentUUID := conn.GetAgentUUID()

		// Remove from active connections
		delete(cm.connections, id)

		// Note: We intentionally keep the connection in history
		// This preserves the connection history for future reference

		fmt.Printf("Connection removed: %s (UUID: %s, Total remaining: %d)\n",
			id, agentUUID, len(cm.connections))
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
		fmt.Printf("[UUID-Track-DEBUG] Connection manager: Retrieved connection %s with UUID %s\n",
			id, conn.GetAgentUUID())
	}
	
	return conn, exists
}

// GetConnectionsByProtocol returns connections filtered by protocol type
func (cm *ConnectionManager) GetConnectionsByProtocol(protocolType interfaces.ProtocolType) []interfaces.Connection {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	fmt.Printf("[CONN-MGR-DEBUG] GetConnectionsByProtocol called for protocol: %v\n", protocolType)

	filtered := make([]interfaces.Connection, 0)
	for _, conn := range cm.connections {
		if conn.GetProtocol() == protocolType {
			filtered = append(filtered, conn)
		}
	}

	return filtered
}

// LogStatus prints the current connection status to the console
func (cm *ConnectionManager) LogStatus() {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	fmt.Printf("Connection Status: %d active connections\n", len(cm.connections))

	// Group by protocol
	protocolCounts := make(map[interfaces.ProtocolType]int)
	for _, conn := range cm.connections {
		protocolCounts[conn.GetProtocol()]++
	}

	// Print counts by protocol
	for protocol, count := range protocolCounts {
		fmt.Printf("  - Protocol %v: %d connections\n", protocol, count)
	}

	// Print first few connection IDs (limited to avoid overwhelming logs)
	maxToShow := 5
	shown := 0
	fmt.Println("Recent connections:")
	for id := range cm.connections {
		if shown >= maxToShow {
			fmt.Printf("  ... and %d more\n", len(cm.connections)-maxToShow)
			break
		}
		fmt.Printf("  - %s\n", id)
		shown++
	}
}

// GetConnectionsByAgentUUID returns all active connections for a given agent UUID
func (cm *ConnectionManager) GetConnectionsByAgentUUID(agentUUID string) []interfaces.Connection {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	var result []interfaces.Connection

	for _, conn := range cm.connections {
		if conn.GetAgentUUID() == agentUUID {
			result = append(result, conn)
		}
	}

	return result
}

// CountByAgentUUID returns the number of connections for a specific agent UUID
func (cm *ConnectionManager) CountByAgentUUID(agentUUID string) int {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	count := 0
	for _, conn := range cm.connections {
		if conn.GetAgentUUID() == agentUUID {
			count++
		}
	}

	return count
}

// GetUniqueAgentUUIDs returns a list of all unique agent UUIDs
func (cm *ConnectionManager) GetUniqueAgentUUIDs() []string {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	// Use a map to track unique UUIDs
	uuidMap := make(map[string]bool)

	for _, conn := range cm.connections {
		uuid := conn.GetAgentUUID()
		if uuid != "" {
			uuidMap[uuid] = true
		}
	}

	// Convert map keys to slice
	uuids := make([]string, 0, len(uuidMap))
	for uuid := range uuidMap {
		uuids = append(uuids, uuid)
	}

	return uuids
}

// GetConnectionHistoryByUUID returns the history of connections for an agent
func (cm *ConnectionManager) GetConnectionHistoryByUUID(agentUUID string) []string {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	// Return copy of the history slice to prevent mutation
	history := cm.connectionHistory[agentUUID]
	result := make([]string, len(history))
	copy(result, history)

	return result
}

// IsReconnection checks if this is a reconnection from a known agent
func (cm *ConnectionManager) IsReconnection(agentUUID string) bool {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	return len(cm.connectionHistory[agentUUID]) > 1
}
