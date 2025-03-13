package h1tls

import (
	"firestarter/internal/interfaces"
	"firestarter/internal/listener"
	"firestarter/internal/router"
	"firestarter/internal/types"
	"fmt"
	"github.com/go-chi/chi/v5"
)

// Factory creates HTTP/1.1 TLS listeners
type Factory struct{}

func (f *Factory) CreateListener(id string, port string, connManager interfaces.ConnectionManager) (types.Listener, error) {
	r := chi.NewRouter()
	router.SetupRoutes(r)

	fmt.Printf("|CREATE| HTTP/1.1 TLS Listener %s configured on port %s\n", id, port)

	return listener.NewConcreteListener(id, port, interfaces.H1TLS, r, connManager), nil
}
