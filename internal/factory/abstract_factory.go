package factory

import (
	"firestarter/internal/protocols/h1c"
	"firestarter/internal/protocols/h2c"
	"firestarter/internal/types"
	"fmt"
	"math/rand"
)

// AbstractFactory decides which protocol-specific factory to use
type AbstractFactory struct {
	factories map[types.ProtocolType]types.ListenerFactory
}

// NewAbstractFactory creates a new AbstractFactory with all registered protocol factories
func NewAbstractFactory() *AbstractFactory {
	return &AbstractFactory{
		factories: map[types.ProtocolType]types.ListenerFactory{
			types.H1C: &h1c.Factory{},
			types.H2C: &h2c.Factory{},
			// Other protocols will be added here as they are implemented
			// types.H1TLS: &h1tls.Factory{},
			// types.H2TLS: &h2tls.Factory{},
			// types.H3: &h3.Factory{},
		},
	}
}

// CreateListener creates a listener with the specified protocol type
func (af *AbstractFactory) CreateListener(protocol types.ProtocolType, port string) (types.Listener, error) {
	factory, ok := af.factories[protocol]
	if !ok {
		return nil, fmt.Errorf("unsupported protocol: %v", protocol)
	}

	// Generate a random ID
	id := fmt.Sprintf("listener_%06d", rand.Intn(1000000))

	return factory.CreateListener(id, port)
}

// Convenience methods for creating specific protocol listeners
func (af *AbstractFactory) CreateH1CListener(port string) (types.Listener, error) {
	return af.CreateListener(types.H1C, port)
}

func (af *AbstractFactory) CreateH2CListener(port string) (types.Listener, error) {
	return af.CreateListener(types.H2C, port)
}
