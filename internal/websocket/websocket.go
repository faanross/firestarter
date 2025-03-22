package websocket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
	"time"
)

// TODO restrict access before release for prod
var upgrader = websocket.Upgrader{
	// Allow connection from any origin for development
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// GlobalWSServer is our global Singleton instance of SocketServer
var GlobalWSServer *SocketServer

// GetGlobalWSServer is the getter function for SocketServer
func GetGlobalWSServer() *SocketServer {
	return GlobalWSServer
}

// SocketServer represents a WebSocket server that manages UI client connections
type SocketServer struct {
	port    int
	clients map[*websocket.Conn]bool // allow for multiple UI in future
	mu      sync.Mutex
}

// NewWebSocketServer is constructor for WebSocket server
func NewWebSocketServer(port int) *SocketServer {
	return &SocketServer{
		port:    port,
		clients: make(map[*websocket.Conn]bool),
	}
}

// StartWebSocketServer initializes and starts the WebSocket server
func StartWebSocketServer(wsp int) {

	// Create and store global instance
	GlobalWSServer = NewWebSocketServer(wsp)

	go func() {
		err := GlobalWSServer.Start()
		if err != nil {
			log.Fatalf("[❌ERR]-> WebSocket server error: %v", err)
		}
	}()

	// Give the WebSocket server a moment to start
	time.Sleep(100 * time.Millisecond)
	fmt.Printf("[🔧WSS] -> WebSocket server is running on :%d.\n", wsp)

}

// Start begins the WebSocket server
func (s *SocketServer) Start() error {
	// Set up HTTP handler for the WebSocket endpoint
	http.HandleFunc("/ws", s.handleWebSocket)

	// Start the server
	addr := fmt.Sprintf(":%d", s.port)

	fmt.Printf("[🔧WSS] -> Starting WebSocket server on port %d.\n", s.port)

	// Start the HTTP server (this is a blocking call)
	return http.ListenAndServe(addr, nil)
}

// handleWebSocket handler that manages the full lifecycle of a client connection
func (s *SocketServer) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Printf("[❌ERR] -> Failed to upgrade connection to WebSocket: %v.", err)
		return
	}

	// Register new client
	s.mu.Lock()
	s.clients[conn] = true
	s.mu.Unlock()

	fmt.Printf("[🔧WSS] -> WebSocket connection established - %s.\n", time.Now().Format(time.RFC3339))

	// Clean up on disconnect
	defer func() {
		conn.Close()
		s.mu.Lock()
		delete(s.clients, conn)
		s.mu.Unlock()
		fmt.Println("[🔧WSS] -> WebSocket connection closed.")
	}()

	// Simple message reading loop
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("[❌ERR] -> Error reading message: %v", err)
			}
			break
		}

		// Process incoming message from ui as a command
		s.processClientMessage(conn, message)
	}
}
