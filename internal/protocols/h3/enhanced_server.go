package h3

import (
	"firestarter/internal/connections"
	"fmt"
	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
	"net"
	"net/http"
	"sync"
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
func (s *EnhancedHTTP3Server) ServeQUICConn(conn quic.Connection) error {
	fmt.Printf("[H3-DEBUG] ServeQUICConn called for connection from: %s\n",
		conn.RemoteAddr().String())

	// Store connection in a map with empty UUID initially
	h3ConnectionUUIDs.Store(conn.RemoteAddr().String(), "")

	// Install header extractor
	s.Server.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract UUID from headers in HTTP/3 requests
		agentUUID := r.Header.Get("X-Agent-UUID")
		if agentUUID != "" {
			// We found a UUID, associate it with this QUIC connection
			fmt.Printf("[HTTP/3] Extracted agent UUID: %s from QUIC connection\n", agentUUID)

			// Update the UUID map
			h3ConnectionUUIDs.Store(conn.RemoteAddr().String(), agentUUID)

			// Update any existing tracked connections
			// This is more complex for HTTP/3 and would need custom implementation
		}

		// Call the original handler
		s.Server.Handler.ServeHTTP(w, r)
	})

	// Get port from listening address
	port := "unknown"
	if s.Server.Addr != "" {
		_, portStr, _ := net.SplitHostPort(s.Server.Addr)
		if portStr != "" {
			port = portStr
		}
	}

	s.observer.OnConnectionEstablished(conn, port)

	fmt.Printf("[H3-SERVER-DEBUG] Observer notified, continuing with standard HTTP/3 handling\n")

	// Continue with normal HTTP/3 handling and return its error
	return s.Server.ServeQUICConn(conn)
}

// Add a global map to track HTTP/3 connection UUIDs
var h3ConnectionUUIDs sync.Map
