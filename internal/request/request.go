package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/JKAravind/TCPtoHTTP/internal/headers"
)

type requestState int

const (
	requestStateStart requestState = iota
	requestStateRequestLine
	requestStateHeaders
	requestStateBody
	requestStateDone
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
	Header      headers.Header
	state       requestState
	Body        []byte
}

var ErrInvalidReq = errors.New("Invaild")

func RequestFromReader(reader io.Reader) (*Request, error) {
	buffer := make([]byte, 8)
	readFromIndex := 0
	readFromParsed := 0

	req := &Request{
		state:  requestStateStart,
		Header: headers.NewHeaders(),
	}

	for req.state != requestStateDone {
		if readFromIndex == len(buffer) {
			temp := buffer
			buffer = make([]byte, readFromIndex*2)
			copy(buffer[:readFromIndex], temp)
		}

		n, err := reader.Read(buffer[readFromIndex:])
		if err != nil {
			fmt.Println(err)
			return req, err
		}

		temp, err := req.parse(buffer, readFromParsed, readFromIndex)
		if err != nil {
			fmt.Println(err)
			return req, err
		}
		readFromParsed += temp
		readFromIndex += n
	}
	return req, nil
}

func (r *Request) parse(data []byte, startToParseFrom, startToReadFrom int) (int, error) {
	fmt.Printf("data %s", string(data))
	totalBytesParsed := 0

	switch r.state {

	case requestStateStart:
		r.state = requestStateRequestLine

	case requestStateRequestLine:
		index := bytes.Index(data, []byte("\r\n"))
		if index == -1 {
			return 0, nil
		}
		r.state = requestStateHeaders
		parsedParts := bytes.Split(data[:index], []byte(" "))
		if len(parsedParts) != 3 {
			return 0, ErrInvalidReq
		}
		r.RequestLine.Method = string(parsedParts[0])
		r.RequestLine.RequestTarget = string(parsedParts[1])
		r.RequestLine.HttpVersion = string(parsedParts[2])
		totalBytesParsed += index + 2
		return index + 2, nil

	case requestStateHeaders:
		consumed, done, err := r.Header.Parse(data[startToParseFrom:])
		if err != nil {
			return 0, err
		}
		if done {
			r.state = requestStateBody
		}
		totalBytesParsed += consumed + 2
		return consumed, nil

	case requestStateBody:

		bodyLength, ok := r.Header.Get("content-length")
		if !ok {
			r.state = requestStateDone
			break
		}
		intBodyLength, _ := strconv.Atoi(bodyLength)

		remaining := intBodyLength - len(r.Body)
		toRead := len(data[startToParseFrom:])
		if toRead > remaining {
			toRead = remaining
		}

		r.Body = data[startToParseFrom : startToReadFrom+toRead]

		if len(r.Body) == intBodyLength {
			r.state = requestStateDone
			return intBodyLength, nil
		}
		return 0, nil

	}
	return 0, nil
}
