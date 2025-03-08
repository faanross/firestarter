package main

import (
	"firestarter/internal/router"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

const serverAddr = ":7777"

func main() {

	r := chi.NewRouter()

	router.SetupRoutes(r)

	server := &http.Server{
		Addr:    serverAddr,
		Handler: r,
	}

	log.Printf("Starting HTTP server on %v", serverAddr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
