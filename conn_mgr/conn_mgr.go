package conn_mgr

import "sync"
import "github.com/sodibus/sodibus/conn"

type ConnMgr struct {
	seqId uint64
	conns map[uint64]*conn.Conn
	connsLock *sync.RWMutex
}

func New() *ConnMgr {
	return &ConnMgr{
		conns: make(map[uint64]*conn.Conn, 128),
		connsLock: &sync.RWMutex{},
	}
}

