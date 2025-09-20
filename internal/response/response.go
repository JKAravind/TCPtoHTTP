package response

import (
	"fmt"
	"io"
	"strconv"

	"github.com/JKAravind/TCPtoHTTP/internal/headers"
)

type StatusCode int

const (
	statusSuccess    StatusCode = 200
	statusWrong      StatusCode = 400
	statusServerFail StatusCode = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	var write []byte = []byte{}
	switch statusCode {
	case statusSuccess:
		fmt.Println("writing this to the rcurl client")
		write = []byte("HTTP/1.1 200 OK\r\n")
	case statusWrong:
		write = []byte("HTTP/1.1 400 Bad Request\r\n")
	case statusServerFail:
		write = []byte("HTTP/1.1 500 Internal Server Error\r\n")
	}
	n, err := w.Write(write)
	fmt.Println("wrote bytes:", n, "err:", err)

	return err
}

func GetDefaultHeaders(contentLen int) headers.Header {
	responseHeader := headers.NewHeaders()
	responseHeader["Content-Length"] = strconv.Itoa(contentLen)
	responseHeader["Connection"] = "close"

	responseHeader["Content-Type"] = "text/plain"

	return responseHeader
}
func WriteHeaders(w io.Writer, headers headers.Header) error {

	for key, value := range headers {
		element := fmt.Sprintf("%s: %s\r\n", key, value)
		_, _ = w.Write([]byte(element))
	}
	_, err := w.Write([]byte("\r\n"))
	return err

}
