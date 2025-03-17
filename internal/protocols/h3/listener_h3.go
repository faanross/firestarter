package h3

import (
	"context"
	"crypto/tls"
	"firestarter/internal/connections"
	"firestarter/internal/interfaces"
	"fmt"
	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
	"net"
	"net/http"
	"time"
)

// HTTP3Listener implements the Listener interface for HTTP/3
type HTTP3Listener struct {
	ID          string
	Port        string
	Protocol    interfaces.ProtocolType
	Router      http.Handler
	CreatedAt   time.Time
	server      *EnhancedHTTP3Server
	udpListener net.PacketConn
	connManager interfaces.ConnectionManager
	quicConfig  *quic.Config
	ctx         context.Context
	cancel      context.CancelFunc
	tlsConfig   *tls.Config
}

// NewHTTP3Listener creates a new HTTP/3 listener
func NewHTTP3Listener(id string, port string, protocol interfaces.ProtocolType,
	router http.Handler, connManager interfaces.ConnectionManager) *HTTP3Listener {

	ctx, cancel := context.WithCancel(context.Background())

	return &HTTP3Listener{
		ID:          id,
		Port:        port,
		Protocol:    protocol,
		Router:      router,
		CreatedAt:   time.Now(),
		connManager: connManager,
		ctx:         ctx,
		cancel:      cancel,
	}
}

// Start implements the Listener interface
func (l *HTTP3Listener) Start() error {
	addr := fmt.Sprintf(":%s", l.Port)

	// Create a UDP listener
	udpConn, err := net.ListenPacket("udp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on UDP: %w", err)
	}
	l.udpListener = udpConn

	fmt.Printf("[H3-DEBUG] Created UDP listener for HTTP/3 on %s: %v\n", addr, udpConn != nil)

	// QUIC configuration - Now includes Versions
	l.quicConfig = &quic.Config{
		MaxIdleTimeout:  30 * time.Second,
		EnableDatagrams: true,
	}

	// Create the HTTP/3 server - Removed Versions field
	h3Server := &http3.Server{
		Handler:         l.Router,
		TLSConfig:       l.tlsConfig,
		EnableDatagrams: true,
	}

	// Set the QUIC config
	h3Server.QUICConfig = l.quicConfig

	fmt.Printf("[H3-DEBUG] Configured QUIC with IdleTimeout: %v, EnableDatagrams: %v\n",
		l.quicConfig.MaxIdleTimeout, l.quicConfig.EnableDatagrams)

	// Create the connection observer
	observer := connections.NewQuicConnectionObserver(l.connManager)

	// Create our enhanced server with connection tracking
	l.server = NewEnhancedHTTP3Server(h3Server, observer)

	fmt.Printf("[H3-DEBUG] Created enhanced HTTP/3 server with observer: %v\n", observer != nil)

	fmt.Printf("[H3-DEBUG] Created HTTP/3 server with TLS config: %v, Handler type: %T\n",
		l.tlsConfig != nil, l.Router)

	// MODIFICATION: Instead of using l.server.Serve, we'll create a QUIC listener
	// and manually handle connections to ensure they pass through our enhanced server
	quicListener, err := quic.Listen(udpConn, l.tlsConfig, l.quicConfig)
	if err != nil {
		return fmt.Errorf("failed to create QUIC listener: %w", err)
	}

	fmt.Printf("[H3-DEBUG] Created QUIC listener for HTTP/3\n")

	// Start accepting connections in a goroutine
	go func() {
		for {
			// Accept a new QUIC connection
			conn, err := quicListener.Accept(l.ctx)
			if err != nil {
				if l.ctx.Err() == nil {
					// Only log errors that aren't due to intentional shutdown
					fmt.Printf("HTTP/3 accept error: %v\n", err)
				}
				return
			}

			fmt.Printf("[H3-DEBUG] Accepted new QUIC connection from %s\n", conn.RemoteAddr())

			// Explicitly pass the connection to our enhanced server to ensure
			// our connection tracking hooks are called
			go func(c quic.Connection) {
				if err := l.server.ServeQUICConn(c); err != nil {
					fmt.Printf("HTTP/3 connection serving error: %v\n", err)
				}
			}(conn)
		}
	}()

	fmt.Printf("[H3-DEBUG] HTTP/3 listener %s started on %s with UDP listener: %v, TLS config present: %v\n",
		l.ID, addr, l.udpListener != nil, l.tlsConfig != nil)

	return nil
}

// Stop implements the Listener interface
func (l *HTTP3Listener) Stop() error {
	if l.server == nil {
		return fmt.Errorf("server not started")
	}

	fmt.Printf("|STOP| Shutting down HTTP/3 listener %s on port %s\n", l.ID, l.Port)

	// Cancel the connregistry to signal shutdown
	l.cancel()

	// Close the server
	err := l.server.Close()
	if err != nil {
		return fmt.Errorf("error shutting down HTTP/3 listener: %v", err)
	}

	// Close the UDP listener
	if l.udpListener != nil {
		l.udpListener.Close()
	}

	fmt.Printf("|STOP| HTTP/3 listener %s on port %s shut down successfully\n", l.ID, l.Port)
	return nil
}

// Interface method implementations
func (l *HTTP3Listener) GetID() string           { return l.ID }
func (l *HTTP3Listener) GetPort() string         { return l.Port }
func (l *HTTP3Listener) GetProtocol() string     { return "HTTP/3" }
func (l *HTTP3Listener) GetCreatedAt() time.Time { return l.CreatedAt }

func (l *HTTP3Listener) SetTLSConfig(config *tls.Config) {
	l.tlsConfig = config

	// If we already have a server, update its TLS config
	if l.server != nil && l.server.Server != nil {
		l.server.Server.TLSConfig = config
	}
}
