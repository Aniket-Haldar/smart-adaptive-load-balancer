package routing

import (
	"log"
	"github.com/sbirmecha99/smart-adaptive-load-balancer/internal/core"
)

type AdaptiveRouter struct {
	rr *RoundRobinRouter
	lc *LeastConnectionsRouter
	rn *RandomRouter
	
	currentAlgo string
	reason       string
	lastPicked   string
}

func NewAdaptiveRouter() *AdaptiveRouter {
	return &AdaptiveRouter{
		rr: NewRoundRobinRouter(),
		lc: NewLeastConnectionsRouter(),
		rn: NewRandomRouter(),
		currentAlgo: "roundrobin",
		reason:      "normal_conditions",
	}
}
func (ar *AdaptiveRouter) GetNextAvailableServer(
	backends []*core.Backend,
) *core.Backend {

	var totalConns int64
	var totalLatency int64
	var totalErrors int64

	aliveCount := 0

	for _, b := range backends {
		b.Mutex.Lock()
		if b.Alive {
			aliveCount++
			totalConns += b.ActiveConns
			totalLatency += int64(b.Latency)
			totalErrors += b.ErrorCount
		}
		b.Mutex.Unlock()
	}

	if aliveCount == 0 {
		log.Println("[ADAPTIVE] no alive backends")
		return nil
	}

	avgConns := totalConns / int64(aliveCount)
	avgLatency := totalLatency / int64(aliveCount)

	if totalErrors > 0 {
	ar.currentAlgo = "random"
	ar.reason = "errors_detected"
	b := ar.rn.GetNextAvailableServer(backends)
	ar.lastPicked = b.Address
	return b
}

if avgConns > 3 {
	ar.currentAlgo = "leastconnections"
	ar.reason = "high_concurrency"
	b := ar.lc.GetNextAvailableServer(backends)
	ar.lastPicked = b.Address
	return b
}

if avgLatency > 2_000_000 {
	ar.currentAlgo = "leastconnections"
	ar.reason = "high_latency"
	b := ar.lc.GetNextAvailableServer(backends)
	ar.lastPicked = b.Address
	return b
}

ar.currentAlgo = "roundrobin"
ar.reason = "normal_conditions"
b := ar.rr.GetNextAvailableServer(backends)
ar.lastPicked = b.Address
return b

}
func (ar *AdaptiveRouter) Name() string {
	if ar.currentAlgo != "" {
		return ar.currentAlgo // 🔹 frontend sees current algo dynamically
	}
	return "adaptive"
}
func (ar *AdaptiveRouter) CurrentAlgo() string { return ar.currentAlgo }
func (ar *AdaptiveRouter) Reason() string      { return ar.reason }
func (ar *AdaptiveRouter) LastPicked() string  { return ar.lastPicked }
