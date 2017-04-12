package callee_mgr

import "sync"
import "log"

type Manager struct {
	groups map[string]*Group
	groupsLock *sync.RWMutex
}

func New() *Manager {
	return &Manager{
		groups: make(map[string]*Group),
		groupsLock: &sync.RWMutex{},
	}
}

func (m *Manager) Resolve(name string) *CalleeId {
	return m.Group(name).Take()
}

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

func (m *Manager) BatchPut(id CalleeId, provides []string) {
	log.Println("CalleeMgr Batch Put", id.ClientId, provides)
	for _, name := range provides {
		m.Group(name).Put(id)
	}
}

func (m *Manager) BatchDel(id CalleeId, provides []string) {
	log.Println("CalleeMgr Batch Del", id.ClientId, provides)
	for _, name := range provides {
		m.Group(name).Del(id)
	}
}

