package config

import (
	"flag"
	"fmt"
	"strings"
	"time"
)

// ProtocolType defines the type of protocol used for communication
type ProtocolType string

// Constants for supported protocol types
const (
	H1C   ProtocolType = "H1C"
	H1TLS ProtocolType = "H1TLS"
	H2C   ProtocolType = "H2C"
	H2TLS ProtocolType = "H2TLS"
	H3    ProtocolType = "H3"
)

// Config holds all configuration options for the agent
type Config struct {
	// Target server information
	TargetHost string
	TargetPort string

	// Protocol configuration
	Protocol ProtocolType

	// Connection management
	ReconnectAttempts int
	ReconnectDelay    time.Duration
	ConnectionTimeout time.Duration
	RequestTimeout    time.Duration

	// Agent identity
	AgentUUID string

	// Health check configuration
	HealthCheckInterval time.Duration
	HealthCheckEndpoint string
}

// DefaultConfig returns a Config with sensible default values
func DefaultConfig() *Config {
	return &Config{
		TargetHost:          "localhost",
		TargetPort:          "7777",
		Protocol:            H1C,
		ReconnectAttempts:   9999,             // practically indefinite, decrease if implementing backup host (TODO)
		ConnectionTimeout:   60 * time.Second, // kernel will try incremental transmissions up until 60 sec
		ReconnectDelay:      30 * time.Minute, // if not able to connect, wait 30 mins, try process again
		RequestTimeout:      5 * time.Minute,  // very generous here since unplanned timeouts can be an issue
		HealthCheckInterval: 30 * time.Second,
		HealthCheckEndpoint: "/",
	}
}

// LoadFromFlags updates the configuration based on command-line flags
func (c *Config) LoadFromFlags() {
	// Server connection flags
	flag.StringVar(&c.TargetHost, "host", c.TargetHost, "Target server hostname or IP address")
	flag.StringVar(&c.TargetPort, "port", c.TargetPort, "Target server port")

	// Protocol flag
	protocol := flag.String("protocol", string(c.Protocol), "Communication protocol (H1C, H1TLS, H2C, H2TLS, H3)")

	// Connection management flags
	flag.IntVar(&c.ReconnectAttempts, "reconnect-attempts", c.ReconnectAttempts, "Number of reconnection attempts before giving up")
	reconnectDelay := flag.Int("reconnect-delay", int(c.ReconnectDelay.Seconds()), "Delay between reconnection attempts in seconds")
	connectionTimeout := flag.Int("connection-timeout", int(c.ConnectionTimeout.Seconds()), "Connection timeout in seconds")
	requestTimeout := flag.Int("request-timeout", int(c.RequestTimeout.Seconds()), "Request timeout in seconds")

	// Health check flags
	healthCheckInterval := flag.Int("health-check-interval", int(c.HealthCheckInterval.Seconds()), "Health check interval in seconds")
	flag.StringVar(&c.HealthCheckEndpoint, "health-check-endpoint", c.HealthCheckEndpoint, "Endpoint to use for health checks")

	// Parse flags
	flag.Parse()

	// Convert string protocol to ProtocolType
	if *protocol != "" {
		protUpper := strings.ToUpper(*protocol)
		switch protUpper {
		case "H1C":
			c.Protocol = H1C
		case "H1TLS":
			c.Protocol = H1TLS
		case "H2C":
			c.Protocol = H2C
		case "H2TLS":
			c.Protocol = H2TLS
		case "H3":
			c.Protocol = H3
		default:
			fmt.Printf("Warning: Unknown protocol '%s', defaulting to H1C\n", *protocol)
			c.Protocol = H1C
		}
	}

	// Convert time values from seconds to Duration
	c.ReconnectDelay = time.Duration(*reconnectDelay) * time.Second
	c.ConnectionTimeout = time.Duration(*connectionTimeout) * time.Second
	c.RequestTimeout = time.Duration(*requestTimeout) * time.Second
	c.HealthCheckInterval = time.Duration(*healthCheckInterval) * time.Second
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.TargetHost == "" {
		return fmt.Errorf("target host cannot be empty")
	}
	if c.TargetPort == "" {
		return fmt.Errorf("target port cannot be empty")
	}
	return nil
}

// String returns a string representation of the configuration
func (c *Config) String() string {
	return fmt.Sprintf(`Agent Configuration:
  Target:                %s:%s
  Protocol:              %s
  Reconnect Attempts:    %d
  Reconnect Delay:       %v
  Connection Timeout:    %v
  Request Timeout:       %v
  Health Check Interval: %v
  Health Check Endpoint: %s`,
		c.TargetHost, c.TargetPort,
		c.Protocol,
		c.ReconnectAttempts,
		c.ReconnectDelay,
		c.ConnectionTimeout,
		c.RequestTimeout,
		c.HealthCheckInterval,
		c.HealthCheckEndpoint)
}
