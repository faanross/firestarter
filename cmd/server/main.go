// cmd/server/main.go
package main

import (
	"bufio"
	"firestarter/internal/connections"
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

// Define port and protocol configurations
var listenerConfigs = []struct {
	Port     string
	Protocol interfaces.ProtocolType
}{
	{"7777", interfaces.H1C}, // HTTP/1.1 on port 7777
	{"8888", interfaces.H2C}, // HTTP/2 on port 8888
	{"9999", interfaces.H2C}, // HTTP/2 on port 9999
}

func main() {
	// Setup signal channel for graceful shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	// Start our Websocket (:8080) for UI integration
	websocket.StartWebSocketServer()
	time.Sleep(5 * time.Second)

	// Create the components
	connectionManager := connections.NewConnectionManager()
	abstractFactory := factory.NewAbstractFactory(connectionManager)
	listenerManager := manager.NewListenerManager()
	listenerService := service.NewListenerService(abstractFactory, listenerManager, connectionManager)

	// ConnectToWebSocket allows our service and websocket to communicate with one another
	listenerService.ConnectToWebSocket()

	// Wait group for synchronization
	var wg sync.WaitGroup

	PressAnyKey()

	// Create and start listeners based on configurations
	for _, config := range listenerConfigs {
		time.Sleep(1 * time.Second)

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
	TestConnectionTracking(listenerService)

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

	// Optional: Clean any leftover newline characters
	fmt.Println() // Add a newline after input for cleaner output
}

// TestConnectionTracking performs manual verification of connection tracking
func TestConnectionTracking(listenerService *service.ListenerService) {
	// Print initial status
	fmt.Println("\n==== CONNECTION TRACKING TEST ====")
	fmt.Println("Initial state (should be 0 connections):")
	listenerService.LogConnectionStatus()

	// Start a monitoring goroutine
	go func() {
		for {
			time.Sleep(5 * time.Second)
			count := listenerService.GetConnectionCount()
			fmt.Printf("\n[%s] Connection monitor: %d active connections\n",
				time.Now().Format(time.RFC3339), count)
			if count > 0 {
				listenerService.LogConnectionStatus()
			}
		}
	}()

	fmt.Println("\nConnection tracking test initialized.")
	fmt.Println("Please make HTTP requests to test endpoints to verify tracking.")
	fmt.Println("====================================")
}
