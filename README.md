## devlog overview
### set-up basic listener, router, handler
- implement HTTP/1.1 clear listener (`/cmd/server/main.go`)
- implement Chi router (`/cmd/server/main.go`)
- implement simple route (`/internal/router/routes.go`)
- implement simple handler (`/internal/router/handlers.go`)

### implement listener factory
- capable of producing multiple listeners
- set this up in `/internal/factory/factory.go`

### implement graceful shutdown + Stop() 
- Use context + sig term to ensure all listeners are shut down gracefully when server closes
- Add a Stop() method to listener to give us ability to intentionally shut down listener