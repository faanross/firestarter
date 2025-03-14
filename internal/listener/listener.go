package listener

import (
	"context"
	"crypto/tls"
	"firestarter/internal/interfaces"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net"
	"net/http"
	"time"
)

// ConcreteListener represents an HTTP server instance
type ConcreteListener struct {
	ID               string
	Port             string
	Protocol         interfaces.ProtocolType
	Router           *chi.Mux
	CreatedAt        time.Time
	server           *http.Server
	handler          http.Handler
	connManager      interfaces.ConnectionManager
	tlsConfig        *tls.Config
	postServerInitFn func(*http.Server)
}

// GetCreatedAt returns time when listener was created
func (l *ConcreteListener) GetCreatedAt() time.Time {
	return l.CreatedAt
}

// SetHandler sets a custom handler
func (l *ConcreteListener) SetHandler(handler http.Handler) {
	l.handler = handler
}

// Start will run the configured listener
func (l *ConcreteListener) Start() error {
	addr := fmt.Sprintf(":%s", l.Port)

	fmt.Printf("|START| %s Listener %s serving on %s\n", l.GetProtocol(), l.ID, addr)

	// Create a standard TCP listener
	tcpListener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to create TCP listener: %v", err)
	}

	// Wrap with our connection tracking listener
	trackingListener := NewConnectionTrackingListener(tcpListener, l.connManager, l.Protocol)

	// Create the server instance
	l.server = &http.Server{
		Addr: addr,
		Handler: func() http.Handler {
			if l.handler != nil {
				return l.handler
			}
			return l.Router
		}(),
		TLSConfig: l.tlsConfig,
	}

	// Call the post-initialization function if set
	if l.postServerInitFn != nil {
		l.postServerInitFn(l.server)
	}

	// If TLS is configured, use ServeTLS, otherwise use Serve
	if l.tlsConfig != nil {
		// Empty strings for cert and key files because we're providing the certificates in TLSConfig
		return l.server.ServeTLS(trackingListener, "", "")
	} else {
		return l.server.Serve(trackingListener)
	}
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
	case interfaces.H1C:
		return "HTTP/1.1 Clear"
	case interfaces.H1TLS:
		return "HTTP/1.1 TLS"
	case interfaces.H2C:
		return "HTTP/2 Clear"
	case interfaces.H2TLS:
		return "HTTP/2 TLS"
	case interfaces.H3:
		return "HTTP/3"
	default:
		return "Unknown Protocol"
	}
}

func (l *ConcreteListener) GetPort() string {
	return l.Port
}

func (l *ConcreteListener) GetID() string {
	return l.ID
}

// NewConcreteListener constructs ConcreteListener struct
func NewConcreteListener(id string, port string, protocol interfaces.ProtocolType, router *chi.Mux, connManager interfaces.ConnectionManager) *ConcreteListener {
	return &ConcreteListener{
		ID:          id,
		Port:        port,
		Protocol:    protocol,
		Router:      router,
		CreatedAt:   time.Now(),
		handler:     nil,
		connManager: connManager,
	}
}

func (l *ConcreteListener) SetTLSConfig(config *tls.Config) {
	l.tlsConfig = config
}

// SetPostServerInitFunc sets a function that will be called after the server is created but before it starts
func (l *ConcreteListener) SetPostServerInitFunc(fn func(*http.Server)) {
	l.postServerInitFn = fn
}
