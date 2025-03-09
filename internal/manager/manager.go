package manager

import (
	"firestarter/internal/types"
	"fmt"
	"sync"
)

// ListenerManager keeps track of all active listeners
type ListenerManager struct {
	listeners map[string]types.Listener
	mu        sync.RWMutex // For thread safety
}

// NewListenerManager creates a new listener manager
func NewListenerManager() *ListenerManager {
	return &ListenerManager{
		listeners: make(map[string]types.Listener),
	}
}

// AddListener adds a listener to the manager
func (m *ListenerManager) AddListener(listener types.Listener) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	id := listener.GetID()
	if _, exists := m.listeners[id]; exists {
		return fmt.Errorf("listener with ID %s already exists", id)
	}

	m.listeners[id] = listener
	return nil
}

// GetListener retrieves a listener by ID
func (m *ListenerManager) GetListener(id string) (types.Listener, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	listener, exists := m.listeners[id]
	if !exists {
		return nil, fmt.Errorf("no listener found with ID %s", id)
	}

	return listener, nil
}

// ListListeners returns a slice of all managed listeners
func (m *ListenerManager) ListListeners() []types.Listener {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]types.Listener, 0, len(m.listeners))
	for _, listener := range m.listeners {
		result = append(result, listener)
	}

	return result
}

// RemoveListener removes a listener from the manager
func (m *ListenerManager) RemoveListener(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.listeners[id]; !exists {
		return fmt.Errorf("no listener found with ID %s", id)
	}

	delete(m.listeners, id)
	return nil
}

// Count returns the number of listeners currently managed
func (m *ListenerManager) Count() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.listeners)
}
