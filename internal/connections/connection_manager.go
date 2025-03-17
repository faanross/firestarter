package connections

import (
	"firestarter/internal/interfaces"
	"fmt"
	"sync"
)

// ConnectionManager implements interfaces.ConnectionManager
type ConnectionManager struct {
	connections map[string]interfaces.Connection
	mu          sync.RWMutex
}

// NewConnectionManager creates a new ConnectionManager with an initialized connections map
func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		connections: make(map[string]interfaces.Connection),
	}
}

func (cm *ConnectionManager) AddConnection(conn interfaces.Connection) {

	fmt.Printf("[CONN-MGR-DEBUG] AddConnection called for connection ID: %s, Protocol: %v\n", conn.GetID(), conn.GetProtocol())

	cm.mu.Lock()
	defer cm.mu.Unlock()

	id := conn.GetID()
	cm.connections[id] = conn

	// Log the addition
	fmt.Printf("Connection added: %s (Protocol: %v, Total active: %d)\n",
		id, conn.GetProtocol(), len(cm.connections))
}

func (cm *ConnectionManager) RemoveConnection(id string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if _, exists := cm.connections[id]; exists {
		delete(cm.connections, id)
		fmt.Printf("Connection removed: %s (Total remaining: %d)\n",
			id, len(cm.connections))
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

// GetConnectionsByAgentUUID returns connections associated with a specific agent UUID
func (cm *ConnectionManager) GetConnectionsByAgentUUID(agentUUID string) []interfaces.Connection {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	var results []interfaces.Connection
	for _, conn := range cm.connections {
		if conn.GetAgentUUID() == agentUUID {
			results = append(results, conn)
		}
	}

	return results
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
