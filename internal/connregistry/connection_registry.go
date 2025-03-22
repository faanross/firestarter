package connregistry

import (
	"firestarter/internal/interfaces"
	"fmt"
	"net"
	"net/http"
	"sync"
)

// Global connection registry instance
var GlobalConnectionRegistry *ConnectionRegistry

// ConnectionRegistry maps between HTTP requests and their underlying TCP connections
type ConnectionRegistry struct {
	connMap        map[string]string            // Map from TCP connection remote address to connection ID
	uuidMap        map[string]string            // Map from connection ID to agent UUID
	connManager    interfaces.ConnectionManager // Connection manager reference
	processedPairs map[string]bool              // Tracks already processed remoteAddr:UUID pairs
	mutex          sync.RWMutex
}

// NewConnectionRegistry creates a new connection registry
func NewConnectionRegistry() *ConnectionRegistry {
	fmt.Println("[🌎REG] -> Global Connection Registry initialized.")
	return &ConnectionRegistry{
		connMap:        make(map[string]string), // remoteAddr : conn ID
		uuidMap:        make(map[string]string), // conn ID : UUID
		processedPairs: make(map[string]bool),   // remoteAddr : UUID
	}
}

// SetConnectionManager sets the connection manager for this registry
func (cr *ConnectionRegistry) SetConnectionManager(cm interfaces.ConnectionManager) {
	cr.mutex.Lock()
	defer cr.mutex.Unlock()
	cr.connManager = cm
	fmt.Println("[🔗LNK] -> Global Registry linked to Connection Manager.")
}

// RegisterConnection associates a TCP connection with a connection ID
func (cr *ConnectionRegistry) RegisterConnection(conn net.Conn, connID string) {
	cr.mutex.Lock()
	defer cr.mutex.Unlock()

	remoteAddr := conn.RemoteAddr().String()
	cr.connMap[remoteAddr] = connID

	fmt.Printf("[🟢NEW] -> New Connection from: %s\n", remoteAddr)
	fmt.Printf("[REG🌎] -> Mapped %s to connection ID %s\n",
		remoteAddr, connID)
}

// RegisterUUID associates a connection ID with an agent UUID from an HTTP request
func (cr *ConnectionRegistry) RegisterUUID(req *http.Request, agentUUID string) {
	// Create a unique key for this request+UUID combination
	pairKey := req.RemoteAddr + ":" + agentUUID

	cr.mutex.RLock()
	alreadyProcessed := cr.processedPairs[pairKey]
	cr.mutex.RUnlock()

	// Skip if we've already processed this exact combination
	if alreadyProcessed {
		return
	}

	cr.mutex.Lock()
	defer cr.mutex.Unlock()

	// Mark as processed
	cr.processedPairs[pairKey] = true

	remoteAddr := req.RemoteAddr
	connID, exists := cr.connMap[remoteAddr]
	if !exists {
		fmt.Printf("Warning: No connection ID found for remote address: %s\n", remoteAddr)
		return
	}

	// Update our UUID map
	cr.uuidMap[connID] = agentUUID
	fmt.Printf("Registry: Associated connection %s with UUID %s\n", connID, agentUUID)

	// Update the actual connection object if we have a connection manager
	if cr.connManager != nil {
		if conn, found := cr.connManager.GetConnection(connID); found {
			// Type assertion to check if we can set the UUID
			if setter, ok := conn.(interface{ SetAgentUUID(string) }); ok {
				setter.SetAgentUUID(agentUUID)
				fmt.Printf("Registry: Updated connection object %s with UUID %s\n", connID, agentUUID)
			} else {
				fmt.Printf("Warning: Connection %s doesn't support SetAgentUUID\n", connID)
			}
		} else {
			fmt.Printf("Warning: Connection %s not found in manager\n", connID)
		}
	}
}

// GetRemoteAddrByConnID retrieves the remote address associated with a connection ID
func (cr *ConnectionRegistry) GetRemoteAddrByConnID(connID string) string {
	cr.mutex.RLock()
	defer cr.mutex.RUnlock()

	// We need to search through the connMap to find the entry where the value is the connID
	for remoteAddr, id := range cr.connMap {
		if id == connID {
			return remoteAddr
		}
	}

	return "" // Return empty string if not found
}

// InitializeConnectionRegistry creates and sets up the global connection registry
func InitializeConnectionRegistry() {
	if GlobalConnectionRegistry == nil {
		GlobalConnectionRegistry = NewConnectionRegistry()
	}
}

// GetConnectionRegistry returns the global connection registry
func GetConnectionRegistry() *ConnectionRegistry {
	return GlobalConnectionRegistry
}

// ConnectRegistryToManager connects the registry to a connection manager
func ConnectRegistryToManager(manager interfaces.ConnectionManager) {
	if GlobalConnectionRegistry != nil {
		GlobalConnectionRegistry.SetConnectionManager(manager)
	} else {
		fmt.Println("[❌ERR] -> Cannot connect registry to manager - registry not initialized")
	}
}
