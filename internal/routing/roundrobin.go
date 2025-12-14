package routing

import "github.com/sbirmecha99/smart-adaptive-load-balancer/internal/core"

// RoundRobinRouter selects backends in circular order
type RoundRobinRouter struct {
	current int
}

func NewRoundRobinRouter() *RoundRobinRouter {
	return &RoundRobinRouter{
		current: 0,
	}
}

// implements Router interface
func (rr *RoundRobinRouter) GetNextAvailableServer(
	backends []*core.Backend,
) *core.Backend {

	n := len(backends)
	if n == 0 {
		return nil
	}

	//checking each server
	for i := 0; i < n; i++ {
		idx := (rr.current + i) % n
		backend := backends[idx]

		if backend.Alive {
			rr.current = idx + 1
			return backend
		}
	}

	//if no alive , then nil returned
	return nil
}
