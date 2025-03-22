package router

import (
	"github.com/go-chi/chi/v5"
)

// SetupRoutes configures all routes for the application
func SetupRoutes(r chi.Router) {

	// Apply middleware to all routes
	r.Use(AgentUUIDHeaderMiddleware)

	// Define our root endpoint
	r.Get("/", RootHandler)

	// Add test endpoints for connection tracking verification
	r.Get("/quick", QuickResponseHandler)
	r.Get("/slow", SlowResponseHandler)

	r.Post("/ping", PingHandler)
}
