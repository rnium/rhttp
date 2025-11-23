package server

import (
	"fmt"
	"io"
	"log"
	"log/slog"
	"net"

	"github.com/rnium/rhttp/internal/request"
	"github.com/rnium/rhttp/internal/router"
)

type Server struct {
	listener net.Listener
	router   *router.Router
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
	handler := s.router.GetHandler(req)
	res := handler(req)
	_, err = res.WriteResponse(conn, req)
	if err != nil {
		fmt.Println(err)
		return
	}
	slog.Info(
		fmt.Sprintf(
			"%d %s %s\n",
			res.StatusCode,
			req.RequestLine.Method,
			req.RequestLine.Target,
		),
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

func Serve(port uint16, router *router.Router) *Server {
	fmt.Println("Starting Server...")
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Cannot initiate listener on port %d\n", port)
	}
	server := &Server{
		listener: listener,
		router:   router,
	}
	go server.acceptConnections()
	fmt.Printf("Listening for tcp connections on port %d\n", port)
	return server
}
