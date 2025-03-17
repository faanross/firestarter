package router

import (
	"github.com/go-chi/chi/v5"
)

// In router/routes_kavqmwsu.go

// SetupRoutes configures all routes for the application
func SetupRoutes(r chi.Router) {

	// Apply middleware to all routes
	r.Use(AgentUUIDMiddleware)
	
	// Define our root endpoint
	r.Get("/", RootHandler)

	// Add test endpoints for connection tracking verification
	r.Get("/quick", QuickResponseHandler)
	r.Get("/slow", SlowResponseHandler)
}
