package api

import (
	"encoding/json"
	"net/http"

	"github.com/sbirmecha99/smart-adaptive-load-balancer/internal/core"
)

func MetricsHandler(pool []*core.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(pool)
	}
}
