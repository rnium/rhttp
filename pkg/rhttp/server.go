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

type Server struct {
	listener  net.Listener
	router    *Router
	wg        sync.WaitGroup
	closeOnce sync.Once
}

func (s *Server) Close() error {
	var err error
	s.closeOnce.Do(func() {
		if s.listener != nil {
			err = s.listener.Close()
		}
		s.wg.Wait()
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

func (s *Server) handleConn(conn io.ReadWriteCloser) {
	defer func() {
		conn.Close()
		s.wg.Done()
	}()

	for {
		req, err := getRequest(conn)
		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			}
			return
		}

		handler := s.router.getHandler(req)
		res, err := s.runHandler(handler, req)
		if err != nil {
			res = response500(err)
		}
		_, err = res.writeResponse(conn, req)
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

		connHeader, ok := req.Headers.Get("connection")
		if ok && strings.EqualFold(connHeader, "close") {
			break
		}
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

func Serve(port uint16, router *Router) *Server {
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
