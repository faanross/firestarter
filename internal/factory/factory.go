package factory

import (
	"context"
	"firestarter/internal/router"
	"fmt"
	"github.com/go-chi/chi/v5"
	"math/rand"
	"net/http"
	"time"
)

type ListenerFactory interface {
	CreateListener(id string, port string) (Listener, error)
}

// Listener represents an HTTP server instance
type Listener struct {
	ID     string
	Port   string
	Router *chi.Mux
	server *http.Server
}

// CreateListener generates a new listener with a random port and unique ID
func (f *ListenerFactory) CreateListener(port string) (*Listener, error) {
	// Generate a random ID (6 digits)
	id := fmt.Sprintf("listener_%06d", rand.Intn(1000000))

	r := chi.NewRouter()

	router.SetupRoutes(r)

	fmt.Printf("|CREATE| Listener %s configured on on port %s\n", id, port)

	return &Listener{
		ID:     id,
		Port:   port,
		Router: r,
	}, nil
}

func (l *Listener) Start() error {
	addr := fmt.Sprintf(":%s", l.Port)
	fmt.Printf("|START| Listener %s serving on %s\n", l.ID, addr)

	// Create the server instance
	l.server = &http.Server{
		Addr:    addr,
		Handler: l.Router,
	}

	return l.server.ListenAndServe()
}

// Stop gracefully shuts down the HTTP server with a timeout
func (l *Listener) Stop() error {
	if l.server == nil {
		return fmt.Errorf("server not started")
	}

	// Create a context with a timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	fmt.Printf("|STOP| Shutting down listener %s on port %s\n", l.ID, l.Port)

	// Shutdown the server gracefully
	err := l.server.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("error shutting down listener %s: %v", l.ID, err)
	}

	fmt.Printf("|STOP| Listener %s on port %s shut down successfully\n", l.ID, l.Port)
	return nil
}

type ProtocolType int

const (
	H1C ProtocolType = iota + 1
	//H1TLS
	H2C
	//H2TLS
	//H3
)

type Listeners interface {
	Start() error
	Stop() error
	GetProtocol() string
	GetPort() int
}
