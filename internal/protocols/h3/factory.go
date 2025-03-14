package h3

import (
	"firestarter/internal/certificates"
	"firestarter/internal/interfaces"
	"firestarter/internal/types"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
)

// Factory creates HTTP/3 listeners
type Factory struct {
	certProvider certificates.CertificateProvider
}

// NewFactory creates a new HTTP/3 factory with the given certificate provider
func NewFactory(certProvider certificates.CertificateProvider) *Factory {
	return &Factory{
		certProvider: certProvider,
	}
}

// CreateListener implements the ListenerFactory interface
func (f *Factory) CreateListener(id string, port string, connManager interfaces.ConnectionManager) (types.Listener, error) {
	// Verify we have a certificate provider
	if f.certProvider == nil {
		return nil, fmt.Errorf("HTTP/3 requires TLS certificates, but no certificate provider was configured")
	}

	// Get the TLS configuration
	tlsConfig, err := f.certProvider.GetTLSConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get TLS configuration for HTTP/3: %w", err)
	}

	// Configure for HTTP/3 (ALPN)
	tlsConfig.NextProtos = []string{"h3", "h3-29"}

	// Create a router and set up routes
	r := chi.NewRouter()

	// Import the router setup function from your existing code
	// This ensures all HTTP/3 routes match your other protocols
	if setupRoutes, err := importRouterSetup(); err == nil {
		setupRoutes(r)
	} else {
		return nil, fmt.Errorf("failed to set up routes for HTTP/3: %w", err)
	}

	// Create the HTTP/3 listener
	listener := NewHTTP3Listener(
		id,
		port,
		interfaces.H3,
		r,
		connManager,
	)

	// Configure TLS
	listener.SetTLSConfig(tlsConfig)

	fmt.Printf("|CREATE| HTTP/3 Listener %s configured on port %s\n", id, port)

	return listener, nil
}

// Helper function to import the router setup function
// This needs to match how your project handles routes
func importRouterSetup() (func(chi.Router), error) {
	// In a real implementation, you would return your actual router setup function
	// For now, we'll use this placeholder that you'll need to replace with your actual routing
	return func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("HTTP/3 Server"))
		})
	}, nil
}
