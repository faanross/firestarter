package factory

import (
	"firestarter/internal/certificates"
	"firestarter/internal/connections"
	"firestarter/internal/interfaces"
	"firestarter/internal/protocols/h1c"
	"firestarter/internal/protocols/h1tls"
	"firestarter/internal/protocols/h2c"
	"firestarter/internal/protocols/h2tls"
	"firestarter/internal/protocols/h3"
	"firestarter/internal/types"
	"fmt"
	"math/rand"
)

// AbstractFactory decides which protocol-specific factory to use
type AbstractFactory struct {
	factories   map[interfaces.ProtocolType]types.ListenerFactory
	connManager *connections.ConnectionManager
}

// NewAbstractFactory creates a new AbstractFactory with all registered protocol factories
func NewAbstractFactory(connManager *connections.ConnectionManager) *AbstractFactory {
	certProvider, err := certificates.GetDefaultCertificateProvider()
	if err != nil {
		fmt.Printf("[⚠️WRN] -> Failed to load certificates: %v\n", err)
		fmt.Println("[⚠️WRN] -> TLS listeners will not be available.")
		// Continue without TLS support
		return &AbstractFactory{
			factories: map[interfaces.ProtocolType]types.ListenerFactory{
				interfaces.H1C: &h1c.Factory{},
				interfaces.H2C: &h2c.Factory{},
			},
			connManager: connManager,
		}
	}

	// Create factories with certProvider for h1tls, h2tls, and h3
	h1tlsFactory := h1tls.NewFactory(certProvider)
	h2tlsFactory := h2tls.NewFactory(certProvider)
	h3Factory := h3.NewFactory(certProvider)

	fmt.Println("[🏭ABS] -> Loaded AbstractFactory with certificates.")
	fmt.Println("[🏭ABS] -> All protocols available as listeners.")

	return &AbstractFactory{
		factories: map[interfaces.ProtocolType]types.ListenerFactory{
			interfaces.H1C:   &h1c.Factory{},
			interfaces.H2C:   &h2c.Factory{},
			interfaces.H1TLS: h1tlsFactory,
			interfaces.H2TLS: h2tlsFactory,
			interfaces.H3:    h3Factory,
		},
		connManager: connManager,
	}
}

// CreateListener creates a listener with the specified protocol type
func (af *AbstractFactory) CreateListener(protocol interfaces.ProtocolType, port string, customID string) (types.Listener, error) {
	factory, ok := af.factories[protocol]
	if !ok {
		return nil, fmt.Errorf("[❌ERR] -> Unsupported protocol: %v", protocol)
	}

	// Use custom ID if provided, otherwise generate a random one
	id := customID
	if id == "" {
		id = fmt.Sprintf("listener_%06d", rand.Intn(1000000))
	}

	// Pass the connection manager along with other parameters
	return factory.CreateListener(id, port, af.connManager)
}
