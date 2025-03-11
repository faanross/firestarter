package connections

import (
	"firestarter/internal/interfaces"
	"net"
	"time"
)

// HTTP1Connection represents an HTTP/1 Clear specific connection
type HTTP1Connection struct {
	BaseConnection
	Conn net.Conn
}

func NewHTTP1Connection(conn net.Conn) *HTTP1Connection {
	return &HTTP1Connection{
		BaseConnection: BaseConnection{
			ID:        GenerateUniqueID(),
			Protocol:  interfaces.H1C,
			CreatedAt: time.Now().UTC(),
		},
		Conn: conn,
	}
}

func (c *HTTP1Connection) GetID() string { return c.ID }
func (c *HTTP1Connection) GetProtocol() interfaces.ProtocolType {
	return c.Protocol
}
func (c *HTTP1Connection) GetCreatedAt() time.Time { return c.CreatedAt }
func (c *HTTP1Connection) GetPort() string         { return c.Port }
func (c *HTTP1Connection) Close() error            { return c.Conn.Close() }
