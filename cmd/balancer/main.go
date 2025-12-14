package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/sbirmecha99/smart-adaptive-load-balancer/internal/api"
	"github.com/sbirmecha99/smart-adaptive-load-balancer/internal/core"
	"github.com/sbirmecha99/smart-adaptive-load-balancer/internal/health"
	"github.com/sbirmecha99/smart-adaptive-load-balancer/internal/proxy/l4"
	"github.com/sbirmecha99/smart-adaptive-load-balancer/internal/proxy/l7"
	"github.com/sbirmecha99/smart-adaptive-load-balancer/internal/routing"
)

func main() {
	mode := os.Getenv("LB_MODE")
	if mode == "" {
		mode = "L7"
	}

	//backend servers
	pool := []*core.Backend{
		{Address: "localhost:9001", Alive: true},
		{Address: "localhost:9002", Alive: true},
	}

	//router will come here
	router := routing.NewRoundRobinRouter()


	//metrics and admin
	http.Handle("/metrics", api.MetricsHandler(pool))
	http.Handle("/admin/add", api.AddServerHandler(&pool))

	//health checker
	checker:= &health.Checker{
		Backends: pool,
		Interval: 5*time.Second,
		Timeout: 2*time.Second,
	}
	checker.Start()
	//mode switch
	if mode == "L4" {
		log.Println("Starting L4 TCP Load Balancer on :8080")

		tcpProxy := &l4.TCPProxy{
			Pool:   pool,
			Router: router,
		}

		log.Fatal(tcpProxy.Start(":8080"))

	} else {
		log.Println("Starting L7 HTTP Load Balancer on :8080")

		httpProxy := &l7.HTTPProxy{
			Pool:   pool,
			Router: router,
		}

		http.Handle("/", httpProxy)
		log.Fatal(http.ListenAndServe(":8080", nil))
	}
}
