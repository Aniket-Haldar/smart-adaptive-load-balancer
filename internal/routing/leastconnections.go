package routing

import (
	"fmt"
	"sync"

	"github.com/sbirmecha99/smart-adaptive-load-balancer/internal/core"
)

// LeastConnectionsRouter selects backends with the fewest active connections
// mutexes to prevent race around conditions
type LeastConnectionsRouter struct {
	mu sync.Mutex
}

func NewLeastConnectionsRouter() *LeastConnectionsRouter {
	return &LeastConnectionsRouter{}
}

// implements Router interface
func (lc *LeastConnectionsRouter) GetNextAvailableServer(
	backends []*core.Backend,
) *core.Backend {

	lc.mu.Lock()
	defer lc.mu.Unlock()

	n := len(backends)
	if n == 0 {
		fmt.Println("No Servers Present")
		return nil
	}

	var selected *core.Backend
	minConns := int64(^uint64(0) >> 1) // intitialises this to the maximum number

	for _, backend := range backends {
		backend.Mutex.Lock()
		alive := backend.Alive
		active := backend.ActiveConns
		backend.Mutex.Unlock()

		if alive && active < minConns {
			minConns = active
			selected = backend
		}
	}

	if selected != nil {
		// Increment connection count as we assign this server
		selected.Mutex.Lock()
		selected.ActiveConns++
		selected.Mutex.Unlock()
	}

	return selected
}
func (lc *LeastConnectionsRouter) Name() string { return "LeastConnections" }