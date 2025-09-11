package server

import (
	"fmt"
	"net"
	"sync/atomic"

	"github.com/blakehulett7/httpfromtcp/internal/response"
)

type Server struct {
	listener net.Listener
	isClosed atomic.Bool
}

func Serve(port int) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return &Server{}, err
	}

	server := &Server{
		listener: listener,
	}

	go server.listen()

	return server, nil
}

func (s *Server) Close() error {
	s.isClosed.Swap(true)
	return s.listener.Close()
}

func (s *Server) listen() {
	for {
		if s.isClosed.Load() {
			fmt.Println("leaving")
			return
		}

		conn, err := s.listener.Accept()
		if err != nil {
			fmt.Println(err)
		}

		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	response.WriteStatusLine(conn, response.StatusOK)
	response.WriteHeaders(conn, response.GetDefaultHeaders(0))
}
