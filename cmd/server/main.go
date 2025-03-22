package main

import (
	"firestarter/internal/connections"
	"firestarter/internal/connregistry"
	"firestarter/internal/factory"
	"firestarter/internal/manager"
	"firestarter/internal/service"
	"firestarter/internal/websocket"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var WebSocketPort = 8080

func main() {
	// Setup channel for SIGINT shutdown signal
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	// Setup all major server components + services
	listenerManager, listenerService := ApplicationSetup()

	// Wait group for synchronization
	var wg sync.WaitGroup

	time.Sleep(1 * time.Second)
	fmt.Printf("Managing %d active listeners.\n",
		listenerManager.Count())

	// Add connection tracking test
	service.ConnectionTrackingUpdate(listenerService)

	// Block until we receive a termination signal
	sig := <-signalChan

	fmt.Printf("\nReceived signal: %v. Starting graceful shutdown...\n", sig)

	// Use the service to stop all listeners
	listenerService.StopAllListeners(&wg)
}

func ApplicationSetup() (*manager.ListenerManager, *service.ListenerService) {
	// Start our Websocket Server for UI integration
	websocket.StartWebSocketServer(WebSocketPort)

	// Create Connection Manager
	connectionManager := connections.NewConnectionManager()

	// Link the Connection Manager to the WebSocket server
	wsServer := websocket.GetGlobalWSServer()
	if wsServer != nil {
		connectionManager.SetWebSocketServer(wsServer)

	} else {
		fmt.Println("[INIT-ERROR] WebSocket server not available for Connection Manager!")
	}

	// Initialize connection registry for UUID tracking
	connregistry.InitializeConnectionRegistry()
	connections.SetConnectionRegistry(connregistry.GetConnectionRegistry())

	// Connect the registry to the connection manager
	connregistry.ConnectRegistryToManager(connectionManager)

	af := factory.NewAbstractFactory(connectionManager)
	lm := manager.NewListenerManager()
	ls := service.NewListenerService(af, lm, connectionManager)

	// ConnectToWebSocket allows our service and websocket to communicate with one another
	ls.ConnectToWebSocket()

	return lm, ls

}
