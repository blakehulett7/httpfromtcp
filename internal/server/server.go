package server

import (
	"bytes"
	"fmt"
	"net"
	"sync/atomic"

	"github.com/blakehulett7/httpfromtcp/internal/request"
	"github.com/blakehulett7/httpfromtcp/internal/response"
)

type Server struct {
	listener net.Listener
	isClosed atomic.Bool
}

func Serve(port int, h Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return &Server{}, err
	}

	server := &Server{
		listener: listener,
	}

	go server.listen(h)

	return server, nil
}

func (s *Server) Close() error {
	s.isClosed.Store(true)
	return s.listener.Close()
}

func (s *Server) listen(h Handler) {
	for {
		conn, err := s.listener.Accept()

		if err != nil {
			if s.isClosed.Load() {
				fmt.Println("leaving")
				return
			}

			fmt.Println(err)
			continue
		}

		go s.handle(conn, h)
	}
}

func (s *Server) handle(conn net.Conn, h Handler) {
	defer conn.Close()

	r, err := request.RequestFromReader(conn)
	if err != nil {
		writeError(conn, &HandlerError{StatusCode: response.StatusBadRequest, Error: err})
		return
	}

	buffer := bytes.NewBuffer(nil)
	h_err := h(buffer, r)
	if h_err != nil {
		writeError(conn, h_err)
		return
	}

	data := buffer.Bytes()

	response.WriteStatusLine(conn, response.StatusOK)
	response.WriteHeaders(conn, response.GetDefaultHeaders(len(data)))
	conn.Write(data)
}
