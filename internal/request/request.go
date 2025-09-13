package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
	state       string
}

var ErrInvalidReq = errors.New("Invaild")

func RequestFromReader(reader io.Reader) (*Request, error) {
	buffer := make([]byte, 8)
	readFromIndex := 0
	readFromParsed := 0

	req := &Request{
		state: "initialised",
	}

	for req.state == "initialised" {
		if readFromIndex == len(buffer) {
			temp := buffer
			buffer = make([]byte, readFromIndex*2)
			copy(buffer[:readFromIndex], temp)
		}

		n, err := reader.Read(buffer[readFromIndex:])
		if err != nil {
			fmt.Println(err)
		}
		temp, err := req.parse(buffer)
		if err != nil {
			fmt.Println(err)
			return req, err
		}
		readFromParsed += temp
		readFromIndex += n

	}

	return req, nil

}

func (r *Request) parse(data []byte) (int, error) {

	index := bytes.Index(data, []byte("\r\n"))
	if index == -1 {
		return 0, nil
	}
	r.state = "done"
	parsedParts := bytes.Split(data[:index], []byte(" "))
	if len(parsedParts) != 3 {
		return 0, ErrInvalidReq
	}
	r.RequestLine.Method = string(parsedParts[0])
	r.RequestLine.RequestTarget = string(parsedParts[1])
	r.RequestLine.HttpVersion = string(parsedParts[2])

	return index, nil

}
