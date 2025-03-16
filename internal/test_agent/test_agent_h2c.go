package test_agent

// HTTP2ClearAgent represents a test agent for HTTP/2 Clear
import (
	"fmt"
	"io"
	"net/http"
	"time"
)

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

	// Create a request with our UUID header
	req, err := http.NewRequest("GET", a.TargetURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Add the agent UUID as a custom header
	req.Header.Add("X-Agent-UUID", a.ID)

	// Execute the request
	resp, err := a.Client.Do(req)
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
