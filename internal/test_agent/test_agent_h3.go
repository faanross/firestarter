package test_agent

import (
	"crypto/tls"
	"fmt"
	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
	"io"
	"net/http"
	"time"
)

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
