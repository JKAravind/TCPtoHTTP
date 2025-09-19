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

	fmt.Println("ready to read")

	req, err := request.RequestFromReader(conn)
	fmt.Println("read Over")

	if err != nil {
		reqError := &HandlerError{StatusCode: 400, Message: "print Bad Req"}
		writeHandlerErr(conn, reqError)
		s.closed.Store(true)
		return
	}
	fmt.Println("Received request:")
	var buf bytes.Buffer
	reqError := handler(&buf, *req)
	if reqError != nil {
		writeHandlerErr(conn, reqError)
		return
	}

	writer := bufio.NewWriter(conn)

	_ = response.WriteStatusLine(writer, response.StatusCode(200))

	headers := response.GetDefaultHeaders(buf.Len())
	_ = response.WriteHeaders(writer, headers)

	writer.Write(buf.Bytes())

	// No body
	writer.Flush()
}

func writeHandlerErr(conn io.Writer, handlerError *HandlerError) {

	_ = response.WriteStatusLine(conn, response.StatusCode(handlerError.StatusCode))
	headers := response.GetDefaultHeaders(0)
	_ = response.WriteHeaders(conn, headers)

}
