package callee_mgr

import "sync"
import "log"

type Group struct {
	name string

	ids []CalleeId
	idsMap map[CalleeId]bool
	idsLock *sync.RWMutex

	cursor int
	cursorLock *sync.Mutex
}

func NewGroup(name string) *Group {
	return &Group{
		name: name,
		ids: make([]CalleeId, 0),
		idsMap: make(map[CalleeId]bool),
		idsLock: &sync.RWMutex{},
		cursorLock: &sync.Mutex{},
	}
}

func (g *Group) Put(id CalleeId) {
	// lock ids and cursor
	g.idsLock.Lock()
	defer func() {
		g.idsLock.Unlock()
	}()

	// check if already exists
	found := g.idsMap[id]

	if found { return }

	// add to map
	g.idsMap[id] = true

	// add to array
	g.ids = append(g.ids, id)
}

func (g *Group) Del(id CalleeId) {
	// lock ids and cursor
	g.idsLock.Lock()
	g.cursorLock.Lock()
	defer func() {
		g.cursorLock.Unlock()
		g.idsLock.Unlock()
	}()

	idx := -1

	// find index of CalleeId
	for i, v := range g.ids {
		if v == id {
			idx = i
		}
	}
	if idx < 0 { return }

	// remove from map
	delete(g.idsMap, id)

	// remove from array
	g.ids = append(g.ids[:idx], g.ids[idx+1:]...)

	// adjust cursor if needed
	if idx < g.cursor {
		g.cursor = g.cursor - 1
	}
}

func (g *Group) Take() *CalleeId {
	log.Println("Take", g.name, "among", len(g.ids), "callees")
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
	if g.cursor >= l { g.cursor = 0 }

	// get
	c := &g.ids[g.cursor]

	// move cursor on
	g.cursor = g.cursor + 1

	return c
}

