package h1tls

import (
	"firestarter/internal/certificates"
	"firestarter/internal/interfaces"
	"firestarter/internal/listener"
	"firestarter/internal/router"
	"firestarter/internal/types"
	"fmt"
	"github.com/go-chi/chi/v5"
)

// Factory creates HTTP/1.1 TLS listeners
type Factory struct {
	certProvider certificates.CertificateProvider
}

// NewFactory creates a new H1TLS factory with the given certificate provider
func NewFactory(certProvider certificates.CertificateProvider) *Factory {
	return &Factory{
		certProvider: certProvider,
	}
}

// CreateListener creates and configures an HTTP/1.1 TLS listener
func (f *Factory) CreateListener(id string, port string, connManager interfaces.ConnectionManager) (types.Listener, error) {
	// Get TLS configuration
	tlsConfig, err := f.certProvider.GetTLSConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get TLS configuration: %w", err)
	}

	// Create router and set up routes
	r := chi.NewRouter()
	router.SetupRoutes(r)

	// Create the listener
	concreteListener := listener.NewConcreteListener(id, port, interfaces.H1TLS, r, connManager)

	// Set the TLS configuration
	concreteListener.SetTLSConfig(tlsConfig)

	fmt.Printf("[👂🏻LSN] -> Listener (%s) created on port %s, protocol %s\n",
		id, port, interfaces.GetProtocolName(interfaces.H1TLS))

	return concreteListener, nil
}
