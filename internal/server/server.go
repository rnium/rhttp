package server

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/rnium/rhttp/internal/request"
)

type Server struct{
	listener net.Listener
}

func (s *Server) Close() {
	s.listener.Close()
}

func (s *Server) handleConn(conn io.ReadWriteCloser) {
	defer conn.Close()
	req, err := request.GetRequest(conn)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Method: %s, Target: %s, Version: %s\n", req.RequestLine.Method, req.RequestLine.Target, req.RequestLine.Version)
	req.Headers.ForEach(func(name, value string) {
		fmt.Printf("%s: %s\n", name, value)
	})
	fmt.Println(string(req.Body))
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