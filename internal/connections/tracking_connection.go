package connections

import (
	"firestarter/internal/interfaces"
	"fmt"
	"net"
	"time"
)

// TrackingConnection wraps a standard net.Conn and handles tracking lifecycle
type TrackingConnection struct {
	// The actual network connection
	conn net.Conn

	// Reference to the connection manager
	manager interfaces.ConnectionManager

	// The tracked connection object
	trackedConn Connection

	// Flag to prevent double-close
	closed bool
}

// NewTrackingConnection creates a new connection that manages its own tracking lifecycle
// NewTrackingConnection creates a connection that manages its own tracking lifecycle
func NewTrackingConnection(conn net.Conn, trackedConn interfaces.Connection,
	manager interfaces.ConnectionManager) *TrackingConnection {
	tc := &TrackingConnection{
		conn:        conn,
		manager:     manager,
		trackedConn: trackedConn,
		closed:      false,
	}

	// Register this connection with the global registry
	connectionRegistry.RegisterConnection(conn, trackedConn.GetID())

	// Register with the connection manager (but UUID will be set later)
	manager.AddConnection(trackedConn)

	return tc
}

// Implement net.Conn interface by delegating to the wrapped connection
func (tc *TrackingConnection) Read(b []byte) (n int, err error) {
	return tc.conn.Read(b)
}

func (tc *TrackingConnection) Write(b []byte) (n int, err error) {
	return tc.conn.Write(b)
}

func (tc *TrackingConnection) Close() error {
	// Prevent double-close and ensure cleanup happens only once
	if tc.closed {
		return nil
	}

	// Mark as closed
	tc.closed = true

	// Remove from connection manager
	fmt.Printf("Connection closed: %s\n", tc.trackedConn.GetID())
	tc.manager.RemoveConnection(tc.trackedConn.GetID())

	// Close the underlying connection
	return tc.conn.Close()
}

func (tc *TrackingConnection) LocalAddr() net.Addr {
	return tc.conn.LocalAddr()
}

func (tc *TrackingConnection) RemoteAddr() net.Addr {
	return tc.conn.RemoteAddr()
}

func (tc *TrackingConnection) SetDeadline(t time.Time) error {
	return tc.conn.SetDeadline(t)
}

func (tc *TrackingConnection) SetReadDeadline(t time.Time) error {
	return tc.conn.SetReadDeadline(t)
}

func (tc *TrackingConnection) SetWriteDeadline(t time.Time) error {
	return tc.conn.SetWriteDeadline(t)
}

// UpdateAgentUUID updates the agent UUID for this connection
func (tc *TrackingConnection) UpdateAgentUUID(agentUUID string) {
	// This is a simplification - in real code we'd need to access the underlying
	// connection and set its AgentUUID field
	if conn, ok := tc.trackedConn.(interface{ SetAgentUUID(string) }); ok {
		conn.SetAgentUUID(agentUUID)
	}
}
