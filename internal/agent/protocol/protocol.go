package protocol

import (
	"context"
	"time"
)

// Protocol defines the interface that all agent-based communication protocols must implement
type Protocol interface {
	// Initialize sets up the protocol with configuration parameters
	Initialize(config Config) error

	// Connect establishes a connection to the server
	Connect(ctx context.Context) error

	// Disconnect terminates the connection to the server
	Disconnect() error

	// IsConnected returns whether the connection is currently active
	IsConnected() bool

	// SendRequest sends a request to the server and returns the response
	SendRequest(ctx context.Context, endpoint string, payload []byte) ([]byte, error)

	// PerformHealthCheck conducts a health check against the server
	PerformHealthCheck(ctx context.Context) error

	// GetLastActivity returns the time of the last successful communication
	GetLastActivity() time.Time

	// Name returns the name of the protocol (e.g., "H1C", "H2C", etc.)
	Name() string
}

// Config holds the common configuration needed by all protocols
type Config struct {
	// TargetHost is the hostname or IP address of the server
	TargetHost string

	// TargetPort is the port the server is listening on
	TargetPort string

	// AgentUUID is the unique identifier for this agent
	AgentUUID string

	// ConnectionTimeout specifies how long to wait for connections
	ConnectionTimeout time.Duration

	// RequestTimeout specifies how long to wait for request responses
	RequestTimeout time.Duration
}
