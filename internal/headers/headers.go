package headers

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

type Header map[string]string

var ErrInvalidHeader = errors.New("invalid header")

func NewHeaders() Header {
	return make(Header)
}

func (header Header) Parse(data []byte) (n int, done bool, err error) {
	// fmt.Println("what is coming", string(data))
	readIndex := 0
	isDone := false

	for {
		index := bytes.Index(data[readIndex:], []byte("\r\n"))
		if index == -1 {
			break
		}
		if index == 0 {
			// this indicates that the \r\n\r\n is twice and it retuns
			isDone = true
			readIndex += 2
			break
		}
		name, value, err := ParseHeader(data[readIndex : readIndex+index])
		if err != nil {
			return 0, false, err
		}
		name = strings.ToLower(name)
		_, ok := header[name]
		if ok {
			header[name] = fmt.Sprintf("%s , %s", header[name], value)
		} else {
			header[name] = value

		}
		readIndex += index + 2

	}
	return readIndex, isDone, nil

}

func ParseHeader(data []byte) (string, string, error) {
	newSlice := bytes.SplitN(bytes.TrimSpace(data), []byte(":"), 2)
	if len(newSlice) != 2 {
		return "", "", ErrInvalidHeader
	}
	name, value := newSlice[0], newSlice[1]

	if bytes.HasSuffix(name, []byte(" ")) {
		return "", "", ErrInvalidHeader
	}

	return string(name), strings.TrimSpace(string(value)), nil

}
