package h2c

import (
	"firestarter/internal/interfaces"
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

func (f *Factory) CreateListener(id string, port string, connManager interfaces.ConnectionManager) (types.Listener, error) {
	r := chi.NewRouter()
	router.SetupRoutes(r)

	h2s := &http2.Server{}

	h2cHandler := h2c.NewHandler(r, h2s)

	fmt.Printf("[👂🏻LSN] -> Listener (%s) created on port %s, protocol %s\n",
		id, port, interfaces.GetProtocolName(interfaces.H2C))

	concreteListener := listener.NewConcreteListener(id, port, interfaces.H2C, r, connManager)

	concreteListener.SetHandler(h2cHandler)

	return concreteListener, nil
}
