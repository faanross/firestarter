package factory

import (
	"firestarter/internal/connections"
	"firestarter/internal/interfaces"
	"firestarter/internal/protocols/h1c"
	"firestarter/internal/protocols/h1tls"
	"firestarter/internal/protocols/h2c"
	"firestarter/internal/types"
	"fmt"
	"math/rand"
)

// AbstractFactory decides which protocol-specific factory to use
type AbstractFactory struct {
	factories   map[interfaces.ProtocolType]types.ListenerFactory
	connManager *connections.ConnectionManager // Add this field
}

// NewAbstractFactory creates a new AbstractFactory with all registered protocol factories
func NewAbstractFactory(connManager *connections.ConnectionManager) *AbstractFactory {
	return &AbstractFactory{
		factories: map[interfaces.ProtocolType]types.ListenerFactory{
			interfaces.H1C:   &h1c.Factory{},
			interfaces.H2C:   &h2c.Factory{},
			interfaces.H1TLS: &h1tls.Factory{},
			// Other protocols will be added here as they are implemented
		},
		connManager: connManager,
	}
}

// CreateListener creates a listener with the specified protocol type
func (af *AbstractFactory) CreateListener(protocol interfaces.ProtocolType, port string) (types.Listener, error) {
	factory, ok := af.factories[protocol]
	if !ok {
		return nil, fmt.Errorf("unsupported protocol: %v", protocol)
	}

	// Generate a random ID
	id := fmt.Sprintf("listener_%06d", rand.Intn(1000000))

	// Pass the connection manager along with other parameters
	return factory.CreateListener(id, port, af.connManager)
}
