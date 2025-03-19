package service

import (
	"fmt"
	"time"
)

// ConnectionTrackingUpdate performs manual verification of connection tracking
func ConnectionTrackingUpdate(listenerService *ListenerService) {
	// Print initial status
	fmt.Println("\n==== CONNECTION TRACKING TEST ====")
	fmt.Println("Initial state (should be 0 connections):")
	listenerService.LogConnectionStatus()

	// Start a monitoring goroutine
	go func() {
		for {
			time.Sleep(40 * time.Second)
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
