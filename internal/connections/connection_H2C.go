package connections

import (
	"firestarter/internal/interfaces"
	"net"
	"time"
)

// HTTP2Connection represents an HTTP/2 specific connection
type HTTP2Connection struct {
	BaseConnection
	Conn net.Conn
}

func NewHTTP2Connection(conn net.Conn) *HTTP2Connection {
	return &HTTP2Connection{
		BaseConnection: BaseConnection{
			ID:        GenerateUniqueID(),
			Protocol:  interfaces.H2C,
			CreatedAt: time.Now().UTC(),
		},
		Conn: conn,
	}
}

func (c *HTTP2Connection) GetID() string { return c.ID }
func (c *HTTP2Connection) GetProtocol() interfaces.ProtocolType {
	return c.Protocol
}
func (c *HTTP2Connection) GetCreatedAt() time.Time { return c.CreatedAt }
func (c *HTTP2Connection) GetPort() string         { return c.Port }
func (c *HTTP2Connection) Close() error            { return c.Conn.Close() }
