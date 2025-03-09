package factory

import (
	"fmt"
	"math/rand"
)

// AbstractFactory decides which protocol-specific factory to use
type AbstractFactory struct {
	factories map[ProtocolType]ListenerFactory
}

// NewAbstractFactory creates a new AbstractFactory with all registered protocol factories
func NewAbstractFactory() *AbstractFactory {
	return &AbstractFactory{
		factories: map[ProtocolType]ListenerFactory{
			H1C: &H1CFactory{},
			// Add other protocol factories as they are implemented
			// H2C: &H2CFactory{},
		},
	}
}

// CreateListener creates a listener with the specified protocol type
func (af *AbstractFactory) CreateListener(protocol ProtocolType, port string) (Listener, error) {
	factory, ok := af.factories[protocol]
	if !ok {
		return nil, fmt.Errorf("unsupported protocol: %v", protocol)
	}

	// Generate a random ID
	id := fmt.Sprintf("listener_%06d", rand.Intn(1000000))

	return factory.CreateListener(id, port)
}

// CreateH1CListener is a convenience method for creating H1C listeners
func (af *AbstractFactory) CreateH1CListener(port string) (Listener, error) {
	return af.CreateListener(H1C, port)
}
