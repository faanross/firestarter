package test_agent

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"time"
)

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
