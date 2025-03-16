package main

import (
	"crypto/tls"
	"fmt"
	"github.com/google/uuid"
	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// Agent defines the common behavior for all test agents
type Agent interface {
	Start() error
	Stop() error
	RunHealthCheck() error
	GetID() string
	GetProtocol() string
}

// BaseAgent provides common functionality for all agents
type BaseAgent struct {
	ID        string
	Protocol  string
	TargetURL string
	Client    *http.Client
	IsRunning bool
	StopChan  chan struct{}
}

// Initialize sets up the base agent
func (a *BaseAgent) Initialize(protocol string, port string) {
	a.ID = uuid.New().String()
	a.Protocol = protocol
	a.StopChan = make(chan struct{})
	a.IsRunning = false

	// Set the target URL based on protocol
	scheme := "http"
	if protocol == "H1TLS" || protocol == "H2TLS" || protocol == "H3" {
		scheme = "https"
	}
	a.TargetURL = fmt.Sprintf("%s://localhost:%s", scheme, port)

	// Log agent creation
	log.Printf("| AGENT %s | Created with ID: %s | Target: %s", protocol, a.ID, a.TargetURL)
}

// GetID returns the agent's unique ID
func (a *BaseAgent) GetID() string {
	return a.ID
}

// GetProtocol returns the agent's protocol
func (a *BaseAgent) GetProtocol() string {
	return a.Protocol
}

// Log prints a formatted log message with the agent identifier
func (a *BaseAgent) Log(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	log.Printf("| AGENT %s | %s", a.Protocol, message)
}

// HTTP1Agent represents a test agent for HTTP/1.1 Clear
type HTTP1Agent struct {
	BaseAgent
}

// NewHTTP1Agent creates a new HTTP/1.1 Clear agent
func NewHTTP1Agent(port string) *HTTP1Agent {
	agent := &HTTP1Agent{}
	agent.Initialize("H1C", port)

	// Configure HTTP/1.1 specific settings
	agent.Client = &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:       1,
			MaxConnsPerHost:    1,
			IdleConnTimeout:    90 * time.Second,
			DisableCompression: true,
			DisableKeepAlives:  false, // Keep connections alive
			ForceAttemptHTTP2:  false, // Force HTTP/1.1
		},
	}

	return agent
}

// Start begins the agent operations
func (a *HTTP1Agent) Start() error {
	a.Log("Starting agent...")

	// Perform initial connection
	err := a.RunHealthCheck()
	if err != nil {
		a.Log("Failed initial connection: %v", err)
		return err
	}

	a.IsRunning = true

	// Start background health check loop
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := a.RunHealthCheck(); err != nil {
					a.Log("Health check failed: %v", err)
				}
			case <-a.StopChan:
				a.Log("Background health checks stopped")
				return
			}
		}
	}()

	a.Log("Agent successfully started")
	return nil
}

// Stop terminates the agent operations
func (a *HTTP1Agent) Stop() error {
	if !a.IsRunning {
		return nil
	}

	a.Log("Stopping agent...")
	close(a.StopChan)
	a.IsRunning = false
	a.Log("Agent stopped")
	return nil
}

// RunHealthCheck performs a connection check
func (a *HTTP1Agent) RunHealthCheck() error {
	a.Log("Performing health check...")

	resp, err := a.Client.Get(a.TargetURL)
	if err != nil {
		return fmt.Errorf("connection failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1024)) // Limit to first 1KB
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	a.Log("Health check successful (Status: %s, Body: %s)", resp.Status, string(body))
	return nil
}

// HTTP2ClearAgent represents a test agent for HTTP/2 Clear
type HTTP2ClearAgent struct {
	BaseAgent
}

// NewHTTP2ClearAgent creates a new HTTP/2 Clear agent
func NewHTTP2ClearAgent(port string) *HTTP2ClearAgent {
	agent := &HTTP2ClearAgent{}
	agent.Initialize("H2C", port)

	// Configure HTTP/2 specific settings
	agent.Client = &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:       1,
			MaxConnsPerHost:    1,
			IdleConnTimeout:    90 * time.Second,
			DisableCompression: true,
			DisableKeepAlives:  false, // Keep connections alive
			ForceAttemptHTTP2:  true,  // Force HTTP/2
		},
	}

	return agent
}

// Start begins the agent operations
func (a *HTTP2ClearAgent) Start() error {
	a.Log("Starting agent...")

	// Perform initial connection
	err := a.RunHealthCheck()
	if err != nil {
		a.Log("Failed initial connection: %v", err)
		return err
	}

	a.IsRunning = true

	// Start background health check loop
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := a.RunHealthCheck(); err != nil {
					a.Log("Health check failed: %v", err)
				}
			case <-a.StopChan:
				a.Log("Background health checks stopped")
				return
			}
		}
	}()

	a.Log("Agent successfully started")
	return nil
}

// Stop terminates the agent operations
func (a *HTTP2ClearAgent) Stop() error {
	if !a.IsRunning {
		return nil
	}

	a.Log("Stopping agent...")
	close(a.StopChan)
	a.IsRunning = false
	a.Log("Agent stopped")
	return nil
}

// RunHealthCheck performs a connection check
func (a *HTTP2ClearAgent) RunHealthCheck() error {
	a.Log("Performing health check...")

	resp, err := a.Client.Get(a.TargetURL)
	if err != nil {
		return fmt.Errorf("connection failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1024)) // Limit to first 1KB
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	a.Log("Health check successful (Status: %s, Body: %s, Protocol: %s)",
		resp.Status, string(body), resp.Proto)
	return nil
}

// HTTP1TLSAgent represents a test agent for HTTP/1.1 TLS
type HTTP1TLSAgent struct {
	BaseAgent
}

// NewHTTP1TLSAgent creates a new HTTP/1.1 TLS agent
func NewHTTP1TLSAgent(port string) *HTTP1TLSAgent {
	agent := &HTTP1TLSAgent{}
	agent.Initialize("H1TLS", port)

	// Configure HTTP/1.1 TLS specific settings
	agent.Client = &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:       1,
			MaxConnsPerHost:    1,
			IdleConnTimeout:    90 * time.Second,
			DisableCompression: true,
			DisableKeepAlives:  false, // Keep connections alive
			ForceAttemptHTTP2:  false, // Force HTTP/1.1
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // Skip certificate validation for testing
			},
		},
	}

	return agent
}

// Start begins the agent operations
func (a *HTTP1TLSAgent) Start() error {
	a.Log("Starting agent...")

	// Perform initial connection
	err := a.RunHealthCheck()
	if err != nil {
		a.Log("Failed initial connection: %v", err)
		return err
	}

	a.IsRunning = true

	// Start background health check loop
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := a.RunHealthCheck(); err != nil {
					a.Log("Health check failed: %v", err)
				}
			case <-a.StopChan:
				a.Log("Background health checks stopped")
				return
			}
		}
	}()

	a.Log("Agent successfully started")
	return nil
}

// Stop terminates the agent operations
func (a *HTTP1TLSAgent) Stop() error {
	if !a.IsRunning {
		return nil
	}

	a.Log("Stopping agent...")
	close(a.StopChan)
	a.IsRunning = false
	a.Log("Agent stopped")
	return nil
}

// RunHealthCheck performs a connection check
func (a *HTTP1TLSAgent) RunHealthCheck() error {
	a.Log("Performing health check...")

	resp, err := a.Client.Get(a.TargetURL)
	if err != nil {
		return fmt.Errorf("connection failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1024)) // Limit to first 1KB
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	a.Log("Health check successful (Status: %s, Body: %s, TLS: %v)",
		resp.Status, string(body), resp.TLS != nil)
	return nil
}

// HTTP2TLSAgent represents a test agent for HTTP/2 TLS
type HTTP2TLSAgent struct {
	BaseAgent
}

// NewHTTP2TLSAgent creates a new HTTP/2 TLS agent
func NewHTTP2TLSAgent(port string) *HTTP2TLSAgent {
	agent := &HTTP2TLSAgent{}
	agent.Initialize("H2TLS", port)

	// Configure HTTP/2 TLS specific settings
	agent.Client = &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:       1,
			MaxConnsPerHost:    1,
			IdleConnTimeout:    90 * time.Second,
			DisableCompression: true,
			DisableKeepAlives:  false, // Keep connections alive
			ForceAttemptHTTP2:  true,  // Force HTTP/2
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // Skip certificate validation for testing
			},
		},
	}

	return agent
}

// Start begins the agent operations
func (a *HTTP2TLSAgent) Start() error {
	a.Log("Starting agent...")

	// Perform initial connection
	err := a.RunHealthCheck()
	if err != nil {
		a.Log("Failed initial connection: %v", err)
		return err
	}

	a.IsRunning = true

	// Start background health check loop
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := a.RunHealthCheck(); err != nil {
					a.Log("Health check failed: %v", err)
				}
			case <-a.StopChan:
				a.Log("Background health checks stopped")
				return
			}
		}
	}()

	a.Log("Agent successfully started")
	return nil
}

// Stop terminates the agent operations
func (a *HTTP2TLSAgent) Stop() error {
	if !a.IsRunning {
		return nil
	}

	a.Log("Stopping agent...")
	close(a.StopChan)
	a.IsRunning = false
	a.Log("Agent stopped")
	return nil
}

// RunHealthCheck performs a connection check
func (a *HTTP2TLSAgent) RunHealthCheck() error {
	a.Log("Performing health check...")

	resp, err := a.Client.Get(a.TargetURL)
	if err != nil {
		return fmt.Errorf("connection failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1024)) // Limit to first 1KB
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	// Check specifically for HTTP/2
	isHTTP2 := resp.ProtoMajor == 2

	a.Log("Health check successful (Status: %s, Body: %s, HTTP/2: %v)",
		resp.Status, string(body), isHTTP2)
	return nil
}

// HTTP3Agent represents a test agent for HTTP/3
type HTTP3Agent struct {
	BaseAgent
	quicClient *http3.Transport
}

// NewHTTP3Agent creates a new HTTP/3 agent
func NewHTTP3Agent(port string) *HTTP3Agent {
	agent := &HTTP3Agent{}
	agent.Initialize("H3", port)

	// Configure HTTP/3 specific settings
	qconf := &quic.Config{
		MaxIdleTimeout: 30 * time.Second,
	}

	// Create the HTTP/3 transport
	transport := &http3.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // Skip certificate validation for testing
		},
		QUICConfig:      qconf,
		EnableDatagrams: false,
	}

	// Store the transport
	agent.quicClient = transport

	// Create an HTTP client using the HTTP/3 transport
	agent.Client = &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}

	return agent
}

// Start begins the agent operations
func (a *HTTP3Agent) Start() error {
	a.Log("Starting agent...")

	// Perform initial connection
	err := a.RunHealthCheck()
	if err != nil {
		a.Log("Failed initial connection: %v", err)
		return err
	}

	a.IsRunning = true

	// Start background health check loop
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := a.RunHealthCheck(); err != nil {
					a.Log("Health check failed: %v", err)
				}
			case <-a.StopChan:
				a.Log("Background health checks stopped")
				return
			}
		}
	}()

	a.Log("Agent successfully started")
	return nil
}

// Stop terminates the agent operations
func (a *HTTP3Agent) Stop() error {
	if !a.IsRunning {
		return nil
	}

	a.Log("Stopping agent...")
	close(a.StopChan)
	a.IsRunning = false
	a.Log("Agent stopped")
	return nil
}

// RunHealthCheck performs a connection check using HTTP/3
func (a *HTTP3Agent) RunHealthCheck() error {
	a.Log("Performing health check...")

	// Use the regular client which now has HTTP/3 transport
	resp, err := a.Client.Get(a.TargetURL)
	if err != nil {
		return fmt.Errorf("HTTP/3 connection failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body with limit
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1024)) // Limit to first 1KB
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	a.Log("Health check successful (Status: %s, Body: %s)",
		resp.Status, string(body))

	fmt.Printf("[H3-CLIENT-DEBUG] HTTP/3 connection successful to %s (Proto: %s)\n", a.TargetURL, resp.Proto)
	
	return nil
}

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
	agents := []Agent{
		NewHTTP1Agent("7777"),
		NewHTTP2ClearAgent("8888"),
		NewHTTP1TLSAgent("9999"),
		NewHTTP2TLSAgent("11111"),
		NewHTTP3Agent("22222"),
	}

	// Start each agent
	for _, agent := range agents {
		wg.Add(1)

		go func(a Agent) {
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
