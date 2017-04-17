package conn

import "github.com/sodibus/packet"

// Delegate delegates some Conn logic to others
type Delegate interface {

	// handle handshake process
	ConnHandshake(c *Conn, f *packet.PacketHandshake) (*packet.PacketReady, error)

	// notify connection handshake successed
	ConnDidStart(c *Conn)

	// notify connection did receive a frame
	ConnDidReceiveFrame(c *Conn, f *packet.Frame)

	// notify connection closed, a optional associated error
	ConnWillClose(c *Conn, err error)
}
