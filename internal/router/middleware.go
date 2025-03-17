package router

import (
	"context"
	"net/http"
)

// Key type for context values
type contextKey string

// Constants for context keys
const (
	AgentUUIDKey contextKey = "agent-uuid"
)

// AgentUUIDMiddleware extracts the agent UUID from request headers
func AgentUUIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract UUID from header
		agentUUID := r.Header.Get("X-Agent-UUID")

		// Store in request context, even if empty
		ctx := context.WithValue(r.Context(), AgentUUIDKey, agentUUID)

		// Call the next handler with the updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetAgentUUIDFromRequest extracts the agent UUID from a request context
func GetAgentUUIDFromRequest(r *http.Request) string {
	if uuid, ok := r.Context().Value(AgentUUIDKey).(string); ok {
		return uuid
	}
	return "" // Return empty string if not found
}
