package sodibus

import "log"
import "github.com/sodibus/packet"
import "github.com/sodibus/sodibus/conn"
import "github.com/sodibus/sodibus/callee"

// Prepare a PacketReady
func (n *Node) ConnHandshake(c *conn.Conn, f *packet.PacketHandshake) (*packet.PacketReady, error) {
	p := &packet.PacketReady{
		Mode:     f.Mode,
		NodeId:   n.id,
		ClientId: c.GetId(),
	}
	return p, nil
}

// Add Conn to internal registry
func (n *Node) ConnDidStart(c *conn.Conn) {
	log.Println("New Conn: id =", c.GetId(), ", callee =", c.IsCallee(), ", provides =", c.GetProvides())
	// put to internal registry
	n.connMgr.Put(c)
	// put to callee manager
	if c.IsCallee() {
		n.calleeMgr.BatchPut(callee.CalleeId{NodeId: n.id, ClientId: c.GetId()}, c.GetProvides())
	}
}

// Handle Frame
func (n *Node) ConnDidReceiveFrame(c *conn.Conn, f *packet.Frame) {
	go n.doConnDidReceiveFrame(c, f)
}

func (n *Node) doConnDidReceiveFrame(c *conn.Conn, f *packet.Frame) {
	m, err := f.Parse()
	if err != nil {
		return
	}

	switch m.(type) {
	case (*packet.PacketCallerSend):
		{
			p := m.(*packet.PacketCallerSend)
			var callee *conn.Conn
			calleeId := n.calleeMgr.Resolve(p.Invocation.CalleeName)
			if calleeId != nil {
				callee = n.connMgr.Get(calleeId.ClientId)
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
						ClientId: c.GetId(),
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

// Handle Conn close
//
// remove Conn from connMgr and calleeMgr if it's a Callee
func (n *Node) ConnWillClose(c *conn.Conn, err error) {
	log.Println("Lost Conn: id =", c.GetId(), ", callee =", c.IsCallee(), ", err =", err)
	// remove from internal registry
	n.connMgr.Del(c.GetId())
	// remove from callee manager
	if c.IsCallee() {
		n.calleeMgr.BatchDel(callee.CalleeId{NodeId: n.id, ClientId: c.GetId()}, c.GetProvides())
	}
}
