package websocket

import (
	"firestarter/internal/types"
)

// MessageType defines the type of WebSocket messages
type MessageType string

const (
	// Message types
	ListenerCreated   MessageType = "listener_created"
	ListenerStopped   MessageType = "listener_stopped"
	ListenersSnapshot MessageType = "listeners_snapshot"
)

// ListenerInfo represents the data about a listener that will be sent to clients
type ListenerInfo struct {
	ID       string `json:"id"`
	Port     string `json:"port"`
	Protocol string `json:"protocol"`
}

// Message is the standard format for all WebSocket messages
type Message struct {
	Type    MessageType `json:"type"`
	Payload interface{} `json:"payload"`
}

// ConvertListener converts a listener to ListenerInfo format
func ConvertListener(listener types.Listener) ListenerInfo {
	return ListenerInfo{
		ID:       listener.GetID(),
		Port:     listener.GetPort(),
		Protocol: listener.GetProtocol(),
	}
}
