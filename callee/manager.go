package callee

import "sync"

// FullID is composed by a NodeID and a ClientID (node-local id)
type FullID struct {
	NodeID   uint64
	ClientID uint64
}

// Manager holds all the FullIDs existed in the whole cluster, providing Round-Robin load balancing by underlaying Group struct
type Manager struct {
	// underlaying groups
	groups map[string]*Group
	// RWMutex lock for alternating groups
	groupsLock *sync.RWMutex
}

// NewManager create a new Manager
func NewManager() *Manager {
	return &Manager{
		groups:     make(map[string]*Group),
		groupsLock: &sync.RWMutex{},
	}
}

// Resolve a FullID for a CalleeName
func (m *Manager) Resolve(name string) *FullID {
	return m.Group(name).Take()
}

// Group find or create a Group with CalleeName
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

// BatchPut put a FullID for mutiple CalleeNames
func (m *Manager) BatchPut(id FullID, provides []string) {
	for _, name := range provides {
		m.Group(name).Put(id)
	}
}

// BatchDel delete a FullID from mutiple CalleeNames
func (m *Manager) BatchDel(id FullID, provides []string) {
	for _, name := range provides {
		m.Group(name).Del(id)
	}
}
