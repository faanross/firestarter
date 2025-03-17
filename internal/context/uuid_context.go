package context

import "sync"

// UUID mapping
var (
	// Map from connection to UUID
	connectionUUIDs = make(map[string]string)
	uuidMutex       sync.RWMutex
)

// SetConnectionUUID associates a connection ID with an agent UUID
func SetConnectionUUID(connectionID, agentUUID string) {
	uuidMutex.Lock()
	defer uuidMutex.Unlock()

	connectionUUIDs[connectionID] = agentUUID
}

// GetConnectionUUID retrieves the UUID for a connection
func GetConnectionUUID(connectionID string) string {
	uuidMutex.RLock()
	defer uuidMutex.RUnlock()

	return connectionUUIDs[connectionID]
}

// RemoveConnectionUUID removes a connection's UUID mapping
func RemoveConnectionUUID(connectionID string) {
	uuidMutex.Lock()
	defer uuidMutex.Unlock()

	delete(connectionUUIDs, connectionID)
}
