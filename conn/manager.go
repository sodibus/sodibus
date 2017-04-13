package conn

import "net"
import "sync"
import "sync/atomic"

// Manager
//
// Manager manages mutiple Conn, assigns unique ids and find Conn by id
type Manager struct {
	seqId uint64
	conns map[uint64]*Conn
	connsLock *sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{
		conns: make(map[uint64]*Conn, 128),
		connsLock: &sync.RWMutex{},
	}
}

// Wrap a net.TCPConn as Conn
func (cm *Manager) Wrap(cn *net.TCPConn) *Conn {
	id := atomic.AddUint64(&cm.seqId, 1)
	return New(cn, id)
}

// Put a Conn, can be queried later
func (cm *Manager) Put(c *Conn) {
	cm.connsLock.Lock()
	defer cm.connsLock.Unlock()
	cm.conns[c.GetId()] = c
}

// Delete a Conn, remove from internal map
func (cm *Manager) Del(id uint64) {
	cm.connsLock.Lock()
	defer cm.connsLock.Unlock()
	delete(cm.conns, id)
}

// Find a Conn by id
func (cm *Manager) Get(id uint64) *Conn {
	cm.connsLock.RLock()
	defer cm.connsLock.RUnlock()
	return cm.conns[id]
}
