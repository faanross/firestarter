package control

import (
	"firestarter/internal/factory"
	"fmt"
	"time"
)

func StopAllListeners(listeners []*factory.Listener) {
	fmt.Println("Shutting down listeners...")
	time.Sleep(1 * time.Second)

	// Gracefully stop each listener
	for _, l := range listeners {
		err := l.Stop()
		if err != nil {
			fmt.Printf("Error stopping listener %s: %v\n", l.ID, err)
		}
	}

	fmt.Println("All listeners shut down. Exiting...")
}
