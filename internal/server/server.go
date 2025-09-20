// Server represents our TCP server
package server

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
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

type Handler = func(io.Writer, request.Request) *HandlerError

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
		writeHandlerErr(conn, &HandlerError{
			StatusCode: 400,
			Message:    "Bad Request\n",
		})
		return
	}

	// Prepare a buffer for the response body
	var body bytes.Buffer

	// Call the handler
	if herr := handler(&body, *req); herr != nil {
		writeHandlerErr(conn, herr)
		return
	}

	// Use bufio.Writer to efficiently write headers + body
	writer := bufio.NewWriter(conn)

	// Write status line
	_ = response.WriteStatusLine(writer, response.StatusCode(200))

	// Write headers (Content-Length includes body)
	headers := response.GetDefaultHeaders(body.Len())
	_ = response.WriteHeaders(writer, headers)

	// Write body
	writer.Write(body.Bytes())

	// Flush everything to the client
	writer.Flush()
}

// writeHandlerErr writes error responses consistently
func writeHandlerErr(conn net.Conn, herr *HandlerError) {

	// Write status line
	_ = response.WriteStatusLine(conn, response.StatusCode(herr.StatusCode))

	// Write headers
	headers := response.GetDefaultHeaders(len(herr.Message))
	_ = response.WriteHeaders(conn, headers)

	// Write error body
	conn.Write([]byte(herr.Message))

	// Flush to client
}
