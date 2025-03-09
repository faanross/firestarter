## devlog overview
### set-up basic listener, router, handler
- implement HTTP/1.1 clear listener
- implement Chi router 
- implement simple route 
- implement simple handler 

### implement listener factory
- capable of producing multiple listeners

### Implement Stop()
- add a `listeners` slice pointer to keep track of all listeners 
- add a Stop() method to listener to give us ability to intentionally shut down listener
- `listeners` can be iterated through to shut down all listeners

### Graceful shutdown
- use context + sig term to ensure shutting down server gracefully exits all listeners
- also integrate wait group into listener start goroutines to ensure proper synchronization during shutdown

### Implement Abstract Factory + H2C
- Create ability for more protocols in a generic manner using abstract factory
- Also implement H2C, test both H1C and H2C - they work

### Add Manager and Service
- Manager new struct to keep all active listeners
- Control is implemented in Service, which acts as master control
- For example creates a centralized method to config, add to manager, start a listener
- Or, stop a listener, remove from manager
- This is going to be very important once we get to frontend implementation to provide "snapshot method"