package websocket

import (
	"encoding/json"
	"firestarter/internal/connregistry"
	"firestarter/internal/interfaces"
	"firestarter/internal/types"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

// MessageType defines the type of WebSocket messages
type MessageType string

const (
	ListenerCreated     MessageType = "listener_created"
	ListenerStopped     MessageType = "listener_stopped"
	ListenersSnapshot   MessageType = "listeners_snapshot"
	ConnectionCreated   MessageType = "connection_created"
	ConnectionStopped   MessageType = "connection_stopped"
	ConnectionsSnapshot MessageType = "connections_snapshot"
)

// ListenerInfo represents the data about a listener that will be sent to UI
type ListenerInfo struct {
	ID        string    `json:"id"`
	Port      string    `json:"port"`
	Protocol  string    `json:"protocol"`
	CreatedAt time.Time `json:"createdAt"`
}

// Message is the standard format for all WebSocket messages
type Message struct {
	Type    MessageType `json:"type"`
	Payload interface{} `json:"payload"`
}

// ConvertListener converts a listener to ListenerInfo format
func ConvertListener(listener types.Listener) ListenerInfo {
	return ListenerInfo{
		ID:        listener.GetID(),
		Port:      listener.GetPort(),
		Protocol:  listener.GetProtocol(),
		CreatedAt: listener.GetCreatedAt(),
	}
}

// Command represents a request from client to server
type Command struct {
	Action  string      `json:"action"`  // What action to perform (e.g., "stop_listener")
	Payload interface{} `json:"payload"` // Parameters for the action
}

// ConnectionInfo represents the data about a connection that will be sent to UI
type ConnectionInfo struct {
	ID         string    `json:"id"`         // Unique identifier for the connection
	Port       string    `json:"port"`       // Port the connection is using
	Protocol   string    `json:"protocol"`   // Protocol type (H1C, H1TLS, etc.)
	CreatedAt  time.Time `json:"createdAt"`  // When the connection was established
	RemoteAddr string    `json:"remoteAddr"` // Client IP address and port
	AgentUUID  string    `json:"agentUUID"`  // UUID of the connected agent
}

// ConvertConnection converts a connection to ConnectionInfo format
func ConvertConnection(conn interfaces.Connection) ConnectionInfo {
	return ConnectionInfo{
		ID:         conn.GetID(),
		Port:       conn.GetPort(),
		Protocol:   interfaces.GetProtocolName(conn.GetProtocol()),
		CreatedAt:  conn.GetCreatedAt(),
		RemoteAddr: getRemoteAddrFromConnection(conn),
		AgentUUID:  conn.GetAgentUUID(),
	}
}

// Helper function to get remote address if available
func getRemoteAddrFromConnection(conn interfaces.Connection) string {
	// Get the registry using the existing getter function
	registry := connregistry.GetConnectionRegistry()

	if registry != nil {
		connID := conn.GetID()
		if remoteAddr := registry.GetRemoteAddrByConnID(connID); remoteAddr != "" {
			return remoteAddr
		}
	}

	// Fallback methods
	if httpConn, ok := conn.(interface{ GetRemoteAddr() string }); ok {
		return httpConn.GetRemoteAddr()
	}

	return "Unknown"
}

// Broadcast sends a message to all connected clients
func (s *SocketServer) Broadcast(msg Message) {
	// Marshall the message
	jsonData, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshalling message: %v", err)
		return
	}

	// Lock before accessing clients
	s.mu.Lock()
	defer s.mu.Unlock()

	// Send to all clients
	for client := range s.clients {
		err := client.WriteMessage(websocket.TextMessage, jsonData)
		if err != nil {
			log.Printf("Error sending message to client: %v", err)
			// Client may be disconnected, clean up
			client.Close()
			delete(s.clients, client)
		}
	}
}

// sendMessage sends a message to a specific client
func (s *SocketServer) sendMessage(conn *websocket.Conn, msg Message) error {
	jsonData, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("error marshalling message: %v", err)
	}

	return conn.WriteMessage(websocket.TextMessage, jsonData)
}

func (s *SocketServer) processClientMessage(conn *websocket.Conn, rawMessage []byte) {
	// Parse the incoming command
	var cmd Command
	err := json.Unmarshal(rawMessage, &cmd)
	if err != nil {
		log.Printf("Error parsing client command: %v", err)
		return
	}

	log.Printf("Received command: %s", cmd.Action)

	// Get the service bridge
	bridge := GetServiceBridge()
	if bridge == nil {
		log.Println("Cannot process command: service bridge not available")
		return
	}

	// Handle different command types
	switch cmd.Action {
	case "get_listeners":
		// Send a snapshot of all listeners
		s.SendListenersSnapshot(conn)

	case "stop_listener":
		// Extract the listener ID from the payload
		payloadMap, ok := cmd.Payload.(map[string]interface{})
		if !ok {
			log.Println("Invalid payload format for stop_listener command")
			return
		}

		idValue, exists := payloadMap["id"]
		if !exists {
			log.Println("Missing 'id' in stop_listener payload")
			return
		}

		id, ok := idValue.(string)
		if !ok {
			log.Println("Listener ID must be a string")
			return
		}

		// Stop the listener using the service bridge
		err := bridge.StopListener(id)
		if err != nil {
			log.Printf("Error stopping listener %s: %v", id, err)
		} else {
			log.Printf("Listener %s stopped successfully", id)
		}
	case "get_connections":
		// Send a snapshot of all connections
		s.SendConnectionsSnapshot(conn)

	case "stop_connection":
		// Extract the connection ID from the payload
		payloadMap, ok := cmd.Payload.(map[string]interface{})
		if !ok {
			log.Println("Invalid payload format for stop_connection command")
			return
		}

		idValue, exists := payloadMap["id"]
		if !exists {
			log.Println("Missing 'id' in stop_connection payload")
			return
		}

		id, ok := idValue.(string)
		if !ok {
			log.Println("Connection ID must be a string")
			return
		}

		// Stop the connection using the service bridge
		err := bridge.StopConnection(id)
		if err != nil {
			log.Printf("Error stopping connection %s: %v", id, err)
		} else {
			log.Printf("Connection %s stopped successfully", id)
		}

	default:
		log.Printf("Unknown command: %s", cmd.Action)
	}
}
