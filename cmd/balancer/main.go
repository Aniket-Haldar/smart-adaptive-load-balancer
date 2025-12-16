package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"

	"github.com/sbirmecha99/smart-adaptive-load-balancer/internal/api"
	"github.com/sbirmecha99/smart-adaptive-load-balancer/internal/core"
	"github.com/sbirmecha99/smart-adaptive-load-balancer/internal/health"
	"github.com/sbirmecha99/smart-adaptive-load-balancer/internal/proxy/l4"
	"github.com/sbirmecha99/smart-adaptive-load-balancer/internal/proxy/l7"
	"github.com/sbirmecha99/smart-adaptive-load-balancer/internal/routing"
)

func startDummyBackend(port string) {
	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello from backend %s", port)
		})
		log.Printf("Starting dummy backend on %s\n", port)
		log.Fatal(http.ListenAndServe(":"+port, mux))
	}()
}

func main() {
	mode := os.Getenv("LB_MODE")
	if mode == "" {
		mode = "L7"
	}

	// Start dummy backends
	startDummyBackend("9001")
	//startDummyBackend("9002")
	startDummyBackend("9003")

	// Backend pool
	pool := []*core.Backend{
		{Address: "localhost:9001", Alive: true},
		{Address: "localhost:9002", Alive: true},
		{Address: "localhost:9003", Alive: true},
	}

	router := routing.NewRoundRobinRouter()

	// 🔹 SINGLE mux for everything
	mux := http.NewServeMux()

	// Metrics & admin
	mux.Handle("/metrics", api.MetricsHandler(pool))
	mux.Handle("/admin/add", api.AddServerHandler(&pool))

	// Health checker
	checker := &health.Checker{
		Backends: pool,
		Interval: 5 * time.Second,
		Timeout:  2 * time.Second,
	}
	checker.Start()

	if mode == "L4" {
		log.Println("Starting L4 TCP Load Balancer on :8080")
		tcpProxy := &l4.TCPProxy{
			Pool:   pool,
			Router: router,
		}
		log.Fatal(tcpProxy.Start(":8080"))
		return
	}

	// L7 HTTP Proxy
	log.Println("Starting L7 HTTP Load Balancer on :8080")
	httpProxy := &l7.HTTPProxy{
		Pool:   pool,
		Router: router,
	}
	mux.Handle("/", httpProxy)

	// 🔹 Apply CORS to entire mux
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type"}),
	)(mux)

	log.Fatal(http.ListenAndServe(":8080", corsHandler))
}
