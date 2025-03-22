package main

import (
	"firestarter/internal/test_agent"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	// Configure logging with timestamps
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)

	fmt.Println("Starting Firestarter Test Agent")
	fmt.Println("===============================")

	// Create a channel to listen for OS signals
	// In main(), replace the signal handling with this pattern
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	// Channel to coordinate shutdown
	shutdownChan := make(chan struct{})

	// Create a separate goroutine to handle signals
	go func() {
		<-signalChan
		fmt.Println("\n🛑 Shutdown signal received, closing all connections...")
		close(shutdownChan)
	}()

	// Wait group to track all agents
	var wg sync.WaitGroup

	// Create agents for each protocol
	agents := []test_agent.Agent{
		test_agent.NewHTTP1Agent("7676"),
		//test_agent.NewHTTP2ClearAgent("8888"),
		//test_agent.NewHTTP1TLSAgent("9999"),
		//test_agent.NewHTTP2TLSAgent("11111"),
		//test_agent.NewHTTP3Agent("22222"),
	}

	// Start each agent
	for _, agent := range agents {
		wg.Add(1)

		go func(a test_agent.Agent) {
			defer wg.Done()

			err := a.Start()
			if err != nil {
				fmt.Printf("Failed to start %s agent: %v\n", a.GetProtocol(), err)
				return
			}

			// Keep agent running until shutdown signal
			<-shutdownChan
			a.Stop()
		}(agent)

		// Small delay to stagger agent starts
		time.Sleep(500 * time.Millisecond)
	}

	// Print summary
	fmt.Println("\nAll agents started. Press Ctrl+C to terminate.")

	// Wait for all agents to complete shutdown
	wg.Wait()
	fmt.Println("🛑 All agents shut down successfully.")

}
