package main

import (
	"errors"
	"firestarter/internal/factory"
	"fmt"
	"net/http"
	"time"
)

var serverPorts = []string{"7777", "8888", "9999"}

func main() {

	listenerFactory := factory.NewListenerFactory()
	// Keep track of created listeners
	var listeners []*factory.Listener

	for _, port := range serverPorts {
		time.Sleep(2 * time.Second)
		l, err := listenerFactory.CreateListener(port)
		if err != nil {
			fmt.Printf("Error creating service: %v\n", err)
			continue
		}
		// Store the listener
		listeners = append(listeners, l)

		time.Sleep(2 * time.Second)
		go func(l *factory.Listener) {
			err := l.Start()
			if err != nil && !errors.Is(err, http.ErrServerClosed) {
				fmt.Printf("Error starting listener %s: %v\n", l.ID, err)
			}
		}(l)
	}

	time.Sleep(2 * time.Second)
	TestListenerStop(listeners)

}

func TestListenerStop(listeners []*factory.Listener) {
	// Wait for 30 seconds before starting the graceful shutdown
	fmt.Println("All listeners started. Will begin shutdown in 15 seconds...")
	time.Sleep(15 * time.Second)

	// Gracefully stop each listener
	for _, l := range listeners {
		err := l.Stop()
		if err != nil {
			fmt.Printf("Error stopping listener %s: %v\n", l.ID, err)
		}
		time.Sleep(2 * time.Second) // Small delay between stopping listeners
	}

	fmt.Println("All listeners shut down. Exiting...")
}
