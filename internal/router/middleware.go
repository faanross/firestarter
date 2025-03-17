package router

import (
	"context"
	"firestarter/internal/connregistry"
	"firestarter/internal/interfaces"
	"fmt"
	"net/http"
	"sync"
)

// Global connection registry instance
var connectionRegistry *connregistry.ConnectionRegistry

var (
	// Track which UUIDs we've already logged for each connection
	processedUUIDs    = make(map[string]bool)
	processedUUIDsMux sync.RWMutex
)

// InitializeConnectionRegistry creates and sets up the global connection registry
func InitializeConnectionRegistry() {
	if connectionRegistry == nil {
		fmt.Println("Initializing global connection registry")
		connectionRegistry = connregistry.NewConnectionRegistry()
	}
}

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
		if agentUUID != "" && connectionRegistry != nil {
			connectionRegistry.RegisterUUID(r, agentUUID)
		}

		// Call the next handler with the updated connregistry
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetAgentUUIDFromRequest extracts the agent UUID from a request connregistry
func GetAgentUUIDFromRequest(r *http.Request) string {
	if uuid, ok := r.Context().Value(AgentUUIDKey).(string); ok {
		return uuid
	}
	return "" // Return empty string if not found
}

// ConnectRegistryToManager connects the registry to a connection manager
func ConnectRegistryToManager(manager interfaces.ConnectionManager) {
	if connectionRegistry != nil {
		connectionRegistry.SetConnectionManager(manager)
	} else {
		fmt.Println("Warning: Cannot connect registry to manager - registry not initialized")
	}
}

// GetConnectionRegistry returns the global connection registry
func GetConnectionRegistry() *connregistry.ConnectionRegistry {
	return connectionRegistry
}
