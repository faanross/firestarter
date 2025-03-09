// cmd/server/main.go
package main

import (
	"firestarter/internal/factory"
	"firestarter/internal/manager"
	"firestarter/internal/service"
	"firestarter/internal/types"
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
	Protocol types.ProtocolType
}{
	{"7777", types.H1C}, // HTTP/1.1 on port 7777
	{"8888", types.H2C}, // HTTP/2 on port 8888
	{"9999", types.H2C}, // HTTP/2 on port 9999
}

func main() {
	// Setup signal channel for graceful shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	// Create the components
	abstractFactory := factory.NewAbstractFactory()
	listenerManager := manager.NewListenerManager()
	listenerService := service.NewListenerService(abstractFactory, listenerManager)

	// Wait group for synchronization
	var wg sync.WaitGroup

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

	time.Sleep(2 * time.Second)
	fmt.Printf("Managing %d active listeners.\n",
		listenerManager.Count())

	// Block until we receive a termination signal
	sig := <-signalChan
	fmt.Printf("\nReceived signal: %v. Starting graceful shutdown...\n", sig)

	// Use the service to stop all listeners
	listenerService.StopAllListeners(&wg)
}
