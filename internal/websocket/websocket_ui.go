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

var WebSocketPort = 8080

var upgrader = websocket.Upgrader{
	// Allow connections from any origin for development
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WebSocketServer represents a WebSocket server that manages client connections
type WebSocketServer struct {
	port    int
	clients map[*websocket.Conn]bool
	mu      sync.Mutex // For thread safety when accessing clients
}

// Global instance of WebSocketServer to be accessed from other packages
var GlobalWSServer *WebSocketServer

// GetGlobalWSServer returns the global WebSocket server instance
func GetGlobalWSServer() *WebSocketServer {
	return GlobalWSServer
}

// NewWebSocketServer creates a new WebSocket server
func NewWebSocketServer(port int) *WebSocketServer {
	return &WebSocketServer{
		port:    port,
		clients: make(map[*websocket.Conn]bool),
	}
}

// Start begins the WebSocket server
func (s *WebSocketServer) Start() error {
	// Set up HTTP handler for the WebSocket endpoint
	http.HandleFunc("/ws", s.handleWebSocket)

	// Start the server
	addr := fmt.Sprintf(":%d", s.port)
	fmt.Printf("WebSocket server starting on %s\n", addr)

	// Start the HTTP server (this is a blocking call)
	return http.ListenAndServe(addr, nil)
}

// handleWebSocket handles WebSocket connections
func (s *WebSocketServer) handleWebSocket(w http.ResponseWriter, r *http.Request) {
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
func (s *WebSocketServer) Broadcast(msg Message) {
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
func (s *WebSocketServer) sendMessage(conn *websocket.Conn, msg Message) error {
	jsonData, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("error marshalling message: %v", err)
	}

	return conn.WriteMessage(websocket.TextMessage, jsonData)
}

// StartWebSocketServer initializes and starts the WebSocket server
func StartWebSocketServer() {
	// Create and store global instance
	GlobalWSServer = NewWebSocketServer(WebSocketPort)

	fmt.Printf("Starting WebSocket server on port %d...\n", WebSocketPort)
	go func() {
		err := GlobalWSServer.Start()
		if err != nil {
			log.Fatalf("WebSocket server error: %v", err)
		}
	}()

	// Give the WebSocket server a moment to start
	time.Sleep(100 * time.Millisecond)
	fmt.Println("WebSocket server is running. You can now connect from the web UI.")
}

func (s *WebSocketServer) processClientMessage(conn *websocket.Conn, rawMessage []byte) {
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

	default:
		log.Printf("Unknown command: %s", cmd.Action)
	}
}

// SendListenersSnapshot sends a snapshot of all current listeners to a client
func (s *WebSocketServer) SendListenersSnapshot(conn *websocket.Conn) {
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
