package connections

import (
	"firestarter/internal/interfaces"
	"net"
	"time"
)

// HTTP1TLSConnection represents an HTTP/1 TLS specific connection
type HTTP1TLSConnection struct {
	BaseConnection
	Conn net.Conn
}

// NewHTTP1TLSConnection creates a new HTTP/1.1 TLS connection
func NewHTTP1TLSConnection(conn net.Conn) *HTTP1TLSConnection {
	return &HTTP1TLSConnection{
		BaseConnection: BaseConnection{
			ID:        GenerateUniqueID(),
			Protocol:  interfaces.H1TLS,
			CreatedAt: time.Now().UTC(),
		},
		Conn: conn,
	}
}

// Implement the Connection interface
func (c *HTTP1TLSConnection) GetID() string { return c.ID }
func (c *HTTP1TLSConnection) GetProtocol() interfaces.ProtocolType {
	return c.Protocol
}
func (c *HTTP1TLSConnection) GetCreatedAt() time.Time { return c.CreatedAt }
func (c *HTTP1TLSConnection) GetPort() string         { return c.Port }
func (c *HTTP1TLSConnection) Close() error            { return c.Conn.Close() }
