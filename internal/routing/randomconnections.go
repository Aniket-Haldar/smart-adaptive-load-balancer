package routing

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/sbirmecha99/smart-adaptive-load-balancer/internal/core"
)

// RandomRouter selects backends randomly
type RandomRouter struct {
	mu  sync.Mutex
	rng *rand.Rand
}

func NewRandomRouter() *RandomRouter {
	return &RandomRouter{
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// implements Router interface
func (rr *RandomRouter) GetNextAvailableServer(
	backends []*core.Backend,
) *core.Backend {

	rr.mu.Lock()
	defer rr.mu.Unlock()

	n := len(backends)
	if n == 0 {
		fmt.Println("No Servers Present")
		return nil
	}

	//we are trying n times to check for alive backedn
	for i := 0; i < n; i++ {
		idx := rr.rng.Intn(n)
		backend := backends[idx]

		backend.Mutex.Lock()
		alive := backend.Alive
		backend.Mutex.Unlock()

		if alive {
			return backend
		}
	}

	return nil
}
func (rn *RandomRouter) Name() string { return "Random" }