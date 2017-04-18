package cluster

// Cluster represents information of a cluster
type Cluster struct {
	id       uint64    // id of current Node
	addr     string    // addr of current Node, will send to other Nodes
	server   *Server   // listening server
	registry *Registry // registry
}

// New create a new Cluster
func New(addr string, laddr string, id uint64) *Cluster {
	return &Cluster{
		id:     id,
		addr:   addr,
		server: NewServer(laddr),
	}
}

// Start start internal components
func (c *Cluster) Start() {
	c.server.Start()
}
