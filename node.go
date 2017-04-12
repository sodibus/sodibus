package sodibus

import "net"
import "log"
import "errors"
import "math/rand"
import "github.com/sodibus/packet"
import "github.com/sodibus/sodibus/conn_mgr"
import "github.com/sodibus/sodibus/callee_mgr"

// locate a Callee across mutiple nodes

type Node struct {
	// node information
	id uint64
	addr string
	listener *net.TCPListener
	// connections
	connMgr *conn_mgr.ConnMgr
	calleeMgr *callee_mgr.Manager
}

func NewNode(addr string) *Node {
	return &Node {
		id: rand.Uint64(),
		addr: addr,
		connMgr: conn_mgr.New(),
		calleeMgr: callee_mgr.New(),
	}
}

// Loops

func (n *Node) Run() error {
	// resolve TCP address to bind
	tcpAddr, err := net.ResolveTCPAddr("tcp", n.addr)
	if err != nil { return err }

	// create listener
	n.listener, err = net.ListenTCP("tcp", tcpAddr)
	if err != nil { return err }

	log.Println("SODIBus", n.id, "listening at", n.addr)

	// accepting
	for {
		// accept
		cn, err := n.listener.AcceptTCP()
		if err == nil {
			// create client, auto atomical id
			c := n.connMgr.Wrap(cn)
			c.SetDelegate(n)
			// start Conn
			go c.Run()
		} else {
			log.Fatal("Failed to accept", err)
			return err
		}
	}
}

// Resolving

func (n *Node) TransportInvocationResult(p *packet.PacketCalleeSend) error {
	conn := n.connMgr.Get(p.Id.ClientId)
	if conn == nil { return errors.New("no callee found") }
	f, err := packet.NewFrameWithPacket(&packet.PacketCallerRecv{
		Id: p.Id.Id,
		Code: packet.ErrorCode_OK,
		Result: p.Result,
	})
	if err != nil { return err }
	err = conn.Send(f)
	return err
}

