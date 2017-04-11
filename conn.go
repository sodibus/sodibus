package sodibus

import "net"
import "errors"
import "sync"
import "github.com/sodibus/packet"
import "github.com/golang/protobuf/proto"

type Conn struct {
	id uint64
	conn *net.TCPConn
	provides []string
	isCallee bool
	sendLock *sync.Mutex
}

type ConnHandler interface {
	GetNodeId() uint64
	ConnDidStart(c *Conn)
	ConnDidReceiveFrame(c *Conn, f *packet.Frame)
	ConnWillClose(c *Conn)
}

func NewConn(conn *net.TCPConn, id uint64) *Conn {
	return &Conn{
		id: id,
		conn: conn,
		sendLock: &sync.Mutex{},
	}
}

func (c *Conn) Run(h ConnHandler) {
	defer func(){
		c.Close(h)
	}()

	err := c.doHandshake(h)
	if err != nil { return }

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

func (c *Conn) Send(f *packet.Frame) error {
	var err error
	c.sendLock.Lock()
	err = f.Write(c.conn)
	c.sendLock.Unlock()
	return err
}

func (c *Conn) Close(h ConnHandler) {
	h.ConnWillClose(c)
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

	r, err := packet.NewFrameWithPacket(&packet.PacketReady{
		Mode: p.Mode,
		NodeId: h.GetNodeId(),
		ClientId: c.id,
	})
	if err != nil { return err }

	err = c.Send(r)
	if err != nil { return err }

	c.isCallee = p.Mode == packet.ClientMode_CALLEE
	c.provides = p.Provides

	h.ConnDidStart(c)

	return err
}
