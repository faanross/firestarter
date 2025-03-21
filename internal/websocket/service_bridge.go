package websocket

import (
	"firestarter/internal/interfaces"
	"firestarter/internal/types"
)

// ServiceBridge acts as contract between the WebSocket server and the service layer
type ServiceBridge interface {
	GetAllListeners() []types.Listener
	StopListener(id string) error
	GetAllConnections() []interfaces.Connection
	StopConnection(id string) error
	IsPortAvailable(port string) bool
	CreateListener(id string, protocol int, port string) (types.Listener, error)
}

// Global service bridge instance
var serviceBridge ServiceBridge

// RegisterServiceBridge sets the service bridge implementation
func RegisterServiceBridge(bridge ServiceBridge) {
	serviceBridge = bridge
}

// GetServiceBridge returns the current service bridge
func GetServiceBridge() ServiceBridge {
	return serviceBridge
}
