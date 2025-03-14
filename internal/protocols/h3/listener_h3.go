package h3

import (
	"context"
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

	// Create a QUIC transport
	transport := &quic.Transport{
		Conn: udpConn,
	}

	// QUIC configuration
	l.quicConfig = &quic.Config{
		MaxIdleTimeout:  30 * time.Second,
		EnableDatagrams: true,
	}

	// Create the HTTP/3 server
	h3Server := &http3.Server{
		Handler:    l.Router,
		QuicConfig: l.quicConfig,
	}

	// Create the connection observer
	observer := connections.NewQuicConnectionObserver(l.connManager)

	// Create our enhanced server with connection tracking
	l.server = NewEnhancedHTTP3Server(h3Server, observer)

	fmt.Printf("|START| HTTP/3 Listener %s serving on %s\n", l.ID, addr)

	// Start the server (non-blocking)
	go func() {
		err := l.server.Serve(transport)
		if err != nil && l.ctx.Err() == nil {
			// Only log errors that aren't due to intentional shutdown
			fmt.Printf("HTTP/3 server error: %v\n", err)
		}
	}()

	return nil
}

// Stop implements the Listener interface
func (l *HTTP3Listener) Stop() error {
	if l.server == nil {
		return fmt.Errorf("server not started")
	}

	fmt.Printf("|STOP| Shutting down HTTP/3 listener %s on port %s\n", l.ID, l.Port)

	// Cancel the context to signal shutdown
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
