package service

import (
	"firestarter/internal/interfaces"
	"firestarter/internal/types"
	"firestarter/internal/websocket"
	"fmt"
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

// StopConnection implements ServiceBridge.StopConnection
func (a *websocketAdapter) StopConnection(id string) error {
	// Find the connection in the connection manager
	conn, found := a.service.GetConnectionManager().GetConnection(id)
	if !found {
		return fmt.Errorf("no connection found with ID %s", id)
	}

	// Log the termination request
	fmt.Printf("Request to terminate connection %s (Protocol: %v, Agent: %s)\n",
		id, conn.GetProtocol(), conn.GetAgentUUID())

	// Close the connection
	err := conn.Close()
	if err != nil {
		return fmt.Errorf("failed to close connection %s: %w", id, err)
	}

	// Explicitly remove from connection manager to ensure proper cleanup
	a.service.GetConnectionManager().RemoveConnection(id)

	return nil
}

// IsPortAvailable implements ServiceBridge.IsPortAvailable
func (a *websocketAdapter) IsPortAvailable(port string) bool {
	return a.service.IsPortAvailable(port)
}
