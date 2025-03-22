package service

import (
	"firestarter/internal/interfaces"
	"firestarter/internal/types"
	"firestarter/internal/websocket"
	"fmt"
	"sync"
)

// ConnectToWebSocket registers this service with the WebSocket server
func (s *ListenerService) ConnectToWebSocket() {
	// Create an adapter that satisfies the ServiceBridge interface
	adapter := &websocketAdapter{service: s}

	// Register the adapter with the WebSocket server
	websocket.RegisterServiceBridge(adapter)
	fmt.Println("[🔗LNK] -> Listener Service registered with WebSocket.")
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
		return fmt.Errorf(" [❌ERR] -> No connection found with ID %s", id)
	}

	// Log the termination request
	fmt.Printf("[🙋🏻REQ] -> Request to terminate connection %s (Protocol: %v, Agent: %s)\n",
		id, conn.GetProtocol(), conn.GetAgentUUID())

	// Close the connection
	err := conn.Close()
	if err != nil {
		return fmt.Errorf("[❌ERR] -> Failed to close connection %s: %w", id, err)
	}

	// Explicitly remove from connection manager to ensure proper cleanup
	a.service.GetConnectionManager().RemoveConnection(id)

	return nil
}

// IsPortAvailable implements ServiceBridge.IsPortAvailable
func (a *websocketAdapter) IsPortAvailable(port string) bool {
	return a.service.IsPortAvailable(port)
}

// CreateListener implements ServiceBridge.CreateListener
func (a *websocketAdapter) CreateListener(id string, protocol int, port string) (types.Listener, error) {
	// Convert the protocol integer to the corresponding ProtocolType
	var protocolType interfaces.ProtocolType
	switch protocol {
	case 1:
		protocolType = interfaces.H1C
	case 2:
		protocolType = interfaces.H1TLS
	case 3:
		protocolType = interfaces.H2C
	case 4:
		protocolType = interfaces.H2TLS
	case 5:
		protocolType = interfaces.H3
	default:
		return nil, fmt.Errorf("[❌ERR] -> Invalid protocol type: %d", protocol)
	}

	// Use a WaitGroup to coordinate the listener start
	var wg sync.WaitGroup

	// Pass the custom ID to the service layer
	listener, err := a.service.CreateAndStartListener(protocolType, port, &wg, id)
	if err != nil {
		return nil, fmt.Errorf("[❌ERR] -> Failed to create listener: %w", err)
	}

	return listener, nil
}
