package websocket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

var WebSocketPort = 8080

var upgrader = websocket.Upgrader{
	// Allow connections from any origin for development
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WebSocketServer represents a simple WebSocket server
type WebSocketServer struct {
	port int
}

// NewWebSocketServer creates a new WebSocket server
func NewWebSocketServer(port int) *WebSocketServer {
	return &WebSocketServer{
		port: port,
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
	defer conn.Close()

	// Log the new connection
	fmt.Println("New WebSocket connection established")

	// Send a welcome message
	err = conn.WriteMessage(websocket.TextMessage, []byte("Connected to Go WebSocket Server"))
	if err != nil {
		log.Printf("Error sending message: %v", err)
		return
	}

	// Simple message reading loop
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}

		// Log the received message
		log.Printf("Received message: %s", message)

		// Echo the message back to the client
		err = conn.WriteMessage(messageType, message)
		if err != nil {
			log.Printf("Error sending message: %v", err)
			break
		}
	}
}

func StartWebSocketServer() {
	// Start WebSocket server in a separate goroutine
	wsServer := NewWebSocketServer(WebSocketPort)
	fmt.Printf("Starting WebSocket server on port %d...\n", WebSocketPort)
	go func() {
		err := wsServer.Start()
		if err != nil {
			log.Fatalf("WebSocket server error: %v", err)
		}
	}()

	// Give the WebSocket server a moment to start
	time.Sleep(100 * time.Millisecond)
	fmt.Println("WebSocket server is running. You can now connect from the web UI.")

}
