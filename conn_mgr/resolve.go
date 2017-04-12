package conn_mgr

import "github.com/Unknwon/com"
import "github.com/sodibus/sodibus/conn"

func (cm *ConnMgr) ResolveCallee(name string) *conn.Conn {
	var c *conn.Conn
	cm.connsLock.RLock()
	for _, v := range cm.conns {
		if v.IsCallee() && com.IsSliceContainsStr(v.GetProvides(), name) {
			c = v
			break
		}
	}
	cm.connsLock.RUnlock()
	return c
}
