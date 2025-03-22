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
	// Get TLS configuration
	tlsConfig, err := f.certProvider.GetTLSConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get TLS configuration: %w", err)
	}

	// Configure for HTTP/2
	tlsConfig.NextProtos = []string{"h2", "http/1.1"}

	// Create router and set up routes
	r := chi.NewRouter()
	router.SetupRoutes(r)

	concreteListener := listener.NewConcreteListener(id, port, interfaces.H2TLS, r, connManager)

	concreteListener.SetTLSConfig(tlsConfig)

	fmt.Printf("[👂🏻LSN] -> Listener (%s) created on port %s, protocol %s\n",
		id, port, interfaces.GetProtocolName(interfaces.H2TLS))

	concreteListener.SetPostServerInitFunc(func(server *http.Server) {
		fmt.Printf("|DEBUG| Configuring HTTP/2 for server on port %s\n", port)
		http2.ConfigureServer(server, &http2.Server{})
	})

	return concreteListener, nil
}
