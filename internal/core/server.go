package core

import "sync"

type Server struct {
	URL         string
	Alive       bool
	Weight      int
	Connections int
	mutex       sync.Mutex
}
