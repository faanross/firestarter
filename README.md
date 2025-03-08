## devlog overview
### set-up basic listener, router, handler
- implement HTTP/1.1 clear listener (`/cmd/server/main.go`)
- implement Chi router (`/cmd/server/main.go`)
- implement simple route (`/internal/router/routes.go`)
- implement simple handler (`/internal/router/handlers.go`)

### implement listener factory
- capable of producing multiple listeners
- set this up in `/internal/factory/factory.go`