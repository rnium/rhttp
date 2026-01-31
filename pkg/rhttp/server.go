package rhttp

import (
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net"
	"strings"
	"sync"
)

type serverStatus int8

const (
	serverIdle serverStatus = iota
	serverRunning
	serverClosed
)

type Server struct {
	status    serverStatus
	listener  net.Listener
	router    *Router
	wg        sync.WaitGroup
	closeOnce sync.Once
	sigQuit   chan bool
}

func (s *Server) Close() error {
	var err error
	s.closeOnce.Do(func() {
		if s.listener != nil {
			err = s.listener.Close()
		}
		close(s.sigQuit)
		s.wg.Wait()
		s.status = serverClosed
	})
	return err
}

func (s *Server) runHandler(handler Handler, req *Request) (res *Response, err error) {
	func() {
		defer func() {
			if r := recover(); r != nil {
				switch v := r.(type) {
				case string:
					err = errors.New(v)
				case error:
					err = v
				default:
					err = errors.New(fmt.Sprint(v))
				}
			}
		}()
		res = handler(req)
	}()
	return
}

func (s *Server) runHttpCycle(conn io.ReadWriter, done chan bool) {
	for {
		req, err := getRequest(conn)
		if err != nil {
			break
		}

		handler := s.router.getHandler(req)
		res, err := s.runHandler(handler, req)
		if err != nil {
			res = response500(err)
		}
		_, err = res.writeResponse(conn, req)
		if err != nil {
			break
		}

		slog.Info(
			fmt.Sprintf(
				"%d %s %s\n",
				res.StatusCode,
				req.RequestLine.Method,
				req.RequestLine.Target,
			),
		)

		connHeader, ok := req.Headers.Get("connection")
		if ok && strings.EqualFold(connHeader, "close") {
			break
		}
	}
	done <- true
}

func (s *Server) handleConn(conn io.ReadWriteCloser) {
	defer func() {
		conn.Close()
		s.wg.Done()
	}()

	done := make(chan bool)
	go s.runHttpCycle(conn, done)

	select {
	case <-s.sigQuit:
	case <-done:
	}
}

func (s *Server) acceptConnections() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			return
		}
		s.wg.Add(1)
		go s.handleConn(conn)
	}
}

func NewServer(router *Router) *Server {
	return &Server{
		status:  serverIdle,
		router:  router,
		sigQuit: make(chan bool),
	}
}

func (s *Server) Start(port uint16) {
	if s.status == serverRunning {
		log.Fatalln("Server is already running")
	}
	s.status = serverRunning
	fmt.Println("Starting Server...")
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Cannot initiate listener on port %d\n", port)
	}
	s.listener = listener
	go s.acceptConnections()
	fmt.Printf("Listening for tcp connections on port %d\n", port)
}