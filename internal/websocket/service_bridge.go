package websocket

import (
	"firestarter/internal/interfaces"
	"firestarter/internal/types"
)

// ServiceBridge defines the interface for accessing listener service functionality
type ServiceBridge interface {
	GetAllListeners() []types.Listener
	StopListener(id string) error
	GetAllConnections() []interfaces.Connection
	GetConnectionCount() int
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
