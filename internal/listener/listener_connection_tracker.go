package listener

import (
	"firestarter/internal/connections"
	"firestarter/internal/interfaces"
	"log"
	"net"
	"net/http"
)

type ConnectionTrackingListener struct {
	net.Listener
	connManager      interfaces.ConnectionManager
	protocol         interfaces.ProtocolType
	requestExtractor func(*http.Request, net.Conn)
}

// NewConnectionTrackingListener creates a new tracking listener
func NewConnectionTrackingListener(l net.Listener, cm interfaces.ConnectionManager, p interfaces.ProtocolType) *ConnectionTrackingListener {
	ctl := &ConnectionTrackingListener{
		Listener:    l,
		connManager: cm,
		protocol:    p,
	}

	// Add a request extractor function that will be called by HTTP handlers
	ctl.requestExtractor = func(req *http.Request, conn net.Conn) {
		// Extract UUID from request and associate it with this connection
		// This will be implemented in our middleware
	}

	return ctl
}

func (ctl *ConnectionTrackingListener) Accept() (net.Conn, error) {
	conn, err := ctl.Listener.Accept()
	if err != nil {
		return nil, err
	}

	var managedConn connections.Connection
	switch ctl.protocol {
	case interfaces.H1C:
		managedConn = connections.NewHTTP1Connection(conn)
	case interfaces.H2C:
		managedConn = connections.NewHTTP2Connection(conn)
	case interfaces.H1TLS:
		managedConn = connections.NewHTTP1TLSConnection(conn)
	case interfaces.H2TLS:
		managedConn = connections.NewHTTP2TLSConnection(conn)
	default:
		log.Printf("Unsupported protocol type: %v", ctl.protocol)
		return conn, nil
	}

	trackingConn := connections.NewTrackingConnection(conn, managedConn, ctl.connManager)

	return trackingConn, nil
}
