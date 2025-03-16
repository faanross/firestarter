package connections

import (
	"firestarter/internal/interfaces"
	"fmt"
	"math/rand"
	"time"
)

// Connection interface defines what all protocol-specific connections must implement
type Connection interface {
	GetID() string
	GetProtocol() interfaces.ProtocolType
	GetCreatedAt() time.Time
	GetPort() string
	Close() error
}

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
