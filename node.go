package sodibus

import "net"
import "log"
import "sync"
import "errors"
import "math/rand"
import "github.com/Unknwon/com"
import "github.com/sodibus/packet"
import "github.com/sodibus/sodibus/conn"

// locate a Callee across mutiple nodes

type CalleeId struct {
	NodeId uint64
	ClientId uint64
}

type Node struct {
	// node information
	id uint64
	addr string
	listener *net.TCPListener
	// connections
	lastConnId uint64
	conns map[uint64]*conn.Conn
	connsLock *sync.RWMutex
}

func NewNode(addr string) *Node {
	return &Node {
		id: rand.Uint64(),
		addr: addr,
		conns: make(map[uint64]*conn.Conn),
		connsLock: &sync.RWMutex{},
	}
}

// Conn Management

func (n *Node) NewConnId() uint64 {
	n.lastConnId = n.lastConnId + 1
	return n.lastConnId
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
			c := conn.New(cn, n.NewConnId(), n)
			// start Conn
			go c.Run()
		} else {
			log.Fatal("Failed to accept", err)
			return err
		}
	}
}

// Resolving

func (n *Node) ResolveCallee(name string) *CalleeId {
	var calleeId *CalleeId
	// find a usable client and send back
	n.connsLock.RLock()
	for _, v := range n.conns {
		if v.IsCallee() && com.IsSliceContainsStr(v.GetProvides(), name) {
			calleeId = &CalleeId{
				NodeId: n.id,
				ClientId: v.GetId(),
			}
			break
		}
	}
	n.connsLock.RUnlock()
	// send nil if nothing found
	return calleeId
}

func (n *Node) TransportInvocation(calleeId *CalleeId, p *packet.PacketCalleeRecv) error {
	conn := n.conns[calleeId.ClientId]
	if conn == nil { return errors.New("no callee found") }
	f, err := packet.NewFrameWithPacket(p)
	if err != nil { return err }
	err = conn.Send(f)
	return err
}

func (n *Node) TransportInvocationResult(p *packet.PacketCalleeSend) error {
	conn := n.conns[p.Id.ClientId]
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

// ConnHandler
