package test_agent

import (
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/http"
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
