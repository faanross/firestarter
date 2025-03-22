package protocol

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

// H1CProtocol implements the Protocol interface for HTTP/1.1 Clear (H1C)
type H1CProtocol struct {
	// Configuration
	config ProtocolConfig

	// HTTP client for connection
	client *http.Client

	// Connection state
	connected     bool
	connectedLock sync.RWMutex

	// Activity tracking
	lastActivity     time.Time
	lastActivityLock sync.RWMutex
}

// NewH1CProtocol creates a new instance of the H1C protocol
func NewH1CProtocol() *H1CProtocol {
	return &H1CProtocol{
		lastActivity: time.Now(),
	}
}

// Initialize sets up the H1C protocol with the provided configuration
func (p *H1CProtocol) Initialize(config ProtocolConfig) error {
	p.config = config

	// Create the HTTP client with appropriate timeouts
	p.client = &http.Client{
		Timeout: config.RequestTimeout,
		Transport: &http.Transport{
			MaxIdleConns:        1,
			IdleConnTimeout:     90 * time.Second,
			DisableCompression:  true,
			MaxConnsPerHost:     1,
			ForceAttemptHTTP2:   false, // Ensure HTTP/1.1 is used
			TLSHandshakeTimeout: config.ConnectionTimeout,
		},
	}

	return nil
}

// Connect establishes a connection to the server
func (p *H1CProtocol) Connect(ctx context.Context) error {
	// Create a simple GET request to check if the server is reachable
	targetURL := fmt.Sprintf("http://%s:%s%s",
		p.config.TargetHost,
		p.config.TargetPort,
		p.config.HealthCheckEndpoint)

	req, err := http.NewRequestWithContext(ctx, "GET", targetURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Add the agent UUID to the request
	req.Header.Set("X-Agent-UUID", p.config.AgentUUID)

	// Send the request
	resp, err := p.client.Do(req)
	if err != nil {
		p.setConnected(false)
		return fmt.Errorf("connection failed: %w", err)
	}
	defer resp.Body.Close()

	// Check if the response is successful (2xx status code)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		p.setConnected(false)
		return fmt.Errorf("server returned non-success status: %d", resp.StatusCode)
	}

	// Update connection status and last activity
	p.setConnected(true)
	p.updateLastActivity()

	return nil
}

// Disconnect terminates the connection to the server
func (p *H1CProtocol) Disconnect() error {
	// HTTP is stateless, so we just mark ourselves as disconnected
	// The actual connection handling is done by the http.Client
	p.setConnected(false)
	return nil
}

// IsConnected returns whether the connection is currently active
func (p *H1CProtocol) IsConnected() bool {
	p.connectedLock.RLock()
	defer p.connectedLock.RUnlock()
	return p.connected
}

// SendRequest sends a request to the server and returns the response
func (p *H1CProtocol) SendRequest(ctx context.Context, endpoint string, payload []byte) ([]byte, error) {
	// Ensure we're connected
	if !p.IsConnected() {
		return nil, fmt.Errorf("not connected to server")
	}

	// Build the full URL
	targetURL := fmt.Sprintf("http://%s:%s%s",
		p.config.TargetHost,
		p.config.TargetPort,
		endpoint)

	// Create the request
	req, err := http.NewRequestWithContext(ctx, "POST", targetURL, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("X-Agent-UUID", p.config.AgentUUID)

	// Send the request
	resp, err := p.client.Do(req)
	if err != nil {
		p.setConnected(false)
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check if the response is successful
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("server returned non-success status: %d", resp.StatusCode)
	}

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Update last activity
	p.updateLastActivity()

	return respBody, nil
}

// PerformHealthCheck conducts a health check against the server
func (p *H1CProtocol) PerformHealthCheck(ctx context.Context) error {
	// Similar to Connect, but we just check if the server is reachable
	targetURL := fmt.Sprintf("http://%s:%s%s",
		p.config.TargetHost,
		p.config.TargetPort,
		p.config.HealthCheckEndpoint)

	req, err := http.NewRequestWithContext(ctx, "GET", targetURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create health check request: %w", err)
	}

	// Add the agent UUID
	req.Header.Set("X-Agent-UUID", p.config.AgentUUID)

	// Send the request
	resp, err := p.client.Do(req)
	if err != nil {
		p.setConnected(false)
		return fmt.Errorf("health check failed: %w", err)
	}
	defer resp.Body.Close()

	// Check if the response is successful
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		p.setConnected(false)
		return fmt.Errorf("health check returned non-success status: %d", resp.StatusCode)
	}

	// Update connection status and last activity
	p.setConnected(true)
	p.updateLastActivity()

	return nil
}

// GetLastActivity returns the time of the last successful communication
func (p *H1CProtocol) GetLastActivity() time.Time {
	p.lastActivityLock.RLock()
	defer p.lastActivityLock.RUnlock()
	return p.lastActivity
}

// Name returns the name of the protocol
func (p *H1CProtocol) Name() string {
	return "H1C"
}

// Helper method to update the last activity time
func (p *H1CProtocol) updateLastActivity() {
	p.lastActivityLock.Lock()
	defer p.lastActivityLock.Unlock()
	p.lastActivity = time.Now()
}

// Helper method to update the connection status
func (p *H1CProtocol) setConnected(connected bool) {
	p.connectedLock.Lock()
	defer p.connectedLock.Unlock()
	p.connected = connected
}
