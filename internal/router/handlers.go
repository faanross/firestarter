package router

import (
	"fmt"
	"net/http"
	"time"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	// Log a message on the server side
	fmt.Println("You hit the endpoint:", r.URL.Path)

	// Send a response to the client
	w.Write([]byte("I'm Mister Derp!"))
}

// QuickResponseHandler returns immediately
func QuickResponseHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Quick request received on:", r.URL.Path)
	w.Write([]byte("Quick response completed"))
}

// SlowResponseHandler simulates a slow API call
func SlowResponseHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Slow request started on:", r.URL.Path)

	// Simulate processing time
	time.Sleep(10 * time.Second)

	fmt.Println("Slow request completed on:", r.URL.Path)
	w.Write([]byte("Slow response completed after 10 seconds"))
}

// PingHandler responds to ping requests from agents
func PingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Ping request received from:", r.RemoteAddr)

	// Extract and log the agent UUID for debugging
	agentUUID := r.Header.Get("X-Agent-UUID")
	if agentUUID != "" {
		fmt.Printf("Ping from agent with UUID: %s\n", agentUUID)
	}

	w.Write([]byte("pong"))
}
