package connections

import (
	"encoding/hex"
	"firestarter/internal/interfaces"
	"fmt"
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
}

func GenerateUniqueID() string {
	timestamp := time.Now().UTC().UnixNano()
	randomBytes := make([]byte, 4)
	randomHex := hex.EncodeToString(randomBytes)

	return fmt.Sprintf("conn_%d_%s", timestamp, randomHex)
}
