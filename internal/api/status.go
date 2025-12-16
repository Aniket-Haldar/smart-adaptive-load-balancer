// file: internal/api/status.go
package api

import (
    "encoding/json"
    "net/http"
    "github.com/sbirmecha99/smart-adaptive-load-balancer/internal/core"
    "github.com/sbirmecha99/smart-adaptive-load-balancer/internal/routing"
)

type StatusResponse struct {
    CurrentAlgo      string          `json:"current_algo"`
    AdaptiveReason   string          `json:"adaptive_reason"`
    SelectedBackend  string          `json:"selected_backend"`
    Backends         []*core.Backend `json:"backends"`
}

func StatusHandler(router routing.Router, pool []*core.Backend) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        algoName := "unknown"
        adaptiveReason := ""
        selectedBackend := ""

        switch rt := router.(type) {
        case *routing.RoundRobinRouter:
            algoName = "roundrobin"
        case *routing.LeastConnectionsRouter:
            algoName = "leastconnections"
        case *routing.RandomRouter:
            algoName = "random"
        case *routing.AdaptiveRouter:
            algoName = "adaptive"
            adaptiveReason = rt.Reason()
            selectedBackend = rt.LastPicked()
        }

        resp := StatusResponse{
            CurrentAlgo:     algoName,
            AdaptiveReason:  adaptiveReason,
            SelectedBackend: selectedBackend,
            Backends:        pool,
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(resp)
    })
}
