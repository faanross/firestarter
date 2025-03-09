package main

import (
	"errors"
	"firestarter/internal/control"
	"firestarter/internal/factory"
	"firestarter/internal/types"
	"fmt"
	"net/http"
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
}

func main() {
	// Setup signal channel for graceful shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	// Create the abstract factory
	listenerFactory := factory.NewAbstractFactory()

	// Keep track of created listeners
	var listeners []types.Listener

	// Create wait group to ensure thread sync
	var wg sync.WaitGroup

	// Create and start listeners based on configurations
	for _, config := range listenerConfigs {
		time.Sleep(1 * time.Second)

		// Create a listener using the abstract factory
		l, err := listenerFactory.CreateListener(config.Protocol, config.Port)
		if err != nil {
			fmt.Printf("Error creating service: %v\n", err)
			continue
		}

		// Log the protocol being used
		fmt.Printf("Created %s listener on port %s\n", l.GetProtocol(), config.Port)

		// Store the listener
		listeners = append(listeners, l)
		time.Sleep(1 * time.Second)

		// Increment WaitGroup counter before starting the goroutine
		wg.Add(1)

		go func(l types.Listener) {
			// Defer the Done() call so it happens even if there's an error
			defer wg.Done()
			err := l.Start()
			if err != nil && !errors.Is(err, http.ErrServerClosed) {
				fmt.Printf("Error starting listener %s: %v\n", l.GetID(), err)
			}
		}(l)
	}

	time.Sleep(2 * time.Second)

	// Block until we receive a termination signal
	sig := <-signalChan
	fmt.Printf("\nReceived signal: %v. Starting graceful shutdown...\n", sig)

	// Stop all listeners
	control.StopAllListeners(listeners, &wg)
}
