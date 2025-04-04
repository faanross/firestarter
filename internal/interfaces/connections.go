package interfaces

import "time"

// ProtocolType defines the supported protocol types
type ProtocolType int

const (
	H1C ProtocolType = iota + 1
	H1TLS
	H2C
	H2TLS
	H3
)

// Connection defines what all protocol-specific connections must implement
type Connection interface {
	GetID() string
	GetProtocol() ProtocolType
	GetCreatedAt() time.Time
	GetPort() string
	GetAgentUUID() string
	SetAgentUUID(string)
	Close() error
}

// ConnectionManager defines the interface for managing connections
type ConnectionManager interface {
	AddConnection(conn Connection)
	RemoveConnection(id string)
	GetAllConnections() []Connection
	Count() int
	GetConnection(id string) (Connection, bool)
}

// Helper function to get protocol name
func GetProtocolName(protocol ProtocolType) string {
	switch protocol {
	case H1C:
		return "HTTP/1.1"
	case H2C:
		return "HTTP/2"
	case H1TLS:
		return "HTTP/1.1 TLS"
	case H2TLS:
		return "HTTP/2 TLS"
	case H3:
		return "HTTP/3"
	default:
		return "Unknown"
	}
}
