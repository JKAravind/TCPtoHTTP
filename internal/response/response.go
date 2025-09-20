package response

import (
	"fmt"
	"net"
	"strconv"

	"github.com/JKAravind/TCPtoHTTP/internal/headers"
)

type StatusCode int

type Writer struct {
	Connection net.Conn
}

const (
	StatusSuccess    StatusCode = 200
	StatusWrong      StatusCode = 400
	StatusServerFail StatusCode = 500
)

func (wr *Writer) WriteStatusLine(statusCode StatusCode) error {
	var write []byte
	switch statusCode {
	case StatusSuccess:
		fmt.Println("writing this to the rcurl client")
		write = []byte("HTTP/1.1 200 OK\r\n")
	case StatusWrong:
		write = []byte("HTTP/1.1 400 Bad Request\r\n")
	case StatusServerFail:
		write = []byte("HTTP/1.1 500 Internal Server Error\r\n")
	}
	n, err := wr.Connection.Write(write)
	fmt.Println("wrote bytes:", n, "err:", err)
	return err
}

func GetDefaultHeaders(contentLen int) headers.Header {
	responseHeader := headers.NewHeaders()
	responseHeader["Content-Length"] = strconv.Itoa(contentLen)
	responseHeader["Connection"] = "close"
	responseHeader["Content-Type"] = "text/html"

	return responseHeader
}
func (writer *Writer) WriteHeaders(headers headers.Header) error {

	for key, value := range headers {
		element := fmt.Sprintf("%s: %s\r\n", key, value)
		_, _ = writer.Connection.Write([]byte(element))
	}
	_, err := writer.Connection.Write([]byte("\r\n"))
	return err

}

func (wr *Writer) WriteBody(p []byte) (int, error) {
	return wr.Connection.Write(p)
}
