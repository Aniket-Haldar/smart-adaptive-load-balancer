package core

import (
	"sync"
)

type ServerPool struct{
	servers []*Backend
	mu sync.RWMutex
}

func NewServerPool() *ServerPool{
	return &ServerPool{}
}

func (p *ServerPool) AddServer(b *Backend){
	p.mu.Lock()
	defer p.mu.Unlock()
	p.servers=append(p.servers, b)
}

func (p *ServerPool) GetServers() []*Backend{
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.servers
}