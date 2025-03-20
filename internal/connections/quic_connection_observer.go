package connections

import (
	"firestarter/internal/interfaces"
	"fmt"
	"github.com/quic-go/quic-go"
	"log"
)

// QuicConnectionObserver observes QUIC connection lifecycle events
type QuicConnectionObserver struct {
	connManager interfaces.ConnectionManager
}

// NewQuicConnectionObserver creates a new observer for QUIC connections
func NewQuicConnectionObserver(connManager interfaces.ConnectionManager) *QuicConnectionObserver {
	return &QuicConnectionObserver{
		connManager: connManager,
	}
}

// OnConnectionEstablished is called when a new QUIC connection is established
func (o *QuicConnectionObserver) OnConnectionEstablished(conn quic.Connection, port string) {

	fmt.Printf("[H3-DEBUG] OnConnectionEstablished called for QUIC connection from: %s\n", conn.RemoteAddr().String())

	// Create a tracked connection object
	trackedConn := NewHTTP3Connection(conn, port)

	fmt.Printf("[H3-OBSERVER-DEBUG] Created HTTP3Connection with ID: %s for protocol: %v\n", trackedConn.GetID(), trackedConn.GetProtocol())

	// Register with connection manager
	o.connManager.AddConnection(trackedConn)

	fmt.Printf("[H3-OBSERVER-DEBUG] HTTP3Connection %s registered with connection manager\n", trackedConn.GetID())

	// Set up connection close monitoring
	go o.monitorConnectionClose(conn, trackedConn.GetID())

	log.Printf("HTTP/3 connection established: %s", trackedConn.GetID())
}

// monitorConnectionClose watches for the QUIC connection to close
func (o *QuicConnectionObserver) monitorConnectionClose(conn quic.Connection, id string) {

	fmt.Printf("[H3-DEBUG] Starting to monitor QUIC connection: %s\n", id)

	// Wait for connection to close using QUIC's connregistry
	<-conn.Context().Done()

	// Deregister from connection manager
	o.connManager.RemoveConnection(id)

	log.Printf("HTTP/3 connection closed: %s", id)
}
