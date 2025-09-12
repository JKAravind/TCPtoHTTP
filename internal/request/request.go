package request

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
}

var ErrInvalidReq = errors.New("Invaild")

func RequestFromReader(reader io.Reader) (*Request, error) {
	ReadArray, err := io.ReadAll(reader)
	if err != nil {
		fmt.Println(err)
	}
	reqLineDict := parseRequestLine(ReadArray)

	if reqLineDict["method"] != "GET" && reqLineDict["method"] != "POST" {
		fmt.Println(err)
		return nil, ErrInvalidReq
	}

	reqLine := RequestLine{
		HttpVersion:   reqLineDict["httpVersion"],
		RequestTarget: reqLineDict["endpoint"],
		Method:        reqLineDict["method"],
	}
	Request := Request{reqLine}
	fmt.Println(Request.RequestLine)
	return &Request, nil

}

func parseRequestLine(reqLine []byte) map[string]string {

	dict := make(map[string]string, 0)

	splitArr := strings.Split(string(reqLine), "\r\n")[0]
	reqLineArr := strings.Split(string(splitArr), " ")
	dict["method"] = reqLineArr[0]
	dict["endpoint"] = reqLineArr[1]
	dict["httpVersion"] = strings.Split(reqLineArr[2], "/")[1]
	fmt.Println(dict)
	return dict

}
