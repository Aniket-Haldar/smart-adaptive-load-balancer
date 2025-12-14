package routing

import "github.com/sbirmecha99/smart-adaptive-load-balancer/internal/core"

// Router defines how a backend server is selected
type Router interface {
	GetNextAvailableServer([]*core.Backend) *core.Backend
}


