package main

import (
	"firestarter/internal/agent/agent"
	"firestarter/internal/agent/config"
	"firestarter/internal/agent/protocol"
	"github.com/google/uuid"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

var (
	// Build-time identity variables
	embeddedUUID string
	buildTime    string

	// Build-time protocol setting
	buildProtocol string

	// Build-time server configuration
	targetHost string
	targetPort string

	// Build-time connection management settings
	reconnectAttempts string
	reconnectDelay    string
	connectionTimeout string
	requestTimeout    string

	// Build-time health check settings
	healthCheckInterval string
	healthCheckEndpoint string
)

func main() {
	// Load default configuration
	cfg := config.DefaultConfig()

	// Apply build-time values to configuration (before flag parsing)
	applyBuildTimeConfig(cfg)

	// Load command-line flags (which will override defaults and build-time values)
	cfg.LoadFromFlags()

	// Agent identity handling
	if embeddedUUID == "" {
		// This should only happen during development
		log.Println("WARNING: No embedded UUID found. Using a temporary UUID.")
		log.Println("In production, build with: go run cmd/build/main.go")
		cfg.AgentUUID = uuid.New().String()
	} else {
		log.Printf("Agent UUID: %s (Built: %s)", embeddedUUID, buildTime)
		cfg.AgentUUID = embeddedUUID
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	// Display final configuration
	log.Println("Agent Configuration:")
	log.Println(cfg)

	// Create the appropriate protocol based on configuration
	var proto protocol.Protocol
	switch cfg.Protocol {
	case config.H1C:
		proto = protocol.NewH1CProtocol()
	case config.H1TLS:
		log.Println("Creating H1TLS protocol handler")
		proto = protocol.NewH1CProtocol() // Placeholder until H1TLS is implemented
	case config.H2C:
		log.Println("Creating H2C protocol handler")
		proto = protocol.NewH1CProtocol() // Placeholder until H2C is implemented
	case config.H2TLS:
		log.Println("Creating H2TLS protocol handler")
		proto = protocol.NewH1CProtocol() // Placeholder until H2TLS is implemented
	case config.H3:
		log.Println("Creating H3 protocol handler")
		proto = protocol.NewH1CProtocol() // Placeholder until H3 is implemented
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

	// After a short delay, test the connection
	time.Sleep(2 * time.Second)
	if err := a.TestConnection(); err != nil {
		log.Printf("Connection test failed: %v", err)
		log.Println("The agent will continue running and retry connections")
	} else {
		log.Println("Connection test successful!")
	}

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

// applyBuildTimeConfig applies build-time values to the configuration
func applyBuildTimeConfig(cfg *config.Config) {
	// Apply protocol if provided at build time
	if buildProtocol != "" {
		switch buildProtocol {
		case "h1c":
			cfg.Protocol = config.H1C
		case "h1tls":
			cfg.Protocol = config.H1TLS
		case "h2c":
			cfg.Protocol = config.H2C
		case "h2tls":
			cfg.Protocol = config.H2TLS
		case "h3":
			cfg.Protocol = config.H3
		}
	}

	// Apply server connection details
	if targetHost != "" {
		cfg.TargetHost = targetHost
	}
	if targetPort != "" {
		cfg.TargetPort = targetPort
	}

	// Apply reconnection settings
	if reconnectAttempts != "" {
		if attempts, err := strconv.Atoi(reconnectAttempts); err == nil {
			cfg.ReconnectAttempts = attempts
		} else {
			log.Printf("Warning: Invalid reconnect attempts value: %s", reconnectAttempts)
		}
	}
	if reconnectDelay != "" {
		if delay, err := time.ParseDuration(reconnectDelay); err == nil {
			cfg.ReconnectDelay = delay
		} else {
			log.Printf("Warning: Invalid reconnect delay format: %s", reconnectDelay)
		}
	}

	// Apply timeout settings
	if connectionTimeout != "" {
		if timeout, err := time.ParseDuration(connectionTimeout); err == nil {
			cfg.ConnectionTimeout = timeout
		} else {
			log.Printf("Warning: Invalid connection timeout format: %s", connectionTimeout)
		}
	}
	if requestTimeout != "" {
		if timeout, err := time.ParseDuration(requestTimeout); err == nil {
			cfg.RequestTimeout = timeout
		} else {
			log.Printf("Warning: Invalid request timeout format: %s", requestTimeout)
		}
	}

	// Apply health check settings
	if healthCheckInterval != "" {
		if interval, err := time.ParseDuration(healthCheckInterval); err == nil {
			cfg.HealthCheckInterval = interval
		} else {
			log.Printf("Warning: Invalid health check interval format: %s", healthCheckInterval)
		}
	}
	if healthCheckEndpoint != "" {
		cfg.HealthCheckEndpoint = healthCheckEndpoint
	}
}
