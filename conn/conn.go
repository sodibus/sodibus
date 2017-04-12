package conn

import "net"
import "sync"

// Conn represents a connection between Node and Caller/Callee
type Conn struct {
	// unique id in Node
	id uint64
	// underlaying *net.TCPConn
	conn *net.TCPConn
	// logic delegate
	delegate Delegate
	// callee provides
	provides []string
	// if this connection is a callee
	isCallee bool
	// mutex lock on sending
	sendLock *sync.Mutex
}

// Create a New Conn
func New(conn *net.TCPConn, id uint64) *Conn {
	return &Conn{
		id: id,
		conn: conn,
		sendLock: &sync.Mutex{},
	}
}

// Set the delegate of a Conn
//
// delegate muste be set before #Run()
func (c *Conn) SetDelegate(delegate Delegate) {
	c.delegate = delegate
}

// Get the ID of a Conn
func (c *Conn) GetId() uint64 {
	return c.id
}

// Get Callee Provides
func (c *Conn) GetProvides() []string {
	return c.provides
}

// If this Conn is a Callee
func (c *Conn) IsCallee() bool {
	return c.isCallee
}

