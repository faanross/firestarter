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

	listenerFactory := factory.NewListenerFactory()
	// Keep track of created listeners
	var listeners []*factory.Listener

	// create wait group to ensure thread sync
	var wg sync.WaitGroup
	
	for _, port := range serverPorts {
		time.Sleep(1 * time.Second)
		l, err := listenerFactory.CreateListener(port)
		if err != nil {
			fmt.Printf("Error creating service: %v\n", err)
			continue
		}
		// Store the listener
		listeners = append(listeners, l)

		time.Sleep(1 * time.Second)
		go func(l *factory.Listener) {
			err := l.Start()
			if err != nil && !errors.Is(err, http.ErrServerClosed) {
				fmt.Printf("Error starting listener %s: %v\n", l.ID, err)
			}
		}(l)
	}

	time.Sleep(2 * time.Second)

	// Block until we receive a termination signal
	sig := <-signalChan
	fmt.Printf("\nReceived signal: %v. Starting graceful shutdown...\n", sig)
	control.StopAllListeners(listeners)

}
