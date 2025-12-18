package health

import (
	"log"
	"net"
	"time"
	"github.com/sbirmecha99/smart-adaptive-load-balancer/internal/core"
)

type Checker struct {
	Pool     *core.ServerPool
	Interval time.Duration
	Timeout  time.Duration
}

// TCP-level health check (works for L4 and L7)
func (c *Checker) Start() {
	ticker := time.NewTicker(c.Interval)

	go func() {
		for range ticker.C {
			backends := c.Pool.GetServers() // SAFE snapshot
			for _, backend := range backends {
				go c.checkBackend(backend)
			}
		}
	}()
}

func (c *Checker) checkBackend(b *core.Backend) {
	start := time.Now()

	conn, err := net.DialTimeout("tcp", b.Address, c.Timeout)

	b.Mutex.Lock()
	defer b.Mutex.Unlock()

	if err != nil {
		b.Alive = false
		log.Printf("[HEALTH] backend DOWN: %s", b.Address)
		return
	}

	_ = conn.Close()
	b.Alive = true
	b.Latency = time.Since(start)
}
