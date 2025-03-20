package router

import (
	"context"
	"firestarter/internal/connregistry"
	"firestarter/internal/interfaces"
	"fmt"
	"net/http"
	"sync"
)

var (
	// Track which UUIDs we've already logged for each connection
	processedUUIDs    = make(map[string]bool)
	processedUUIDsMux sync.RWMutex
)

// Key type for connregistry values
type contextKey string

// Constants for connregistry keys
const (
	AgentUUIDKey contextKey = "agent-uuid"
)

// AgentUUIDMiddleware extracts the agent UUID from request headers
func AgentUUIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract UUID from header
		agentUUID := r.Header.Get("X-Agent-UUID")

		// Log the extraction
		if agentUUID != "" {
			// Generate a unique key for this connection+UUID combination
			connUUIDKey := r.RemoteAddr + ":" + agentUUID

			processedUUIDsMux.RLock()
			alreadyProcessed := processedUUIDs[connUUIDKey]
			processedUUIDsMux.RUnlock()

			if !alreadyProcessed {
				processedUUIDsMux.Lock()
				processedUUIDs[connUUIDKey] = true
				processedUUIDsMux.Unlock()

				fmt.Printf("[UUID-Track-DEBUG] Middleware: Extracted agent UUID: %s from request to %s (Remote: %s)\n",
					agentUUID, r.URL.Path, r.RemoteAddr)
			}
		}

		// Store in request connregistry
		ctx := context.WithValue(r.Context(), AgentUUIDKey, agentUUID)

		// Look up the connection from our registry
		if agentUUID != "" && connregistry.GlobalConnectionRegistry != nil {
			connregistry.GlobalConnectionRegistry.RegisterUUID(r, agentUUID)
		}

		// Call the next handler with the updated connregistry
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// ConnectRegistryToManager connects the registry to a connection manager
func ConnectRegistryToManager(manager interfaces.ConnectionManager) {
	if connregistry.GlobalConnectionRegistry != nil {
		connregistry.GlobalConnectionRegistry.SetConnectionManager(manager)
	} else {
		fmt.Println("Warning: Cannot connect registry to manager - registry not initialized")
	}
}
