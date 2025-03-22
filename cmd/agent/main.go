package main

import (
	"firestarter/internal/agent/agent"
	"firestarter/internal/agent/protocol"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"firestarter/internal/agent/config"
)

func main() {
	// Load configuration
	cfg := config.DefaultConfig()
	cfg.LoadFromFlags()

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		fmt.Printf("Configuration error: %v\n", err)
		os.Exit(1)
	}

	// Display configuration
	fmt.Println(cfg)
	
	// Create the appropriate protocol based on configuration
	var proto protocol.Protocol

	switch cfg.Protocol {
	case config.H1C:
		proto = protocol.NewH1CProtocol()
	default:
		log.Fatalf("Unsupported protocol: %s", cfg.Protocol)
	}

	// Create and initialize the agent
	a := agent.NewAgent(proto)
	if err := a.Initialize(cfg); err != nil {
		log.Fatalf("Failed to initialize agent: %v", err)
	}

	// Set up signal handling for graceful shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// Start the agent
	if err := a.Start(); err != nil {
		log.Fatalf("Failed to start agent: %v", err)
	}
	log.Println("Agent started successfully")

	// Wait for termination signal
	sig := <-signalChan
	log.Printf("Received signal: %v, initiating graceful shutdown...", sig)

	// Define a timeout for graceful shutdown
	shutdownTimeout := 10 * time.Second
	log.Printf("Allowing up to %v for cleanup...", shutdownTimeout)

	// Create a timeout channel
	timeoutChan := time.After(shutdownTimeout)

	// Create a channel to signal completion of cleanup
	doneChan := make(chan struct{})

	// Perform cleanup in a goroutine
	go func() {
		// Stop the agent
		if err := a.Stop(); err != nil {
			log.Printf("Error during agent shutdown: %v", err)
		} else {
			log.Println("Agent shutdown successful")
		}

		// Report final status
		if a.IsConnected() {
			log.Println("WARNING: Agent still shows as connected after shutdown")
		} else {
			log.Println("Agent connection properly closed")
		}

		// Signal completion
		close(doneChan)
	}()

	// Wait for cleanup to complete or timeout
	select {
	case <-doneChan:
		log.Println("Graceful shutdown completed successfully")
	case <-timeoutChan:
		log.Println("Shutdown timeout exceeded, forcing exit")
	}

	log.Println("Agent exiting")
}
