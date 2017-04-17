package sodibus

import "log"
import "github.com/sodibus/packet"
import "github.com/sodibus/sodibus/conn"
import "github.com/sodibus/sodibus/callee"

// ConnHandshake provides handshake logic
func (n *Node) ConnHandshake(c *conn.Conn, f *packet.PacketHandshake) (*packet.PacketReady, error) {
	p := &packet.PacketReady{
		Mode:     f.Mode,
		NodeId:   n.id,
		ClientId: c.GetID(),
	}
	return p, nil
}

// ConnDidStart provides logic on Conn did successfully complete handshake
func (n *Node) ConnDidStart(c *conn.Conn) {
	log.Println("New Conn: id =", c.GetID(), ", callee =", c.IsCallee(), ", provides =", c.GetProvides())
	// put to internal registry
	n.connMgr.Put(c)
	// put to callee manager
	if c.IsCallee() {
		n.calleeMgr.BatchPut(callee.FullID{NodeID: n.id, ClientID: c.GetID()}, c.GetProvides())
	}
}

// ConnDidReceivePacket handles frame
func (n *Node) ConnDidReceivePacket(c *conn.Conn, m packet.Packet) {
	go n.doConnDidReceivePacket(c, m)
}

func (n *Node) doConnDidReceivePacket(c *conn.Conn, m packet.Packet) {
	switch m.(type) {
	case (*packet.PacketCallerSend):
		{
			p := m.(*packet.PacketCallerSend)
			var callee *conn.Conn
			calleeID := n.calleeMgr.Resolve(p.Invocation.CalleeName)
			if calleeID != nil {
				callee = n.connMgr.Get(calleeID.ClientID)
			}
			if callee == nil {
				log.Println("Callee named", p.Invocation.CalleeName, "not found")
				r, _ := packet.NewFrameWithPacket(&packet.PacketCallerRecv{
					Id:     p.Id,
					Code:   packet.ErrorCode_NO_CALLEE,
					Result: "",
				})
				c.Send(r)
			} else {
				f, _ := packet.NewFrameWithPacket(&packet.PacketCalleeRecv{
					Id: &packet.InvocationId{
						Id:       p.Id,
						ClientId: c.GetID(),
						NodeId:   n.id,
					},
					Invocation: p.Invocation,
				})
				callee.Send(f)
			}
		}
	case (*packet.PacketCalleeSend):
		{
			p := m.(*packet.PacketCalleeSend)
			n.TransportInvocationResult(p)
		}
	}
}

// ConnWillClose removes Conn from connMgr and from calleeMgr if it's a Callee
func (n *Node) ConnWillClose(c *conn.Conn, err error) {
	log.Println("Lost Conn: id =", c.GetID(), ", callee =", c.IsCallee(), ", err =", err)
	// remove from internal registry
	n.connMgr.Del(c.GetID())
	// remove from callee manager
	if c.IsCallee() {
		n.calleeMgr.BatchDel(callee.FullID{NodeID: n.id, ClientID: c.GetID()}, c.GetProvides())
	}
}
