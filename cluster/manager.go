package cluster

import "net"

type Manager struct {
	id uint64
	addr string
	listener net.Listener
}

