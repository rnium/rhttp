package server

import (
	"fmt"
	"io"
	"log"
	"net"
)

type Server struct{
	listener net.Listener
}

func (s *Server) Close() {
	s.listener.Close()
}

func (s *Server) handleConn(conn io.ReadWriteCloser) {
	defer conn.Close()
}

func (s *Server) acceptConnections() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			return
		}
		go s.handleConn(conn)
	}
}


func Serve(port uint16) *Server {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Cannot initiate listener on port %d\n", port)
	}
	server := &Server{
		listener: listener,
	}
	go server.acceptConnections()
	fmt.Printf("Listening for tcp connections on port %d\n", port)
	return server
}