// Server represents our TCP server
package server

import (
	"fmt"
	"net"
	"strconv"
	"sync/atomic"

	"github.com/JKAravind/TCPtoHTTP/internal/request"
	"github.com/JKAravind/TCPtoHTTP/internal/response"
)

type Server struct {
	closed   atomic.Bool
	listener net.Listener
}

type HandlerError struct {
	StatusCode int
	Message    string
}

type Handler func(w *response.Writer, req *request.Request)

// Serve starts the server on the given port
func Serve(port int, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return nil, fmt.Errorf("listen error: %v", err)
	}

	server := &Server{
		closed:   atomic.Bool{},
		listener: listener,
	}
	server.closed.Store(false)

	// start accepting connections in a goroutine
	go server.listen(handler)

	return server, nil
}

// Close shuts down the server
func (s *Server) Close() error {
	s.closed.Store(true)
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

// listen accepts connections and handles each in its own goroutine
func (s *Server) listen(handler Handler) {
	for !s.closed.Load() {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				// expected error after server is closed
				return
			}
			fmt.Println("accept error:", err)
			continue
		}
		go s.handle(conn, handler)
	}
}

// handle writes a fixed HTTP response and closes the connection
func (s *Server) handle(conn net.Conn, handler Handler) {
	defer conn.Close()

	// Parse the request
	req, err := request.RequestFromReader(conn)
	if err != nil {
		fmt.Println(err)
		return
	}
	w := &response.Writer{Connection: conn}
	handler(w, req)

}

// writeHandlerErr writes error responses consistently
