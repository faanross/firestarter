package connections

import (
	"firestarter/internal/interfaces"
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
func (o *QuicConnectionObserver) OnConnectionEstablished(conn quic.Connection) {
	// Create a tracked connection object
	trackedConn := NewHTTP3Connection(conn)

	// Register with connection manager
	o.connManager.AddConnection(trackedConn)

	// Set up connection close monitoring
	go o.monitorConnectionClose(conn, trackedConn.GetID())

	log.Printf("HTTP/3 connection established: %s", trackedConn.GetID())
}

// monitorConnectionClose watches for the QUIC connection to close
func (o *QuicConnectionObserver) monitorConnectionClose(conn quic.Connection, id string) {
	// Wait for connection to close using QUIC's context
	<-conn.Context().Done()

	// Deregister from connection manager
	o.connManager.RemoveConnection(id)

	log.Printf("HTTP/3 connection closed: %s", id)
}
