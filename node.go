package sodibus

import (
	"fmt"
	"log"
	"net"

	"github.com/sodibus/packet"
	"github.com/sodibus/sodibus/callee"
)
import "math/rand"

import "github.com/sodibus/sodibus/conn"

// Node is a SODIBus server
//
// Node should have a unique Id, it holds mutiple Callee/Callers
type Node struct {
	// node information
	id       uint64
	addr     string
	listener *net.TCPListener
	// connections
	connMgr   *conn.Manager
	calleeMgr *callee.Manager
}

// NewNode create a new Node
func NewNode(addr string) *Node {
	return &Node{
		id:        rand.Uint64(),
		addr:      addr,
		connMgr:   conn.NewManager(),
		calleeMgr: callee.NewManager(),
	}
}

// Run is the main loop of Node
//
// should run in a goroutine or main routine
func (n *Node) Run() error {
	// resolve TCP address to bind
	tcpAddr, err := net.ResolveTCPAddr("tcp", n.addr)
	if err != nil {
		return err
	}

	// create listener
	n.listener, err = net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return err
	}

	log.Println("SODIBus", n.id, "listening at", n.addr)

	// accepting
	for {
		// accept
		cn, err := n.listener.AcceptTCP()
		if err == nil {
			cn.SetKeepAlive(true)
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

// TransportInvocationResult transports a CalleeSend packet to it's orignal Caller
//TODO: use cluster transport system
func (n *Node) TransportInvocationResult(p *packet.PacketCalleeSend) error {
	conn := n.connMgr.Get(p.Id.ClientId)
	if conn == nil {
		return fmt.Errorf("Callee Not Found with ID %v", p.Id)
	}
	f, err := packet.NewFrameWithPacket(&packet.PacketCallerRecv{
		Id:     p.Id.Id,
		Code:   packet.ErrorCode_OK,
		Result: p.Result,
	})
	if err != nil {
		return err
	}
	err = conn.Send(f)
	return err
}
