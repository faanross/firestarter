package control

import (
	"firestarter/internal/factory"
	"fmt"
	"sync"
	"time"
)

func StopAllListeners(listeners []*factory.Listener, wg *sync.WaitGroup) {
	fmt.Println("Shutting down listeners...")

	// Gracefully stop each listener
	for _, l := range listeners {
		time.Sleep(1 * time.Second)
		err := l.Stop()
		if err != nil {
			fmt.Printf("Error stopping listener %s: %v\n", l.ID, err)
		}
	}

	fmt.Println("All listeners shutdown initiated. Waiting for server goroutines to complete...")

	// Wait for all server goroutines to finish
	wg.Wait()

	fmt.Println("All server goroutines terminated. Exiting...")
}
