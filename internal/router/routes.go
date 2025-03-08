package router

import "github.com/go-chi/chi/v5"

func SetupRoutes(r chi.Router) {
	// Define our root endpoint
	r.Get("/", RootHandler)
}
