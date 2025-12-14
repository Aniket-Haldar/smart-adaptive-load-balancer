package api

import (
	"encoding/json"
	"net/http"

	"github.com/sbirmecha99/smart-adaptive-load-balancer/internal/core"
)

func AddServerHandler(pool *[]*core.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var b core.Backend
		if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		b.Alive = true
		*pool = append(*pool, &b)
		w.WriteHeader(http.StatusCreated)
	}
}
