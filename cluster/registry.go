package cluster

import "github.com/sodibus/packet"

// Registry registry of existing Callees
type Registry struct {
}

// NewRegistry create a new Registry
func NewRegistry() *Registry {
	return &Registry{}
}

// Put add a callee to addr
func (r *Registry) Put(id packet.CalleeId, provides []string) {
}

// Del remove a callee by id
func (r *Registry) Del(id packet.CalleeId) {
}

// DelByNodeID remove all callee attached with specified host
func (r *Registry) DelByNodeID(id uint64) {
}

// Get get a node by id
func (r *Registry) Get(id uint64) {
}
