package sodibus

import "net"
import "errors"
import "github.com/sodibus/packet"
import "github.com/golang/protobuf/proto"

type Conn struct {
	id uint64
	conn *net.TCPConn
	provides []string
	isCallee bool

	sendChan chan *packet.Frame
	stopChan chan bool
}

type ConnHandler interface {
	ConnDidStart(c *Conn)
	ConnDidReceiveFrame(c *Conn, f *packet.Frame)
	ConnWillClose(c *Conn)
}

func NewConn(conn *net.TCPConn, id uint64) *Conn {
	return &Conn{
		id: id,
		conn: conn,
		sendChan: make(chan *packet.Frame, 64),
		stopChan: make(chan bool, 1),
	}
}

func (c *Conn) recvLoop(h ConnHandler) {
	defer func(){
		c.Close(h)
	}()
	for {
		f, err := packet.ReadFrame(c.conn)
		if err != nil {
			_, ok := err.(packet.UnsynchronizedError)
			if ok { continue } else { break }
		} else {
			go h.ConnDidReceiveFrame(c, f)
		}
	}
}

func (c *Conn) sendLoop(h ConnHandler) {
	select {
		case f := <- c.sendChan: {
			f.Write(c.conn)
		}
		case _ = <- c.stopChan: {
			c.doClose()
			break
		}
	}
}

func (c *Conn) Run(h ConnHandler) {
	err := c.doHandshake(h)
	if err != nil {
		c.Close(h)
	}
	go c.sendLoop(h)
	c.recvLoop(h)
}

func (c *Conn) Send(f *packet.Frame) {
	c.sendChan <- f
}

func (c *Conn) Close(h ConnHandler) {
	h.ConnWillClose(c)
	c.stopChan <- true
}

func (c *Conn) doClose() {
	c.conn.Close()
}

func (c *Conn) doHandshake(h ConnHandler) error {
	f, err := packet.ReadFrame(c.conn)
	if err != nil { return err }

	var m proto.Message
	m, err = f.Parse()
	if err != nil { return err }

	p, ok := m.(*packet.PacketHandshake)
	if !ok { return errors.New("Not a Handshake Packet") }

	c.isCallee = p.Mode == packet.ClientMode_CALLEE
	c.provides = p.Provides

	h.ConnDidStart(c)

	return err
}
