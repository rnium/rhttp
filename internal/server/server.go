package server

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/rnium/rhttp/internal/headers"
	"github.com/rnium/rhttp/internal/request"
	"github.com/rnium/rhttp/internal/response"
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
	res := response.NewResponse(req)
	headers := headers.NewHeaders()
	_ = headers.Set("server", "rhttp")
	_ = headers.Set("content-type", "application/json")
	_, err = res.WriteStatusLine(conn, response.StatusForbidden)
	if err != nil {
		fmt.Println(err)
	}
	_, err = res.WriteResponseHeaders(conn, headers)
	if err != nil {
		fmt.Println(err)
	}
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