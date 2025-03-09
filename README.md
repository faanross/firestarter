## devlog overview
### set-up basic listener, router, handler
- implement HTTP/1.1 clear listener (`/cmd/server/main.go`)
- implement Chi router (`/cmd/server/main.go`)
- implement simple route (`/internal/router/routes.go`)
- implement simple handler (`/internal/router/handlers.go`)

### implement listener factory
- capable of producing multiple listeners
- set this up in `/internal/factory/factory.go`

### Implement Stop()
- add a `listeners` slice pointer to keep track of all listeners 
- add a Stop() method to listener to give us ability to intentionally shut down listener
- `listeners` can be iterated through to shut down all listeners

### Graceful shutdown
- use context + sig term to ensure shutting down server gracefully exits all listeners
- also integrate wait group into listener start goroutines to ensure proper synchronization during shutdown