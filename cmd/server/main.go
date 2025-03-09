package main

import (
	"errors"
	"firestarter/internal/control"
	"firestarter/internal/factory"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var serverPorts = []string{"7777", "8888", "9999"}

func main() {
	// Setup signal channel for graceful shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	// Create the abstract factory
	listenerFactory := factory.NewAbstractFactory()

	// Keep track of created listeners
	var listeners []factory.Listener

	// Create wait group to ensure thread sync
	var wg sync.WaitGroup

	for _, port := range serverPorts {
		time.Sleep(1 * time.Second)

		// Create a listener using the abstract factory
		l, err := listenerFactory.CreateH1CListener(port)
		if err != nil {
			fmt.Printf("Error creating service: %v\n", err)
			continue
		}

		// Store the listener
		listeners = append(listeners, l)
		time.Sleep(1 * time.Second)

		// Increment WaitGroup counter BEFORE starting the goroutine
		wg.Add(1)

		go func(l factory.Listener) {
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
