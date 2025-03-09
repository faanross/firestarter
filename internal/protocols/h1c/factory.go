package h1c

import (
	"firestarter/internal/listener"
	"firestarter/internal/router"
	"firestarter/internal/types"
	"fmt"
	"github.com/go-chi/chi/v5"
	"math/rand"
)

// Factory creates HTTP/1.1 cleartext listeners
type Factory struct{}

func (f *Factory) CreateListener(id string, port string) (types.Listener, error) {
	// If ID is empty, generate a random one
	if id == "" {
		id = fmt.Sprintf("listener_%06d", rand.Intn(1000000))
	}

	r := chi.NewRouter()
	router.SetupRoutes(r)

	fmt.Printf("|CREATE| H1C Listener %s configured on port %s\n", id, port)

	return listener.NewConcreteListener(id, port, types.H1C, r), nil
}
