package sodibus

import "log"
import "github.com/sodibus/packet"
import "github.com/sodibus/sodibus/conn"

// Prepare a PacketReady
func (n *Node) ConnHandshake(c *conn.Conn, f *packet.PacketHandshake) (*packet.PacketReady, error) {
	p := &packet.PacketReady{
		Mode: f.Mode,
		NodeId: n.id,
		ClientId: c.GetId(),
	}
	return p, nil
}

// Add Conn to internal registry
func (n *Node) ConnDidStart(c *conn.Conn) {
	log.Println("New Conn: id =", c.GetId(), ", callee =", c.IsCallee(), ", provides =", c.GetProvides())
	// put to internal registry
	n.connsLock.Lock()
	n.conns[c.GetId()] = c
	n.connsLock.Unlock()
}

// Handle Frame
func (n *Node) ConnDidReceiveFrame(c *conn.Conn, f *packet.Frame) {
	go n.doConnDidReceiveFrame(c, f)
}

func (n *Node) doConnDidReceiveFrame(c *conn.Conn, f *packet.Frame) {
	m, err := f.Parse()
	if err != nil { return }
	switch m.(type) {
		case (*packet.PacketCallerSend): {
			p := m.(*packet.PacketCallerSend)
			log.Println("Invoke from", c.GetId(), ", callee_name =", p.Invocation.CalleeName , ", method =", p.Invocation.MethodName, ", arguments =", p.Invocation.Arguments)
			calleeId := n.ResolveCallee(p.Invocation.CalleeName)
			if calleeId == nil {
				log.Println("Callee named", p.Invocation.CalleeName, "not found")
				r, _ := packet.NewFrameWithPacket(&packet.PacketCallerRecv{
					Id: p.Id,
					Code: packet.ErrorCode_NO_CALLEE,
					Result: "",
				})
				c.Send(r)
			} else {
				r := &packet.PacketCalleeRecv{
					Id: &packet.InvocationId{
						Id: p.Id,
						ClientId: c.GetId(),
						NodeId: n.id,
					},
					Invocation: p.Invocation,
				}
				n.TransportInvocation(calleeId, r)
			}
		}
		case (*packet.PacketCalleeSend): {
			p := m.(*packet.PacketCalleeSend)
			n.TransportInvocationResult(p)
		}
	}
}

func (n *Node) ConnWillClose(c *conn.Conn, err error) {
	log.Println("Lost Conn: id =", c.GetId(), ", callee =", c.IsCallee(), ", err =", err)
	// remove from internal registry
	n.connsLock.Lock()
	delete(n.conns, c.GetId())
	n.connsLock.Unlock()
}

