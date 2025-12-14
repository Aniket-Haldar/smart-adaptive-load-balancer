package core

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type ReverseProxy struct {
	backendURL string
	proxy      *httputil.ReverseProxy
}

func NewReverseProxy(backendURL string) *ReverseProxy {
	backend, _ := url.Parse(backendURL)

	return &ReverseProxy{
		backendURL: backendURL,
		proxy:      httputil.NewSingleHostReverseProxy(backend),
	}
}

func (rp *ReverseProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Forwarding request to %s : %s\n", rp.backendURL, r.URL.Path)
	rp.proxy.ServeHTTP(w, r)
}
