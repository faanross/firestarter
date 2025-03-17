package router

import (
	"fmt"
	"net/http"
)

// HTTP2HeaderLogger logs HTTP/2 specific header information
func HTTP2HeaderLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only act on HTTP/2 requests
		if r.ProtoMajor == 2 {
			// Log HTTP/2 specific information
			fmt.Printf("[HTTP/2] Request to %s has %d headers\n",
				r.URL.Path, len(r.Header))

			// Check if we have any HTTP/2 specific issues with our UUID header
			if r.Header.Get("X-Agent-UUID") == "" {
				// Check for lowercase variant (HTTP/2 can normalize headers)
				if r.Header.Get("x-agent-uuid") != "" {
					fmt.Printf("[HTTP/2] Found lowercase UUID header, HTTP/2 may have normalized it\n")
					// Copy the lowercase value to our standard casing
					r.Header.Set("X-Agent-UUID", r.Header.Get("x-agent-uuid"))
				}
			}
		}

		next.ServeHTTP(w, r)
	})
}
