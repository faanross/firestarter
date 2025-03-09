package h2c

import (
	"firestarter/internal/listener"
	"firestarter/internal/router"
	"firestarter/internal/types"
	"fmt"
	"github.com/go-chi/chi/v5"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

// Factory creates HTTP/2 cleartext listeners
type Factory struct{}

func (f *Factory) CreateListener(id string, port string) (types.Listener, error) {

	// Create a router and set up routes
	r := chi.NewRouter()
	router.SetupRoutes(r)

	// Configure for HTTP/2
	h2s := &http2.Server{}

	// Wrap the router with h2c handler
	// This allows HTTP/2 connections over cleartext TCP
	h2cHandler := h2c.NewHandler(r, h2s)

	fmt.Printf("|CREATE| HTTP/2 Listener %s configured on port %s\n", id, port)

	// Create a concrete listener with the H2C protocol type
	concreteListener := listener.NewConcreteListener(id, port, types.H2C, r)

	// Set the H2C handler to be used when starting the server
	concreteListener.SetHandler(h2cHandler)

	return concreteListener, nil
}
