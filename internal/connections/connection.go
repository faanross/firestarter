package connections

import (
	"firestarter/internal/interfaces"
	"fmt"
	"math/rand"
	"time"
)

// BaseConnection provides common fields and functionality for all connection types
type BaseConnection struct {
	ID        string
	Protocol  interfaces.ProtocolType
	Port      string
	CreatedAt time.Time
	AgentUUID string
}

func GenerateUniqueID() string {
	return fmt.Sprintf("connection_%06d", rand.Intn(1000000))
}

// GetAgentUUID returns the agent UUID for this connection
func (bc *BaseConnection) GetAgentUUID() string {
	return bc.AgentUUID
}

// SetAgentUUID updates the agent UUID for this connection
func (bc *BaseConnection) SetAgentUUID(uuid string) {
	if uuid != "" && bc.AgentUUID != uuid {
		fmt.Printf("[UUID-Track-DEBUG] BaseConnection: Updating connection %s with agent UUID: %s\n", bc.ID, uuid)
		bc.AgentUUID = uuid
	}
}
