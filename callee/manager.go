package callee

import "sync"

// CalleeId
//
// composed by a NodeId and a ClientId (node-local id)
type CalleeId struct {
	NodeId uint64
	ClientId uint64
}

// Manager
//
// Manager holds all the CalleeIds existed in the whole cluster, providing Round-Robin load balancing by underlaying Group struct
type Manager struct {
	// underlaying groups
	groups map[string]*Group
	// RWMutex lock for alternating groups
	groupsLock *sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{
		groups: make(map[string]*Group),
		groupsLock: &sync.RWMutex{},
	}
}

// Resolve a CalleeId for a CalleeName
func (m *Manager) Resolve(name string) *CalleeId {
	return m.Group(name).Take()
}

// Find or create a Group with CalleeName
func (m *Manager) Group(name string) *Group {
	var g *Group
	// find
	m.groupsLock.RLock()
	g = m.groups[name]
	m.groupsLock.RUnlock()

	if g == nil {
		// second find
		m.groupsLock.Lock()
		g = m.groups[name]

		// create
		if g == nil {
			g = NewGroup(name)
			m.groups[name] = g
		}
		m.groupsLock.Unlock()
	}
	return g
}

// Put a CalleeId for mutiple CalleeNames
func (m *Manager) BatchPut(id CalleeId, provides []string) {
	for _, name := range provides {
		m.Group(name).Put(id)
	}
}

// Delete a CalleeId from mutiple CalleeNames
func (m *Manager) BatchDel(id CalleeId, provides []string) {
	for _, name := range provides {
		m.Group(name).Del(id)
	}
}

