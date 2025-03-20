package main

import (
	"bufio"
	"firestarter/internal/connections"
	"firestarter/internal/connregistry"
	"firestarter/internal/factory"
	"firestarter/internal/interfaces"
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

var listenerConfigs = []struct {
	Port     string
	Protocol interfaces.ProtocolType
}{
	{"7777", interfaces.H1C},    // HTTP/1.1 on port 7777
	{"8888", interfaces.H2C},    // HTTP/2 on port 8888
	{"9999", interfaces.H1TLS},  // HTTP/1.1 TLS on port 9999
	{"11111", interfaces.H2TLS}, // HTTP/2 TLS on port 11111
	{"22222", interfaces.H3},    // HTTP/3  on port 22222
}

func main() {
	// Setup channel for SIGINT shutdown signal
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	// Setup all major server components + services
	listenerManager, listenerService := ApplicationSetup()

	// Wait group for synchronization
	var wg sync.WaitGroup

	PressAnyKey()

	// Create and start listeners based on configurations
	for _, config := range listenerConfigs {
		time.Sleep(300 * time.Millisecond)

		// Use the service to create and start the listener
		_, err := listenerService.CreateAndStartListener(
			config.Protocol,
			config.Port,
			&wg,
		)

		if err != nil {
			fmt.Printf("Error creating and starting listener: %v\n", err)
			continue
		}

	}

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

// PressAnyKey displays a message and waits for the user to press any key before continuing
func PressAnyKey() {
	fmt.Println("Press any key to start creating listeners...")

	// Create a reader to read a single byte from stdin
	reader := bufio.NewReader(os.Stdin)
	_, _ = reader.ReadString('\n')

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
