package routing

import (
	"fmt"
	"net/http"
)

type RoundRobinLB struct {
	servers []*Server
	next    int
}

func NewRoundRobinLB(servers []*Server) *RoundRobinLB {
	return &RoundRobinLB{
		servers: servers,
		next:    0,
	}
}

// GetNextAvailableServer returns the next available backend server in a round-robin manner.
func (lb *RoundRobinLB) GetNextAvailableServer() *Server {

	numServers := len(lb.servers)

	// Start searching from the next index
	start := lb.next

	for i := 0; i < numServers; i++ {
		serverIndex := (start + i) % numServers
		server := lb.servers[serverIndex]

		server.mutex.Lock()
		alive := server.Alive
		server.mutex.Unlock()

		if alive {
			// update the next index for next iteration
			lb.next = (serverIndex + 1) % numServers
			return server
		}
	}
	// No available servers found, return nil
	return nil
}
func (lb *RoundRobinLB) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	server := lb.GetNextAvailableServer()
	if server != nil {
		proxy := NewReverseProxy(server.URL)

		fmt.Print("server is ", server)
		// Set the logger for the ReverseProxy
		proxy.ServerHttp(w, r)
	} else {
		fmt.Print("server is ", server)
		// TODO:
		// Handle the case when no available server is found
	}
}
