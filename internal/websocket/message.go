package websocket

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
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

// Message is the standard format for all WebSocket messages
type Message struct {
	Type    MessageType `json:"type"`
	Payload interface{} `json:"payload"`
}

// Command represents a request from client to server
type Command struct {
	Action  string      `json:"action"`  // What Action to perform (e.g., "stop_listener")
	Payload interface{} `json:"payload"` // Parameters for the action
}

// Broadcast sends a message to all connected clients
func (s *SocketServer) Broadcast(msg Message) {
	// Marshall the message
	jsonData, err := json.Marshal(msg)
	if err != nil {
		log.Printf("[❌ERR] -> Error marshalling message: %v.", err)
		return
	}

	// Lock before accessing clients
	s.mu.Lock()
	defer s.mu.Unlock()

	// Send to all clients
	for client := range s.clients {
		err := client.WriteMessage(websocket.TextMessage, jsonData)
		if err != nil {
			log.Printf("[❌ERR] -> Error sending message to client: %v", err)
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
		return fmt.Errorf("[❌ERR] -> Error marshalling message: %v", err)
	}

	return conn.WriteMessage(websocket.TextMessage, jsonData)
}

func (s *SocketServer) processClientMessage(conn *websocket.Conn, rawMessage []byte) {
	// Parse the incoming command
	var cmd Command
	err := json.Unmarshal(rawMessage, &cmd)
	if err != nil {
		log.Printf("[❌ERR] -> Error parsing client command: %v", err)
		return
	}

	fmt.Printf("[🖥️WUI] -> Received command: %s.\n", convertText(cmd.Action))

	// Get the service bridge
	bridge := GetServiceBridge()

	if bridge == nil {
		log.Println("[❌ERR] -> Service Bridge not available.")
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
			log.Println("[❌ERR] -> Invalid payload format for stop_listener command")
			return
		}

		idValue, exists := payloadMap["id"]
		if !exists {
			log.Println("[❌ERR] -> Missing 'id' in stop_listener payload")
			return
		}

		id, ok := idValue.(string)
		if !ok {
			log.Println("[❌ERR] -> Listener ID must be a string")
			return
		}

		// Stop the listener using the service bridge
		err := bridge.StopListener(id)
		if err != nil {
			log.Printf("[❌ERR] -> Error stopping listener %s: %v.", id, err)
		} else {
			fmt.Printf("[🛑STP] -> Listener %s stopped successfully.\n", id)
		}
	case "get_connections":
		// Send a snapshot of all connections
		s.SendConnectionsSnapshot(conn)

	case "stop_connection":
		// Extract the connection ID from the payload
		payloadMap, ok := cmd.Payload.(map[string]interface{})
		if !ok {
			log.Println("[❌ERR] -> Invalid payload format for stop_connection command")
			return
		}

		idValue, exists := payloadMap["id"]
		if !exists {
			log.Println("[❌ERR] -> Missing 'id' in stop_connection payload")
			return
		}

		id, ok := idValue.(string)
		if !ok {
			log.Println("[❌ERR] -> Connection ID must be a string")
			return
		}

		// Stop the connection using the service bridge
		err := bridge.StopConnection(id)
		if err != nil {
			log.Printf("[❌ERR] -> Error stopping connection %s: %v", id, err)
		} else {
			fmt.Printf("[🛑STP] -> Connection %s stopped successfully.\n", id)
		}

	default:
		log.Printf("[❌ERR] -> Unknown command: %s.", cmd.Action)
	}
}

func convertText(action string) string {
	switch action {
	case "get_listeners":
		return "Get Listeners Snapshot"
	case "stop_listener":
		return "Stop Listener"
	case "get_connection":
		return "Get Connections Snapshot"
	case "stop_connection":
		return "Stop Connection"
	default:
		return "Unknown"
	}
}
