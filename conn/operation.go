package conn

import "errors"
import "github.com/sodibus/packet"
import "github.com/golang/protobuf/proto"

// Recv loop
//
// should run in a seperated goroutine
func (c *Conn) Run() {
	var err error

	// defer to close
	defer func(){
		c.Close(err)
	}()

	// execute handshake
	err = c.doHandshake()
	if err != nil { return }

	// read loop
	for {
		f, err := packet.ReadFrame(c.conn)
		if err != nil {
			_, ok := err.(packet.UnsynchronizedError)
			if ok { continue } else { break } // ignore UnsynchronizedError
		} else {
			c.delegate.ConnDidReceiveFrame(c, f)
		}
	}
}

// Send a Frame
//
// uses a Mutex internally
func (c *Conn) Send(f *packet.Frame) error {
	var err error
	c.sendLock.Lock()
	err = f.Write(c.conn)
	c.sendLock.Unlock()
	return err
}

// Close a Conn
//
// notify delegate and close underlaying connection
func (c *Conn) Close(err error) {
	c.delegate.ConnWillClose(c, err)
	c.conn.Close()
}

// Execute Handshake Process
//
// execute handshake process, using delegate internally
func (c *Conn) doHandshake() error {
	// read a frame
	f, err := packet.ReadFrame(c.conn)
	if err != nil { return err }

	// parse handshake packet
	var m proto.Message
	m, err = f.Parse()
	if err != nil { return err }

	p, ok := m.(*packet.PacketHandshake)
	if !ok { return errors.New("Not a Handshake Packet") }

	// handshake with delegate
	r, err := c.delegate.ConnHandshake(c, p)
	if err != nil { return err }

	rf, err := packet.NewFrameWithPacket(r)
	if err != nil { return err }

	err = c.Send(rf)
	if err != nil { return err }

	// update values
	c.isCallee = p.Mode == packet.ClientMode_CALLEE
	c.provides = p.Provides

	// notify delegate
	c.delegate.ConnDidStart(c)

	return err
}

