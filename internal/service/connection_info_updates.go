package service

import (
	"firestarter/internal/interfaces"
	"fmt"
	"time"
)

// ConnectionTrackingUpdate performs manual verification of connection tracking
func ConnectionTrackingUpdate(listenerService *ListenerService, cm time.Duration) {

	fmt.Println("[🔌CON] -> Connection Passive Monitoring initialized.")
	fmt.Println("[🔌CON] -> Connection monitor updates set for:", cm)
	fmt.Println("")
	// Start a monitoring goroutine
	go func() {
		for {
			count := listenerService.GetConnectionCount()
			fmt.Println("=================>🔌CONNECTION PASSIVE MONITOR🔌<================")
			fmt.Printf("[🔌CON] -> Current Timestamp: [%s]\n", time.Now().Format(time.RFC3339))
			fmt.Printf("[🔌CON] -> %d Active Connections\n", count)
			if count > 0 {
				listenerService.LogConnectionStatus()
			}
			fmt.Println("=================================================================")
			fmt.Println()
			time.Sleep(cm)
		}

	}()
}

// LogConnectionStatus prints comprehensive connection status information
func (s *ListenerService) LogConnectionStatus() {
	currentConnections := s.connManager.GetAllConnections()

	// List all connections
	if len(currentConnections) > 0 {
		for _, conn := range currentConnections {
			fmt.Printf("[🔌CON] -> ID: %s, Protocol: %s, Created: %s\n",
				conn.GetID(),
				interfaces.GetProtocolName(conn.GetProtocol()),
				conn.GetCreatedAt().Format(time.RFC3339))
		}
	}

	fmt.Println("=================================================================")
	fmt.Println()
}
