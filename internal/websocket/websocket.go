package websocket

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
	"time"
)

var upgrader = websocket.Upgrader{
	// Allow connection from any origin for development
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// SocketServer represents a WebSocket server that manages client connection
type SocketServer struct {
	port    int
	clients map[*websocket.Conn]bool
	mu      sync.Mutex
}

// GlobalWSServer is our global Singleton instance of SocketServer
var GlobalWSServer *SocketServer

// GetGlobalWSServer is the getter function for SocketServer
func GetGlobalWSServer() *SocketServer {
	return GlobalWSServer
}

// NewWebSocketServer creates a new WebSocket server
func NewWebSocketServer(port int) *SocketServer {
	return &SocketServer{
		port:    port,
		clients: make(map[*websocket.Conn]bool),
	}
}

// StartWebSocketServer initializes and starts the WebSocket server
func StartWebSocketServer(ws int) {
	fmt.Printf("\n==========>🔧CREATING WEBSOCKET SERVER🔧<==========\n")

	// Create and store global instance
	GlobalWSServer = NewWebSocketServer(ws)

	fmt.Printf("[🔧WS] -> Starting WebSocket server on port %d.\n", ws)
	go func() {
		err := GlobalWSServer.Start()
		if err != nil {
			log.Fatalf("[❌WS]-> WebSocket server error: %v", err)
		}
	}()

	// Give the WebSocket server a moment to start
	time.Sleep(100 * time.Millisecond)
	fmt.Printf("[🔧WS] -> WebSocket server is running on :%d.\n", ws)
	fmt.Println("[🖥️UI] -> You can now connect from the web UI.")
	fmt.Printf("___________________________________________________\n\n")
}

// Start begins the WebSocket server
func (s *SocketServer) Start() error {
	// Set up HTTP handler for the WebSocket endpoint
	http.HandleFunc("/ws", s.handleWebSocket)

	// Start the server
	addr := fmt.Sprintf(":%d", s.port)

	// Start the HTTP server (this is a blocking call)
	return http.ListenAndServe(addr, nil)
}

// handleWebSocket handles WebSocket connection
func (s *SocketServer) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection to WebSocket: %v", err)
		return
	}

	// Register new client
	s.mu.Lock()
	s.clients[conn] = true
	s.mu.Unlock()

	fmt.Println("New WebSocket connection established")

	// Send a welcome message
	welcomeMsg := Message{
		Type:    "welcome",
		Payload: "Connected to FirestarterC2 WebSocket Server",
	}

	err = s.sendMessage(conn, welcomeMsg)
	if err != nil {
		log.Printf("Error sending welcome message: %v", err)
	}

	// Send a snapshot of all current listeners
	s.SendListenersSnapshot(conn)

	// Send a snapshot of all current connections
	s.SendConnectionsSnapshot(conn)

	// Clean up on disconnect
	defer func() {
		conn.Close()
		s.mu.Lock()
		delete(s.clients, conn)
		s.mu.Unlock()
		fmt.Println("WebSocket connection closed")
	}()

	// Simple message reading loop
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error reading message: %v", err)
			}
			break
		}

		// Log the received message
		log.Printf("Received message: %s", message)

		// Process incoming message from ui as a command
		s.processClientMessage(conn, message)
	}
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
