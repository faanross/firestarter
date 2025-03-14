package h3

import (
	"firestarter/internal/connections"
	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
)

// EnhancedHTTP3Server extends the standard HTTP/3 server with connection tracking
type EnhancedHTTP3Server struct {
	*http3.Server
	observer *connections.QuicConnectionObserver
}

// NewEnhancedHTTP3Server creates a new HTTP/3 server with connection tracking
func NewEnhancedHTTP3Server(server *http3.Server, observer *connections.QuicConnectionObserver) *EnhancedHTTP3Server {
	return &EnhancedHTTP3Server{
		Server:   server,
		observer: observer,
	}
}

// ServeQUICConn intercepts QUIC connections for tracking before handling
func (s *EnhancedHTTP3Server) ServeQUICConn(conn quic.Connection) {
	// Notify our observer about the new connection
	s.observer.OnConnectionEstablished(conn)

	// Continue with normal HTTP/3 handling
	s.Server.ServeQUICConn(conn)
}
