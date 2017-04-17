package conn

import "net"
import "sync"
import "sync/atomic"

// Manager manages mutiple Conns, assigns unique ids and find Conn by id
type Manager struct {
	seqID     uint64
	conns     map[uint64]*Conn
	connsLock *sync.RWMutex
}

// NewManager creates a new Manager
func NewManager() *Manager {
	return &Manager{
		conns:     make(map[uint64]*Conn, 128),
		connsLock: &sync.RWMutex{},
	}
}

// Wrap wraps a net.TCPConn as Conn, but not put to internal map yet.
func (cm *Manager) Wrap(cn *net.TCPConn) *Conn {
	id := atomic.AddUint64(&cm.seqID, 1)
	return New(cn, id)
}

// Put puts a Conn, can be queried later
func (cm *Manager) Put(c *Conn) {
	cm.connsLock.Lock()
	defer cm.connsLock.Unlock()
	cm.conns[c.GetID()] = c
}

// Del deletes a Conn, remove from internal map
func (cm *Manager) Del(id uint64) {
	cm.connsLock.Lock()
	defer cm.connsLock.Unlock()
	delete(cm.conns, id)
}

// Get finds a Conn by id
func (cm *Manager) Get(id uint64) *Conn {
	cm.connsLock.RLock()
	defer cm.connsLock.RUnlock()
	return cm.conns[id]
}
