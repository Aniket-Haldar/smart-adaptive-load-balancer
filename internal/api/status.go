package api
import (
	"encoding/json"
	"net/http"

	"github.com/sbirmecha99/smart-adaptive-load-balancer/internal/core"
	"github.com/sbirmecha99/smart-adaptive-load-balancer/internal/routing"
)

// StatusResponse is what we return to the client
type StatusResponse struct {
	CurrentAlgo     string          `json:"current_algo"`
	AdaptiveReason  string          `json:"adaptive_reason"`
	SelectedBackend string          `json:"selected_backend"`
	Backends        []*core.Backend `json:"backends"`
}

// StatusHandler dynamically reports router status and backends
func StatusHandler(router *routing.AdaptiveRouter, getBackends func() []*core.Backend) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := StatusResponse{
			CurrentAlgo:     router.CurrentAlgo(),
			AdaptiveReason:  router.Reason(),
			SelectedBackend: router.LastPicked(),
			Backends:        getBackends(),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})
}   