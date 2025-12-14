package routing

import "github.com/sbirmecha99/smart-adaptive-load-balancer/internal/core"

type Router interface {
	GetNextAvailableServer(backends []*core.Backend) *core.Backend
} 