package routing

import (
	"log"
	"fmt"
	"time"
	"sync"
	"github.com/sbirmecha99/smart-adaptive-load-balancer/internal/core"
)

type AdaptiveRouter struct {
	pool *core.ServerPool

	rr *RoundRobinRouter
	lc *LeastConnectionsRouter
	rn *RandomRouter

	currentAlgo string
	reason       string
	lastPicked   string
}
type Decision struct {
    Time     time.Time `json:"time"`

    Algo     string `json:"algo"`
    Reason   string `json:"reason"`
    Backend  string `json:"backend"`
}

var (
	DecisionLog []Decision
	DecisionMu  sync.Mutex
)



func NewAdaptiveRouter(pool *core.ServerPool) *AdaptiveRouter {
	return &AdaptiveRouter{
		pool:        pool,
		rr:          NewRoundRobinRouter(),
		lc:          NewLeastConnectionsRouter(),
		rn:          NewRandomRouter(),
		currentAlgo: "roundrobin",
		reason:      "normal_conditions",
	}
}

// ✅ SINGLE ENTRY POINT
func (ar *AdaptiveRouter) Pick() *core.Backend {
	backends := ar.pool.GetServers()
	if len(backends) == 0 {
		log.Println("[ADAPTIVE] no backends in pool")
		return nil
	}

	var totalConns int64
	var totalLatency int64
	var totalErrors int64
	var maxConns int64
	aliveCount := 0

	for _, b := range backends {
		b.Mutex.Lock()
		if b.Alive {
			aliveCount++
			totalConns += b.ActiveConns
			totalLatency += int64(b.Latency)
			totalErrors += b.ErrorCount
			if b.ActiveConns > maxConns {
				maxConns = b.ActiveConns
			}
		}
		b.Mutex.Unlock()
	}

	if aliveCount == 0 {
		log.Println("[ADAPTIVE] no alive backends")
		return nil
	}

	avgConns := totalConns / int64(aliveCount)
	avgLatency := totalLatency / int64(aliveCount)
	avgLatencyMs := avgLatency / int64(time.Millisecond)
	errorRate := float64(totalErrors) / float64(totalConns+1)

	log.Printf("[ADAPTIVE] algo=%s reason=%s picked=%s avgConns=%d maxConns=%d avgLatencyMs=%d errorRate=%.2f",
	ar.currentAlgo,
	ar.reason,
	ar.lastPicked,
	avgConns,
	maxConns,
	avgLatency/int64(time.Millisecond),
	errorRate,
	)

	var selected *core.Backend

	// 1️⃣ Error dominance
	if errorRate > 0.3 {
		ar.currentAlgo = "random"
		ar.reason = fmt.Sprintf("high_error_rate (%.2f)", errorRate)
		selected = ar.rn.GetNextAvailableServer(backends)

		// 2️⃣ Load skew
	} else if maxConns > 3 {
		ar.currentAlgo = "leastconnections"
		ar.reason = "high_concurrency"
		selected = ar.lc.GetNextAvailableServer(backends)

		// 3️⃣ Latency degradation
	} else if avgLatencyMs > 200 {
		ar.currentAlgo = "leastconnections"
		ar.reason = "high_latency"
		selected = ar.lc.GetNextAvailableServer(backends)

		// 4️⃣ Default
	} else {
		ar.currentAlgo = "roundrobin"
		ar.reason = "normal_conditions"
		selected = ar.rr.GetNextAvailableServer(backends)
	}

	if selected != nil {
	ar.lastPicked = selected.Address

	DecisionMu.Lock()
	DecisionLog = append(DecisionLog, Decision{
		Time:    time.Now(),
		Algo:    ar.currentAlgo,
		Reason:  ar.reason,
		Backend: selected.Address,
	})
	DecisionMu.Unlock()
}

	return selected
}
// adaptive.go
func (ar *AdaptiveRouter) GetNextAvailableServer(_ []*core.Backend) *core.Backend {
    return ar.Pick()
}
func (ar *AdaptiveRouter) Name() string        { return "adaptive" }
func (ar *AdaptiveRouter) CurrentAlgo() string { return ar.currentAlgo }
func (ar *AdaptiveRouter) Reason() string      { return ar.reason }
func (ar *AdaptiveRouter) LastPicked() string  { return ar.lastPicked }
