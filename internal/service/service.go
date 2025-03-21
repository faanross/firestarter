// internal/service/service.go
package service

import (
	"firestarter/internal/connections"
	"firestarter/internal/factory"
	"firestarter/internal/interfaces"
	"firestarter/internal/manager"
	"firestarter/internal/types"
	"firestarter/internal/websocket"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

// ListenerService coordinates listener lifecycle operations
type ListenerService struct {
	factory     *factory.AbstractFactory
	manager     *manager.ListenerManager
	connManager *connections.ConnectionManager
}

// NewListenerService creates a new listener service
func NewListenerService(factory *factory.AbstractFactory, manager *manager.ListenerManager, connManager *connections.ConnectionManager) *ListenerService {
	return &ListenerService{
		factory:     factory,
		manager:     manager,
		connManager: connManager,
	}
}

// CreateAndStartListener creates a listener, registers it with the manager, and starts it
func (s *ListenerService) CreateAndStartListener(protocol interfaces.ProtocolType, port string, wg *sync.WaitGroup) (types.Listener, error) {
	// Create the listener
	listener, err := s.factory.CreateListener(protocol, port)
	if err != nil {
		return nil, fmt.Errorf("failed to create listener: %w", err)
	}

	// Register with the manager
	err = s.manager.AddListener(listener)
	if err != nil {
		return nil, fmt.Errorf("failed to register listener: %w", err)
	}

	// Broadcast the creation to WebSocket clients
	wsServer := websocket.GetGlobalWSServer()
	if wsServer != nil {
		listenerInfo := websocket.ConvertListener(listener)
		wsServer.Broadcast(websocket.Message{
			Type:    websocket.ListenerCreated,
			Payload: listenerInfo,
		})
	}

	// Start the listener in a goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := listener.Start()
		if err != nil {
			// Check if this is just a server closed error (expected during shutdown)
			if err.Error() != "http: Server closed" {
				fmt.Printf("Error starting listener %s: %v\n", listener.GetID(), err)
				// Remove from manager if it failed to start unexpectedly
				_ = s.manager.RemoveListener(listener.GetID())
			}
		}
	}()

	return listener, nil
}

// StopListener stops a listener and removes it from the manager
func (s *ListenerService) StopListener(id string) error {
	// Get the listener from the manager
	listener, err := s.manager.GetListener(id)
	if err != nil {
		return err
	}

	// Broadcast the removal to WebSocket clients
	wsServer := websocket.GetGlobalWSServer()
	if wsServer != nil {
		listenerInfo := websocket.ConvertListener(listener)
		wsServer.Broadcast(websocket.Message{
			Type:    websocket.ListenerStopped,
			Payload: listenerInfo,
		})
	}

	// Stop the listener
	err = listener.Stop()
	if err != nil {
		return fmt.Errorf("failed to stop listener: %w", err)
	}

	// Remove from manager
	err = s.manager.RemoveListener(id)
	if err != nil {
		return fmt.Errorf("failed to remove listener from manager: %w", err)
	}

	return nil
}

// StopAllListeners stops all managed listeners
func (s *ListenerService) StopAllListeners(wg *sync.WaitGroup) {
	fmt.Println("Shutting down all listeners...")

	// Get all listeners
	listeners := s.manager.ListListeners()

	wsServer := websocket.GetGlobalWSServer()

	// Stop each listener
	for _, listener := range listeners {
		id := listener.GetID()

		// Broadcast the removal to WebSocket clients
		if wsServer != nil {
			listenerInfo := websocket.ConvertListener(listener)
			wsServer.Broadcast(websocket.Message{
				Type:    websocket.ListenerStopped,
				Payload: listenerInfo,
			})
		}

		err := listener.Stop()
		if err != nil {
			fmt.Printf("Error stopping listener %s: %v\n", id, err)
		}

		// Remove from manager
		_ = s.manager.RemoveListener(id)
	}

	fmt.Println("All listeners shutdown initiated. Waiting for goroutines to complete...")
	wg.Wait()
	fmt.Println("All server goroutines terminated. Exiting...")
}

// GetAllListeners returns all managed listeners
func (s *ListenerService) GetAllListeners() []types.Listener {
	return s.manager.ListListeners()
}

// GetManager returns the manager instance
func (s *ListenerService) GetManager() *manager.ListenerManager {
	return s.manager
}

// GetConnectionManager is the getter for our connection manager
func (s *ListenerService) GetConnectionManager() *connections.ConnectionManager {
	return s.connManager
}

// Change this function signature
func (s *ListenerService) GetAllConnections() []interfaces.Connection {
	return s.connManager.GetAllConnections()
}

// GetConnectionsByProtocol returns connections filtered by protocol
func (s *ListenerService) GetConnectionsByProtocol(protocol interfaces.ProtocolType) []interfaces.Connection {
	allConnections := s.connManager.GetAllConnections()
	filteredConnections := make([]interfaces.Connection, 0)

	for _, conn := range allConnections {
		if conn.GetProtocol() == protocol {
			filteredConnections = append(filteredConnections, conn)
		}
	}

	return filteredConnections
}

// GetConnectionCount returns the total number of active connections
func (s *ListenerService) GetConnectionCount() int {
	return s.connManager.Count()
}

// LogConnectionStatus prints comprehensive connection status information
func (s *ListenerService) LogConnectionStatus() {
	connections := s.connManager.GetAllConnections()
	fmt.Printf("\n==== CONNECTION STATUS REPORT ====\n")
	fmt.Printf("Total active connections: %d\n", len(connections))

	// Group by protocol
	protocolCounts := make(map[interfaces.ProtocolType]int)

	fmt.Printf("[CONN-STATUS-DEBUG] Found protocol counts: %v\n", protocolCounts)

	for _, conn := range connections {
		protocolCounts[conn.GetProtocol()]++
	}

	// Print counts by protocol with percentage
	fmt.Println("\nBreakdown by protocol:")
	for protocol, count := range protocolCounts {
		percentage := float64(count) / float64(len(connections)) * 100
		fmt.Printf("  - %s: %d connections (%.1f%%)\n",
			interfaces.GetProtocolName(protocol), count, percentage)
	}

	// List a sample of connections
	if len(connections) > 0 {
		fmt.Println("\nSample connections:")
		maxToShow := 3
		shown := 0
		for _, conn := range connections {
			if shown >= maxToShow {
				break
			}
			fmt.Printf("  - ID: %s, Protocol: %s, Created: %s\n",
				conn.GetID(),
				interfaces.GetProtocolName(conn.GetProtocol()),
				conn.GetCreatedAt().Format(time.RFC3339))
			shown++
		}

		if len(connections) > maxToShow {
			fmt.Printf("  (and %d more...)\n", len(connections)-maxToShow)
		}
	}

	fmt.Println("==================================\n")
}

// StartConnectionMonitor begins periodic monitoring of connection status
func (s *ListenerService) StartConnectionMonitor(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for {
			<-ticker.C
			connectionCount := s.GetConnectionCount()
			if connectionCount > 0 {
				fmt.Printf("\n[%s] Connection monitor:\n", time.Now().Format(time.RFC3339))
				s.LogConnectionStatus()
			}
		}
	}()
	fmt.Printf("Connection monitor started (interval: %s)\n", interval)
}

// ConnectionStats represents statistics about the active connections
type ConnectionStats struct {
	TotalConnections      int
	ConnectionsByProtocol map[interfaces.ProtocolType]int
	OldestConnection      time.Time
	NewestConnection      time.Time
	AverageAgeSeconds     float64
}

// GetConnectionStats returns statistics about the current connections
func (s *ListenerService) GetConnectionStats() ConnectionStats {
	connections := s.connManager.GetAllConnections()
	stats := ConnectionStats{
		TotalConnections:      len(connections),
		ConnectionsByProtocol: make(map[interfaces.ProtocolType]int),
	}

	if len(connections) == 0 {
		return stats
	}

	// Initialize with first connection values
	stats.OldestConnection = connections[0].GetCreatedAt()
	stats.NewestConnection = connections[0].GetCreatedAt()

	// Calculate total age for average
	var totalAgeSeconds float64

	for _, conn := range connections {
		// Update protocol counts
		stats.ConnectionsByProtocol[conn.GetProtocol()]++

		// Check for oldest/newest
		createdAt := conn.GetCreatedAt()
		if createdAt.Before(stats.OldestConnection) {
			stats.OldestConnection = createdAt
		}
		if createdAt.After(stats.NewestConnection) {
			stats.NewestConnection = createdAt
		}

		// Add to total age
		ageSeconds := time.Since(createdAt).Seconds()
		totalAgeSeconds += ageSeconds
	}

	// Calculate average age
	stats.AverageAgeSeconds = totalAgeSeconds / float64(len(connections))

	return stats
}

// BroadcastConnectionStatus sends connection statistics to all WebSocket clients
func (s *ListenerService) BroadcastConnectionStatus() {
	wsServer := websocket.GetGlobalWSServer()
	if wsServer == nil {
		return
	}

	stats := s.GetConnectionStats()

	// Create a simplified structure for the WebSocket message
	type connectionStatusPayload struct {
		TotalConnections  int            `json:"totalConnections"`
		ByProtocol        map[string]int `json:"byProtocol"`
		AverageAgeSeconds float64        `json:"averageAgeSeconds"`
	}

	// Convert protocol type keys to strings for JSON
	byProtocolStr := make(map[string]int)
	for proto, count := range stats.ConnectionsByProtocol {
		byProtocolStr[interfaces.GetProtocolName(proto)] = count
	}

	payload := connectionStatusPayload{
		TotalConnections:  stats.TotalConnections,
		ByProtocol:        byProtocolStr,
		AverageAgeSeconds: stats.AverageAgeSeconds,
	}

	// Define a new message type for connection status
	const ConnectionStatus websocket.MessageType = "connection_status"

	// Broadcast the message
	wsServer.Broadcast(websocket.Message{
		Type:    ConnectionStatus,
		Payload: payload,
	})
}

// IsPortAvailable checks if the specified port is available for binding
func (s *ListenerService) IsPortAvailable(port string) bool {
	// Try to bind to the port to see if it's available
	listener, err := net.Listen("tcp", ":"+port)

	// If there was an error, the port is not available
	if err != nil {
		log.Printf("[❌ERR] -> Port %s is not available: %v", port, err)
		return false
	}

	// If we get here, the port is available, so close the listener and return true
	listener.Close()
	fmt.Printf("[✅SCS] -> Port %s is available\n", port)
	return true
}
