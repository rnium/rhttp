package server

import (
	"fmt"
	"io"
	"log"
	"log/slog"
	"net"

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
	data := []byte("hello world")
	res := response.NewResponse(response.StatusBadRequest, data, nil)
	n, err := res.WriteResponse(conn, req)
	if err != nil {
		fmt.Println(err)
	}
	slog.Info(
		fmt.Sprintf("Written %d bytes\n", n),
	)
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
	fmt.Println("Starting Server...")
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