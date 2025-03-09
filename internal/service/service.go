// internal/service/service.go
package service

import (
	"firestarter/internal/factory"
	"firestarter/internal/manager"
	"firestarter/internal/types"
	"fmt"
	"sync"
)

// ListenerService coordinates listener lifecycle operations
type ListenerService struct {
	factory *factory.AbstractFactory
	manager *manager.ListenerManager
}

// NewListenerService creates a new listener service
func NewListenerService(factory *factory.AbstractFactory, manager *manager.ListenerManager) *ListenerService {
	return &ListenerService{
		factory: factory,
		manager: manager,
	}
}

// CreateAndStartListener creates a listener, registers it with the manager, and starts it
func (s *ListenerService) CreateAndStartListener(protocol types.ProtocolType, port string, wg *sync.WaitGroup) (types.Listener, error) {
	// Create the listener
	listener, err := s.factory.CreateListener(protocol, port)
	if err != nil {
		return nil, fmt.Errorf("failed to create listener: %w", err)
	}

	// Register with the manager
	err = s.manager.AddListener(listener)
	if err != nil {
		return nil, fmt.Errorf("failed to register listener: %w", err)
	}

	// Start the listener in a goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := listener.Start()
		if err != nil {
			// Check if this is just a server closed error (expected during shutdown)
			if err.Error() != "http: Server closed" {
				fmt.Printf("Error starting listener %s: %v\n", listener.GetID(), err)
				// Remove from manager if it failed to start unexpectedly
				_ = s.manager.RemoveListener(listener.GetID())
			}
		}
	}()

	return listener, nil
}

// StopListener stops a listener and removes it from the manager
func (s *ListenerService) StopListener(id string) error {
	// Get the listener from the manager
	listener, err := s.manager.GetListener(id)
	if err != nil {
		return err
	}

	// Stop the listener
	err = listener.Stop()
	if err != nil {
		return fmt.Errorf("failed to stop listener: %w", err)
	}

	// Remove from manager
	err = s.manager.RemoveListener(id)
	if err != nil {
		return fmt.Errorf("failed to remove listener from manager: %w", err)
	}

	return nil
}

// StopAllListeners stops all managed listeners
func (s *ListenerService) StopAllListeners(wg *sync.WaitGroup) {
	fmt.Println("Shutting down all listeners...")

	// Get all listeners
	listeners := s.manager.ListListeners()

	// Stop each listener
	for _, listener := range listeners {
		id := listener.GetID()
		fmt.Printf("Stopping listener %s\n", id)
		err := listener.Stop()
		if err != nil {
			fmt.Printf("Error stopping listener %s: %v\n", id, err)
		}

		// Remove from manager
		_ = s.manager.RemoveListener(id)
	}

	fmt.Println("All listeners shutdown initiated. Waiting for goroutines to complete...")
	wg.Wait()
	fmt.Println("All server goroutines terminated. Exiting...")
}

// GetAllListeners returns all managed listeners
func (s *ListenerService) GetAllListeners() []types.Listener {
	return s.manager.ListListeners()
}

// GetManager returns the manager instance
func (s *ListenerService) GetManager() *manager.ListenerManager {
	return s.manager
}
