package main

import (
	"fmt"
	"os"

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

	// TODO: Initialize and start the agent using this configuration

	fmt.Println("Agent initialized with configuration")
}
