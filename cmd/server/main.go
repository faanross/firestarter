package main

import (
	"log"
	"net/http"
)

const serverAddr = ":7777"

func main() {
	server := &http.Server{
		Addr: serverAddr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Default handler that does nothing
		}),
	}

	log.Printf("Starting HTTP server on %v", serverAddr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
