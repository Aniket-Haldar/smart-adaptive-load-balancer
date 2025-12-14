package proxy

import (
	"fmt"
	"net/http/httputil"
	"net/url"
)

// NewReverseProxy creates a reverse proxy for a backend address
func NewReverseProxy(address string) *httputil.ReverseProxy {

	backendURL, err := url.Parse(address)
	if err != nil {
		fmt.Println("Not a Proper URL")
	}

	// Single-host reverse proxy forwards all requests to this backend
	return httputil.NewSingleHostReverseProxy(backendURL)
}
