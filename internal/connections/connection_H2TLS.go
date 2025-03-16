package connections

import (
	"firestarter/internal/interfaces"
	"net"
	"time"
)

// HTTP2TLSConnection represents an HTTP/2 TLS specific connection
type HTTP2TLSConnection struct {
	BaseConnection
	Conn net.Conn
}

// NewHTTP2TLSConnection creates a new HTTP/2 TLS connection
func NewHTTP2TLSConnection(conn net.Conn) *HTTP2TLSConnection {
	return &HTTP2TLSConnection{
		BaseConnection: BaseConnection{
			ID:        GenerateUniqueID(),
			Protocol:  interfaces.H2TLS,
			CreatedAt: time.Now().UTC(),
		},
		Conn: conn,
	}
}

// Connection interface implementation
func (c *HTTP2TLSConnection) GetID() string { return c.ID }
func (c *HTTP2TLSConnection) GetProtocol() interfaces.ProtocolType {
	return c.Protocol
}
func (c *HTTP2TLSConnection) GetCreatedAt() time.Time { return c.CreatedAt }
func (c *HTTP2TLSConnection) GetPort() string         { return c.Port }
func (c *HTTP2TLSConnection) Close() error            { return c.Conn.Close() }
func (c *HTTP2TLSConnection) GetAgentUUID() string    { return c.AgentUUID }
