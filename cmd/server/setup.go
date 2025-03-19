package main

import (
	"firestarter/internal/connections"
	"firestarter/internal/factory"
	"firestarter/internal/manager"
	"firestarter/internal/router"
	"firestarter/internal/service"
	"firestarter/internal/websocket"
	"fmt"
)

func ApplicationSetup() (*manager.ListenerManager, *service.ListenerService) {
	// Start our Websocket (:8080) for UI integration
	websocket.StartWebSocketServer()

	// Initialize connection registry for UUID tracking
	router.InitializeConnectionRegistry()
	connections.SetConnectionRegistry(router.GetConnectionRegistry())

	// Create the components
	connectionManager := connections.NewConnectionManager()
	// Connect the registry to the connection manager
	router.ConnectRegistryToManager(connectionManager)

	// Link the Connection Manager to the WebSocket server
	wsServer := websocket.GetGlobalWSServer()
	if wsServer != nil {
		connectionManager.SetWebSocketServer(wsServer)
		fmt.Println("[INIT] WebSocket server linked to Connection Manager")
	} else {
		fmt.Println("[INIT-ERROR] WebSocket server not available for Connection Manager!")
	}

	af := factory.NewAbstractFactory(connectionManager)
	lm := manager.NewListenerManager()
	ls := service.NewListenerService(af, lm, connectionManager)

	return lm, ls

}
