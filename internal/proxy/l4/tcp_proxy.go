package l4

import (
	"io"
	"net"
	"time"

	"github.com/sbirmecha99/smart-adaptive-load-balancer/internal/core"
	"github.com/sbirmecha99/smart-adaptive-load-balancer/internal/routing"
)

type TCPProxy struct {
	Pool   []*core.Backend
	Router routing.Router
}

func (p *TCPProxy) Start(listenAddr string) error {
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return err
	}

	for {
		clientConn, err := listener.Accept()
		if err != nil {
			continue
		}

		go p.handleConnection(clientConn)
	}
}

func (p *TCPProxy) handleConnection(clientConn net.Conn) {
	start := time.Now()

	backend := p.Router.GetNextAvailableServer(p.Pool)
	if backend == nil {
		clientConn.Close()
		return
	}

	backend.Mutex.Lock()
	backend.ActiveConns++
	backend.Mutex.Unlock()

	serverConn, err := net.Dial("tcp", backend.Address)
	if err != nil {
		backend.Mutex.Lock()
		backend.ActiveConns--
		backend.ErrorCount++
		backend.Mutex.Unlock()
		clientConn.Close()
		return
	}

	// Bidirectional copy
	go io.Copy(serverConn, clientConn)
	io.Copy(clientConn, serverConn)

	clientConn.Close()
	serverConn.Close()

	backend.Mutex.Lock()
	backend.ActiveConns--
	lat := time.Since(start)
	backend.Latency = (backend.Latency + lat) / 2

	backend.Mutex.Unlock()
}
