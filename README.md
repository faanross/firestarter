## DEVLOG OVERVIEW
### Set-up Basic Listener, Router, Handler
- Implement HTTP/1.1 clear listener
- Implement Chi router 
- Implement simple route 
- Implement simple handler 

### Implement Listener Factory
- Capable of producing multiple listeners

### Implement Stop()
- Add a `listeners` slice pointer to keep track of all listeners 
- Add a Stop() method to listener to give us ability to intentionally shut down listener
- `listeners` can be iterated through to shut down all listeners

### Graceful Shutdown
- Use context + sig term to ensure shutting down server gracefully exits all listeners
- Also integrate wait group into listener start goroutines to ensure proper synchronization during shutdown

### Implement Abstract Factory + H2C
- Create ability for more protocols in a generic manner using abstract factory
- Also implement H2C, test both H1C and H2C - they work

### Add Manager and Service
- Manager new struct to keep all active listeners
- Control is implemented in Service, which acts as master control
- For example creates a centralized method to config, add to manager, start a listener
- Or, stop a listener, remove from manager
- This is going to be very important once we get to frontend implementation to provide "snapshot method"

### Frontend
- Add basic web ui frontend, not-integrated
- Add websocket component in server + frontend, can connect

### ListenerTable Update
- Once listener is created, server sends single event update to UI, updating table
- However no persistence yet, ie if we refresh websocket connection table empties
- Create "snapshot" method to solve this
- Whenever UI connects to server, server immediately send list of all active listeners
- This achieves persistence

### Implement UI -> Server control
- Add `command` structure, allows frontend to send commands to server
- Add a `Stop` button in Listeners table, we can now stop individual listeners from UI