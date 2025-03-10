package types

import (
	"github.com/go-chi/chi/v5"
	"time"
)

// ProtocolType defines the supported protocol types
type ProtocolType int

const (
	H1C ProtocolType = iota + 1
	H1TLS
	H2C
	H2TLS
	H3
)

// Listener interface defines methods that all listener types must implement
type Listener interface {
	Start() error
	Stop() error
	GetProtocol() string
	GetPort() string
	GetID() string
	GetCreatedAt() time.Time
}

// ListenerFactory interface defines methods for creating listeners
type ListenerFactory interface {
	CreateListener(id string, port string) (Listener, error)
}

// RouterProvider is a helper type to avoid some package cycles
type RouterProvider func() *chi.Mux
