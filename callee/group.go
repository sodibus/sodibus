package callee

import "sync"

// Group is a group of Callees share a same Callee Name
type Group struct {
	// name
	name string

	// array of FullIDs
	ids []FullID
	// map of FullIDs
	idsMap map[FullID]bool
	// RWMutex to protect mutation
	idsLock *sync.RWMutex

	// cursor of round-robin balacing
	cursor int
	// lock of cursor
	cursorLock *sync.Mutex
}

// NewGroup creates a new Group
func NewGroup(name string) *Group {
	return &Group{
		name:       name,
		ids:        make([]FullID, 0),
		idsMap:     make(map[FullID]bool),
		idsLock:    &sync.RWMutex{},
		cursorLock: &sync.Mutex{},
	}
}

// Put adds a FullID to the group, if not added before
func (g *Group) Put(id FullID) {
	// lock ids and cursor
	g.idsLock.Lock()
	defer func() {
		g.idsLock.Unlock()
	}()

	// check if already exists
	found := g.idsMap[id]

	if found {
		return
	}

	// add to map
	g.idsMap[id] = true

	// add to array
	g.ids = append(g.ids, id)
}

// Del deletes a FullID from the group, if added before
//
// this will also adjust round-robin cursor automatically
func (g *Group) Del(id FullID) {
	// lock ids and cursor
	g.idsLock.Lock()
	g.cursorLock.Lock()
	defer func() {
		g.cursorLock.Unlock()
		g.idsLock.Unlock()
	}()

	idx := -1

	// find index of FullID
	for i, v := range g.ids {
		if v == id {
			idx = i
		}
	}
	if idx < 0 {
		return
	}

	// remove from map
	delete(g.idsMap, id)

	// remove from array
	g.ids = append(g.ids[:idx], g.ids[idx+1:]...)

	// adjust cursor if needed
	if idx < g.cursor {
		g.cursor = g.cursor - 1
	}
}

// Take takes a FullID by Round-Robin balancing
//
// returns nil if nothing found
func (g *Group) Take() *FullID {
	// lock ids and cursor
	g.idsLock.RLock()
	g.cursorLock.Lock()
	defer func() {
		g.cursorLock.Unlock()
		g.idsLock.RUnlock()
	}()

	l := len(g.ids)

	// returns nil if empty
	if l == 0 {
		return nil
	}

	// reset cursor if exceeded
	if g.cursor >= l {
		g.cursor = 0
	}

	// get
	c := &g.ids[g.cursor]

	// move cursor on
	g.cursor = g.cursor + 1

	return c
}
