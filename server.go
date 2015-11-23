package tftp

import (
	"fmt"
	"net"
)

type Server struct {
	conn *net.UDPConn
}

func (s *Server) Listen(addr *net.UDPAddr) (err error) {
	s.conn, err = net.ListenUDP("udp", addr)
	if err != nil {
		return
	}
	return
}

func (s *Server) Close() {
	s.conn.Close()
}

func (s *Server) Serve() {

	buf := make([]byte, 2048)
	for {
		n, addr, err := s.conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error:", err.Error())
		} else {
			fmt.Println("Received", n, "p from", addr.String())
			fmt.Printf("% x", buf[:n])
			fmt.Println()
		}
	}
}
