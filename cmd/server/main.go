package main

import (
	"firestarter/internal/connections"
	"firestarter/internal/connregistry"
	"firestarter/internal/factory"
	"firestarter/internal/manager"
	"firestarter/internal/service"
	"firestarter/internal/websocket"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var WebSocketPort = 8080

var connectionMonitor = time.Minute * 5

func main() {
	// Setup channel for SIGINT shutdown signal
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	// Setup all major server components + services
	listenerService := ApplicationSetup()

	// Wait group for synchronization
	var wg sync.WaitGroup

	time.Sleep(1 * time.Second)

	// Add connection tracking test
	service.ConnectionTrackingUpdate(listenerService, connectionMonitor)

	// Block until we receive a termination signal
	sig := <-signalChan

	// Use the service to stop all listeners
	fmt.Printf("\nReceived signal: %v. Starting graceful shutdown...\n", sig)
	listenerService.StopAllListeners(&wg)
}

func ApplicationSetup() *service.ListenerService {
	fmt.Println("===============>⚙️PERFORMING APPLICATION SETUP⚙️<===============")

	// Start our Websocket Server for UI integration
	websocket.StartWebSocketServer(WebSocketPort)

	// Create Connection Manager
	connectionManager := connections.NewConnectionManager()

	// Link the Connection Manager to the WebSocket server
	wsServer := websocket.GetGlobalWSServer()
	if wsServer != nil {
		connectionManager.SetWebSocketServer(wsServer)

	} else {
		log.Println("[❌ERR] -> WebSocket server not available for Connection Manager.")
	}

	// Initialize connection registry for UUID tracking
	connregistry.InitializeConnectionRegistry()
	connections.SetConnectionRegistry(connregistry.GetConnectionRegistry())

	// Connect the registry to the connection manager
	connregistry.ConnectRegistryToManager(connectionManager)

	af := factory.NewAbstractFactory(connectionManager)
	lm := manager.NewListenerManager()
	ls := service.NewListenerService(af, lm, connectionManager)

	// ConnectToWebSocket registers Listeners Service with WSS -> Allows UI to execute commands on server
	ls.ConnectToWebSocket()

	fmt.Println("================================================================")
	fmt.Println()
	fmt.Println("[🖥️WUI] -> YOU CAN NOW CONNECT FROM THE WEB UI <- [WUI🖥️]")
	fmt.Println()

	return ls

}
