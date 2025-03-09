package factory

type ProtocolType int

const (
	H1C ProtocolType = iota + 1
	//H1TLS
	H2C
	//H2TLS
	//H3
)

// Listener interface defines methods that all listener types must implement
type Listener interface {
	Start() error
	Stop() error
	GetProtocol() string
	GetPort() string
	GetID() string
}

// ListenerFactory interface defines methods for creating listeners
type ListenerFactory interface {
	CreateListener(id string, port string) (Listener, error)
}
