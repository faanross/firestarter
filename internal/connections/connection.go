package connections

import (
	"encoding/hex"
	"firestarter/internal/types"
	"fmt"
	"time"
)

// Connection interface defines what all protocol-specific connections must implement
type Connection interface {
	GetID() string
	GetProtocol() types.ProtocolType
	GetCreatedAt() time.Time
	GetPort() string
	Close() error
}

// BaseConnection provides common fields and functionality for all connection types
type BaseConnection struct {
	ID        string
	Protocol  types.ProtocolType
	Port      string
	CreatedAt time.Time
}

func GenerateUniqueID() string {
	timestamp := time.Now().UTC().UnixNano()
	randomBytes := make([]byte, 4)
	randomHex := hex.EncodeToString(randomBytes)

	return fmt.Sprintf("conn_%d_%s", timestamp, randomHex)
}
