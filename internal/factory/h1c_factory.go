package factory

import (
	"firestarter/internal/router"
	"fmt"
	"github.com/go-chi/chi/v5"
	"math/rand"
)

// H1CFactory creates HTTP/1.1 cleartext listeners
type H1CFactory struct{}

func (f *H1CFactory) CreateListener(id string, port string) (Listener, error) {
	// If ID is empty, generate a random one
	if id == "" {
		id = fmt.Sprintf("listener_%06d", rand.Intn(1000000))
	}

	r := chi.NewRouter()
	router.SetupRoutes(r)

	fmt.Printf("|CREATE| H1C Listener %s configured on port %s\n", id, port)

	return &ConcreteListener{
		ID:       id,
		Port:     port,
		Protocol: H1C,
		Router:   r,
	}, nil
}
