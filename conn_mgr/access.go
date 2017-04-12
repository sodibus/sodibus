package conn_mgr

import "net"
import "sync/atomic"
import "github.com/sodibus/sodibus/conn"

func (cm *ConnMgr) Wrap(cn *net.TCPConn) *conn.Conn {
	id := atomic.AddUint64(&cm.seqId, 1)
	return conn.New(cn, id)
}

func (cm *ConnMgr) Put(c *conn.Conn) {
	cm.connsLock.Lock()
	cm.conns[c.GetId()] = c
	cm.connsLock.Unlock()
}

func (cm *ConnMgr) Del(id uint64) {
	cm.connsLock.Lock()
	delete(cm.conns, id)
	cm.connsLock.Unlock()
}

func (cm *ConnMgr) Get(id uint64) *conn.Conn {
	var c *conn.Conn
	cm.connsLock.RLock()
	c = cm.conns[id]
	cm.connsLock.RUnlock()
	return c
}
