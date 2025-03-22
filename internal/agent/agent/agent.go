package agent

import (
	"context"
	"firestarter/internal/agent/config"
	"firestarter/internal/agent/protocol"
	"fmt"
	"log"
	"sync"
	"time"
)

// Agent represents the core agent functionality
type Agent struct {
	// Configuration
	config *config.Config

	// Protocol implementation to use for communication
	protocol protocol.Protocol

	// Agent state tracking
	running      bool
	runningLock  sync.RWMutex
	stopChan     chan struct{}
	healthTicker *time.Ticker

	// Error tracking
	lastError     error
	lastErrorLock sync.RWMutex

	// Connection attempt tracking
	connectionAttempts int
}

// NewAgent creates a new agent instance with the specified protocol
func NewAgent(protocol protocol.Protocol) *Agent {
	return &Agent{
		protocol:           protocol,
		connectionAttempts: 0,
		stopChan:           make(chan struct{}),
	}
}

// Initialize sets up the agent with the provided configuration
func (a *Agent) Initialize(cfg *config.Config) error {
	log.Printf("Initializing agent with %s protocol", cfg.Protocol)

	a.config = cfg

	// Convert from agent config to protocol config
	protocolCfg := protocol.ProtocolConfig{
		TargetHost:          cfg.TargetHost,
		TargetPort:          cfg.TargetPort,
		AgentUUID:           cfg.AgentUUID,
		ConnectionTimeout:   cfg.ConnectionTimeout,
		RequestTimeout:      cfg.RequestTimeout,
		HealthCheckEndpoint: cfg.HealthCheckEndpoint,
	}

	// Initialize the protocol
	err := a.protocol.Initialize(protocolCfg)
	if err != nil {
		return fmt.Errorf("failed to initialize protocol: %w", err)
	}

	return nil
}

// Start begins agent operations, establishing a connection and starting health checks
func (a *Agent) Start() error {
	// Prevent starting twice
	if a.isRunning() {
		return fmt.Errorf("agent is already running")
	}

	log.Printf("Starting agent, targeting %s:%s using %s protocol",
		a.config.TargetHost, a.config.TargetPort, a.config.Protocol)

	// Attempt initial connection
	if err := a.connect(); err != nil {
		log.Printf("Initial connection failed: %v", err)
		// We don't return an error here, as the agent will keep trying to connect
	}

	// Mark as running and set up health checks
	a.setRunning(true)
	a.healthTicker = time.NewTicker(a.config.HealthCheckInterval)

	// Start health check goroutine
	go a.healthCheckLoop()

	return nil
}

// Stop gracefully shuts down the agent
func (a *Agent) Stop() error {
	if !a.isRunning() {
		return nil // Already stopped
	}

	log.Println("Stopping agent...")

	// Signal health check loop to stop
	close(a.stopChan)

	// Stop the ticker if it exists
	if a.healthTicker != nil {
		a.healthTicker.Stop()
	}

	// Disconnect from server
	if a.protocol.IsConnected() {
		if err := a.protocol.Disconnect(); err != nil {
			log.Printf("Error disconnecting: %v", err)
			// Continue with shutdown anyway
		}
	}

	// Mark as not running
	a.setRunning(false)
	log.Println("Agent stopped")

	return nil
}

// connect attempts to establish a connection to the server
func (a *Agent) connect() error {
	log.Println("Attempting to connect to server...")

	// Reset connection attempts if we were previously connected
	if a.protocol.IsConnected() {
		a.connectionAttempts = 0
	}

	// Create a context with the connection timeout
	ctx, cancel := context.WithTimeout(context.Background(), a.config.ConnectionTimeout)
	defer cancel()

	// Attempt to connect
	err := a.protocol.Connect(ctx)
	if err != nil {
		a.connectionAttempts++
		a.setLastError(err)
		log.Printf("Connection attempt %d failed: %v", a.connectionAttempts, err)
		return err
	}

	// Success
	log.Println("Successfully connected to server")
	a.connectionAttempts = 0
	a.setLastError(nil)
	return nil
}

// reconnect implements the reconnection logic with exponential backoff
func (a *Agent) reconnect() {
	// Check if we've exceeded max attempts
	if a.config.ReconnectAttempts > 0 && a.connectionAttempts >= a.config.ReconnectAttempts {
		log.Printf("Exceeded maximum reconnection attempts (%d), giving up", a.config.ReconnectAttempts)
		return
	}

	// Calculate delay with exponential backoff and jitter
	delay := a.config.ReconnectDelay
	for i := 0; i < a.connectionAttempts && i < 8; i++ {
		delay *= 2
	}

	// Add jitter (±20%)
	jitterFactor := 0.8 + (0.4 * float64(time.Now().Nanosecond()) / float64(1000000000))
	delay = time.Duration(float64(delay) * jitterFactor)

	// Cap the maximum delay
	maxDelay := 1 * time.Hour
	if delay > maxDelay {
		delay = maxDelay
	}

	log.Printf("Waiting %v before reconnection attempt %d", delay, a.connectionAttempts+1)

	// Wait for the delay or until the agent is stopped
	select {
	case <-time.After(delay):
		// Try to connect again
		_ = a.connect() // Ignore error, we're already handling it
	case <-a.stopChan:
		// Agent is stopping, abort reconnect
		return
	}
}

// healthCheckLoop runs periodic health checks in a separate goroutine
func (a *Agent) healthCheckLoop() {
	log.Printf("Starting health check loop with interval: %v", a.config.HealthCheckInterval)

	for {
		select {
		case <-a.healthTicker.C:
			// Skip if not running
			if !a.isRunning() {
				continue
			}

			// If not connected, try to reconnect
			if !a.protocol.IsConnected() {
				a.reconnect()
				continue
			}

			// Perform health check
			ctx, cancel := context.WithTimeout(context.Background(), a.config.RequestTimeout)
			err := a.protocol.PerformHealthCheck(ctx)
			cancel()

			if err != nil {
				log.Printf("Health check failed: %v", err)
				// If health check fails, we need to reconnect
				a.reconnect()
			} else {
				log.Printf("Health check successful")
			}

		case <-a.stopChan:
			// Agent is stopping
			log.Println("Health check loop terminating")
			return
		}
	}
}

// SendRequest sends a request to the server and returns the response
func (a *Agent) SendRequest(endpoint string, payload []byte) ([]byte, error) {
	if !a.isRunning() {
		return nil, fmt.Errorf("agent is not running")
	}

	if !a.protocol.IsConnected() {
		return nil, fmt.Errorf("agent is not connected to server")
	}

	ctx, cancel := context.WithTimeout(context.Background(), a.config.RequestTimeout)
	defer cancel()

	return a.protocol.SendRequest(ctx, endpoint, payload)
}

// IsConnected returns whether the agent is currently connected to the server
func (a *Agent) IsConnected() bool {
	return a.isRunning() && a.protocol.IsConnected()
}

// GetLastError returns the last error encountered by the agent
func (a *Agent) GetLastError() error {
	a.lastErrorLock.RLock()
	defer a.lastErrorLock.RUnlock()
	return a.lastError
}

// Helper method to safely check if the agent is running
func (a *Agent) isRunning() bool {
	a.runningLock.RLock()
	defer a.runningLock.RUnlock()
	return a.running
}

// Helper method to safely set the running state
func (a *Agent) setRunning(running bool) {
	a.runningLock.Lock()
	defer a.runningLock.Unlock()
	a.running = running
}

// Helper method to safely set the last error
func (a *Agent) setLastError(err error) {
	a.lastErrorLock.Lock()
	defer a.lastErrorLock.Unlock()
	a.lastError = err
}
