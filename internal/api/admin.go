package api

import (
	"net/http"

	"github.com/sbirmecha99/smart-adaptive-load-balancer/internal/core"
)

func AddServerHandler(pool *core.ServerPool) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        addr := r.URL.Query().Get("address")
        if addr == "" {
            http.Error(w, "missing address", http.StatusBadRequest)
            return
        }

        pool.AddServer(&core.Backend{Address: addr, Alive: true})
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("server added"))
    }
}
