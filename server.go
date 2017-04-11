package sodibus

import "net"
import "math/rand"
import "sync/atomic"
import _ "github.com/sodibus/packet"

type Client struct {
	Id uint64
	inited bool
	conn *net.TCPConn
	callee bool
	provides []string
}

type Server struct {
	seqId uint64
	Id uint64
	Addr string
	clients map[uint64]*Client
	listener *net.TCPListener
}

func NewServer(addr string) *Server {
	return &Server{ Addr: addr, Id: rand.Uint64() }
}

func (s *Server) Run() error {
	tcpAddr, err := net.ResolveTCPAddr("tcp", s.Addr)
	if err != nil { return err }
	s.listener, err = net.ListenTCP("tcp", tcpAddr)
	if err != nil { return err }
	s.listenLoop()
	return err
}

func (s *Server) listenLoop() {
	for {
		conn, err := s.listener.AcceptTCP()
		if err != nil {
			println("Failed to Accept: ", err)
		} else {
			id := atomic.AddUint64(&s.seqId, 1)
			s.clients[id] = &Client{ Id: id, conn: conn }
			go s.handleLoop(conn, id)
		}
	}
}

func (s *Server) handleLoop(conn *net.TCPConn, id uint64) {
	// read loop
	for {
	}
}
