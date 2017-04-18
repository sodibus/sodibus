package cluster

import "net"

// Server cluster server, receving requests from other nodes
type Server struct {
	addr     string
	listener net.Listener
}

// NewServer create a new cluster server
func NewServer(addr string) *Server {
	return &Server{
		addr: addr,
	}
}

// Start start the main loop
func (s *Server) Start() {
	go s.run()
}

func (s *Server) run() {
}
