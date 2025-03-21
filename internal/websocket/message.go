package websocket

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"strconv"
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
	case "check_port":
		// Extract the port from the payload
		payloadMap, ok := cmd.Payload.(map[string]interface{})
		if !ok {
			log.Println("[❌ERR] -> Invalid payload format for check_port command")
			return
		}

		portValue, exists := payloadMap["port"]
		if !exists {
			log.Println("[❌ERR] -> Missing 'port' in check_port payload")
			return
		}

		port, ok := portValue.(string)
		if !ok {
			log.Println("[❌ERR] -> Port must be a string")
			return
		}

		// Check if the port is available
		isAvailable := bridge.IsPortAvailable(port)

		// Send the result back to the client
		response := Message{
			Type: "port_check_result",
			Payload: map[string]interface{}{
				"port":        port,
				"isAvailable": isAvailable,
			},
		}

		// Send response directly to the requesting client, not broadcast
		err := s.sendMessage(conn, response)
		if err != nil {
			log.Printf("[❌ERR] -> Error sending port check result: %v", err)
		}
	case "create_listener":
		// Extract the parameters from the payload
		payloadMap, ok := cmd.Payload.(map[string]interface{})
		if !ok {
			log.Printf("[❌ERR] -> Invalid payload format for create_listener command")
			return
		}

		// Extract the port
		portValue, exists := payloadMap["port"]
		if !exists {
			log.Printf("[❌ERR] -> Missing 'port' in create_listener payload")
			return
		}

		// Handle port as string or number
		var port string
		switch v := portValue.(type) {
		case string:
			port = v
		case float64:
			port = fmt.Sprintf("%.0f", v)
		default:
			log.Printf("[❌ERR] -> Port must be a string or number, got %T", portValue)
			return
		}

		// Extract the protocol
		protocolValue, exists := payloadMap["protocol"]
		if !exists {
			log.Printf("[❌ERR] -> Missing 'protocol' in create_listener payload")
			return
		}

		// Protocol should be a number (though it might come as a float64 from JSON)
		var protocol int
		switch v := protocolValue.(type) {
		case float64:
			protocol = int(v)
		case string:
			// Try to parse string as int
			p, err := strconv.Atoi(v)
			if err != nil {
				log.Printf("[❌ERR] -> Protocol must be a valid number, got %v", v)
				return
			}
			protocol = p
		default:
			log.Printf("[❌ERR] -> Protocol must be a number, got %T", protocolValue)
			return
		}

		// Extract the ID (which is optional)
		var id string
		if idValue, exists := payloadMap["id"]; exists {
			if idStr, ok := idValue.(string); ok {
				id = idStr
			}
		}

		// Create the listener
		listener, err := bridge.CreateListener(id, protocol, port)
		if err != nil {
			log.Printf("[❌ERR] -> Failed to create listener: %v", err)

			// Send error message back to client
			errorResponse := Message{
				Type: "listener_creation_error",
				Payload: map[string]interface{}{
					"message": err.Error(),
				},
			}
			s.sendMessage(conn, errorResponse)
			return
		}

		// Success - send response with created listener details
		successResponse := Message{
			Type: "listener_created",
			Payload: map[string]interface{}{
				"id":       listener.GetID(),
				"port":     listener.GetPort(),
				"protocol": listener.GetProtocol(),
			},
		}
		s.sendMessage(conn, successResponse)

		fmt.Printf("[🆕NEW] -> Listener %s created successfully via WebSocket\n", listener.GetID())
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
	case "get_connections":
		return "Get Connections Snapshot"
	case "stop_connection":
		return "Stop Connection"
	case "check_port":
		return "Check Port Availability"
	case "create_listener":
		return "Create New Listener"
	default:
		return "Unknown"
	}
}
