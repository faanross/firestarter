package h2tls

import (
	"firestarter/internal/certificates"
	"firestarter/internal/interfaces"
	"firestarter/internal/listener"
	"firestarter/internal/router"
	"firestarter/internal/types"
	"fmt"
	"github.com/go-chi/chi/v5"
	"golang.org/x/net/http2"
	"net/http"
)

// Factory creates HTTP/2 TLS listeners
type Factory struct {
	certProvider certificates.CertificateProvider
}

// NewFactory creates a new H2TLS factory with the given certificate provider
func NewFactory(certProvider certificates.CertificateProvider) *Factory {
	return &Factory{
		certProvider: certProvider,
	}
}

// CreateListener creates and configures an HTTP/2 TLS listener
func (f *Factory) CreateListener(id string, port string, connManager interfaces.ConnectionManager) (types.Listener, error) {
	// Get TLS configuration from the certificate provider
	tlsConfig, err := f.certProvider.GetTLSConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get TLS configuration: %w", err)
	}

	// Configure specifically for HTTP/2
	// HTTP/2 requires specific TLS configuration with ALPN
	tlsConfig.NextProtos = []string{"h2", "http/1.1"}

	// Create router and set up routes
	r := chi.NewRouter()
	router.SetupRoutes(r)

	// Create the HTTP/2 server configuration
	h2srv := &http2.Server{
		// HTTP/2 specific settings
		MaxConcurrentStreams: 250,     // Limit concurrent streams per connection
		MaxReadFrameSize:     1 << 20, // 1MB max frame size
		IdleTimeout:          300,     // Seconds before idle connection is closed
	}

	// Create the concrete listener
	concreteListener := listener.NewConcreteListener(id, port, interfaces.H2TLS, r, connManager)

	// Set the TLS configuration
	concreteListener.SetTLSConfig(tlsConfig)

	// Configure the server to use HTTP/2
	// This will modify the HTTP server after it's created but before it starts
	concreteListener.SetPostServerInitFunc(func(server *http.Server) {
		// Enable HTTP/2 support on the server
		http2.ConfigureServer(server, h2srv)
	})

	fmt.Printf("|CREATE| HTTP/2 TLS Listener %s configured on port %s\n", id, port)

	return concreteListener, nil
}
