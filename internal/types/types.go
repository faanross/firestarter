package types

import (
	"firestarter/internal/interfaces"
	"github.com/go-chi/chi/v5"
	"time"
)

// Listener interface defines methods that all listener types must implement
type Listener interface {
	Start() error
	Stop() error
	GetProtocol() interfaces.ProtocolType
	GetPort() string
	GetID() string
	GetCreatedAt() time.Time
}

// ListenerFactory interface defines methods for creating listeners
type ListenerFactory interface {
	CreateListener(id string, port string, connManager interfaces.ConnectionManager) (Listener, error)
}

// RouterProvider is a helper type to avoid some package cycles
type RouterProvider func() *chi.Mux
