package websocket

import (
	"firestarter/internal/interfaces"
	"firestarter/internal/types"
	"time"
)

// MessageType defines the type of WebSocket messages
type MessageType string

const (
	ListenerCreated     MessageType = "listener_created"
	ListenerStopped     MessageType = "listener_stopped"
	ListenersSnapshot   MessageType = "listeners_snapshot"
	ConnectionCreated   MessageType = "connection_created"
	ConnectionStopped   MessageType = "connection_stopped"
	ConnectionsSnapshot MessageType = "connections_snapshot"
)

// ListenerInfo represents the data about a listener that will be sent to UI
type ListenerInfo struct {
	ID        string    `json:"id"`
	Port      string    `json:"port"`
	Protocol  string    `json:"protocol"`
	CreatedAt time.Time `json:"createdAt"`
}

// Message is the standard format for all WebSocket messages
type Message struct {
	Type    MessageType `json:"type"`
	Payload interface{} `json:"payload"`
}

// ConvertListener converts a listener to ListenerInfo format
func ConvertListener(listener types.Listener) ListenerInfo {
	return ListenerInfo{
		ID:        listener.GetID(),
		Port:      listener.GetPort(),
		Protocol:  listener.GetProtocol(),
		CreatedAt: listener.GetCreatedAt(),
	}
}

// Command represents a request from client to server
type Command struct {
	Action  string      `json:"action"`  // What action to perform (e.g., "stop_listener")
	Payload interface{} `json:"payload"` // Parameters for the action
}

// ConnectionInfo represents the data about a connection that will be sent to UI
type ConnectionInfo struct {
	ID         string    `json:"id"`
	Port       string    `json:"port"`
	Protocol   string    `json:"protocol"`
	CreatedAt  time.Time `json:"createdAt"`
	RemoteAddr string    `json:"remoteAddr"`
	AgentUUID  string    `json:"agentUUID"`
}

// ConvertConnection converts a connection to ConnectionInfo format
func ConvertConnection(conn interfaces.Connection) ConnectionInfo {
	return ConnectionInfo{
		ID:   conn.GetID(),
		Port: conn.GetPort(),
		//Protocol:   getProtocolName(conn.GetProtocol()),
		CreatedAt:  conn.GetCreatedAt(),
		RemoteAddr: getRemoteAddrFromConnection(conn),
		AgentUUID:  conn.GetAgentUUID(),
	}
}

// Helper function to get remote address if available
func getRemoteAddrFromConnection(conn interfaces.Connection) string {
	// This is a placeholder, I will later implement ability to
	// retrieve the remote address from the connection
	return "Unknown"
}
