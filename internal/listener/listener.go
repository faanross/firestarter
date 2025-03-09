package listener

import (
	"context"
	"firestarter/internal/types"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"time"
)

// ConcreteListener represents an HTTP server instance
type ConcreteListener struct {
	ID       string
	Port     string
	Protocol types.ProtocolType
	Router   *chi.Mux
	server   *http.Server
}

func (l *ConcreteListener) Start() error {
	addr := fmt.Sprintf(":%s", l.Port)
	fmt.Printf("|START| Listener %s serving on %s\n", l.ID, addr)

	// Create the server instance
	l.server = &http.Server{
		Addr:    addr,
		Handler: l.Router,
	}

	return l.server.ListenAndServe()
}

func (l *ConcreteListener) Stop() error {
	if l.server == nil {
		return fmt.Errorf("server not started")
	}

	// Create a context with a timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Printf("|STOP| Shutting down listener %s on port %s\n", l.ID, l.Port)

	// Shutdown the server gracefully
	err := l.server.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("error shutting down listener %s: %v", l.ID, err)
	}

	fmt.Printf("|STOP| Listener %s on port %s shut down successfully\n", l.ID, l.Port)
	return nil
}

func (l *ConcreteListener) GetProtocol() string {
	switch l.Protocol {
	case types.H1C:
		return "HTTP/1.1"
	case types.H1TLS:
		return "HTTP/1.1 (TLS)"
	case types.H2C:
		return "HTTP/2"
	case types.H2TLS:
		return "HTTP/2 (TLS)"
	case types.H3:
		return "HTTP/3"
	default:
		return "Unknown"
	}
}

func (l *ConcreteListener) GetPort() string {
	return l.Port
}

func (l *ConcreteListener) GetID() string {
	return l.ID
}

// NewConcreteListener creates a new concrete listener with the specified parameters
func NewConcreteListener(id string, port string, protocol types.ProtocolType, router *chi.Mux) types.Listener {
	return &ConcreteListener{
		ID:       id,
		Port:     port,
		Protocol: protocol,
		Router:   router,
	}
}
