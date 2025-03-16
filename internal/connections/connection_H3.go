package connections

import (
	"firestarter/internal/interfaces"
	"github.com/quic-go/quic-go"
	"time"
)

// HTTP3Connection represents an HTTP/3 specific connection over QUIC
type HTTP3Connection struct {
	BaseConnection
	QUICConn quic.Connection
}

// NewHTTP3Connection creates a new HTTP/3 connection
func NewHTTP3Connection(conn quic.Connection) *HTTP3Connection {
	return &HTTP3Connection{
		BaseConnection: BaseConnection{
			ID:        GenerateUniqueID(),
			Protocol:  interfaces.H3,
			CreatedAt: time.Now().UTC(),
		},
		QUICConn: conn,
	}
}

// Connection interface implementation
func (c *HTTP3Connection) GetID() string                        { return c.ID }
func (c *HTTP3Connection) GetProtocol() interfaces.ProtocolType { return c.Protocol }
func (c *HTTP3Connection) GetCreatedAt() time.Time              { return c.CreatedAt }
func (c *HTTP3Connection) GetPort() string                      { return c.Port }
func (c *HTTP3Connection) Close() error                         { return c.QUICConn.CloseWithError(0, "closed by server") }
func (c *HTTP3Connection) GetAgentUUID() string                 { return c.AgentUUID }
