package service

import (
	"firestarter/internal/interfaces"
	"firestarter/internal/types"
	"firestarter/internal/websocket"
)

// ConnectToWebSocket registers this service with the WebSocket server
func (s *ListenerService) ConnectToWebSocket() {
	// Create an adapter that satisfies the ServiceBridge interface
	adapter := &websocketAdapter{service: s}

	// Register the adapter with the WebSocket server
	websocket.RegisterServiceBridge(adapter)
}

// Private adapter type that implements the ServiceBridge interface
type websocketAdapter struct {
	service *ListenerService
}

// GetAllListeners implements ServiceBridge.GetAllListeners
func (a *websocketAdapter) GetAllListeners() []types.Listener {
	return a.service.GetAllListeners()
}

// StopListener implements ServiceBridge.StopListener
func (a *websocketAdapter) StopListener(id string) error {
	return a.service.StopListener(id)
}

// GetAllConnections implements ServiceBridge.GetAllConnections
func (a *websocketAdapter) GetAllConnections() []interfaces.Connection {
	return a.service.GetAllConnections()
}

// GetConnectionCount implements ServiceBridge.GetConnectionCount
func (a *websocketAdapter) GetConnectionCount() int {
	return a.service.GetConnectionCount()
}
