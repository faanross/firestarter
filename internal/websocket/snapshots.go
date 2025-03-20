package websocket

import (
	"github.com/gorilla/websocket"
	"log"
)

// SendListenersSnapshot sends a snapshot of all current listeners to a client
func (s *SocketServer) SendListenersSnapshot(conn *websocket.Conn) {
	// Check if we have access to the service
	bridge := GetServiceBridge()
	if bridge == nil {
		log.Println("Cannot send snapshot: service bridge not available")
		return
	}

	// Get all listeners from the service
	listeners := bridge.GetAllListeners()

	// Convert listeners to info objects
	listenerInfos := make([]ListenerInfo, 0, len(listeners))
	for _, listener := range listeners {
		listenerInfos = append(listenerInfos, ConvertListener(listener))
	}

	// Create and send the snapshot message
	snapshotMsg := Message{
		Type:    ListenersSnapshot,
		Payload: listenerInfos,
	}

	err := s.sendMessage(conn, snapshotMsg)
	if err != nil {
		log.Printf("Error sending listeners snapshot: %v", err)
	} else {
		log.Printf("Sent snapshot with %d listeners", len(listeners))
	}
}

// SendConnectionsSnapshot sends a snapshot of all current connections to a client
func (s *SocketServer) SendConnectionsSnapshot(conn *websocket.Conn) {
	// Check if we have access to the service
	bridge := GetServiceBridge()
	if bridge == nil {
		log.Println("Cannot send connection snapshot: service bridge not available")
		return
	}

	// Get all connections from the service
	connections := bridge.GetAllConnections()

	// Convert connections to info objects
	connectionInfos := make([]ConnectionInfo, 0, len(connections))
	for _, connection := range connections {
		connectionInfos = append(connectionInfos, ConvertConnection(connection))
	}

	// Create and send the snapshot message
	snapshotMsg := Message{
		Type:    ConnectionsSnapshot,
		Payload: connectionInfos,
	}

	err := s.sendMessage(conn, snapshotMsg)
	if err != nil {
		log.Printf("Error sending connections snapshot: %v", err)
	} else {
		log.Printf("Sent snapshot with %d connections", len(connections))
	}
}
