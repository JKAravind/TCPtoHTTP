package main

import (
	"fmt"
	"log"
	"net"

	"github.com/JKAravind/TCPtoHTTP/internal/request"
)

const port = ":42069"

func main() {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("error listening for TCP traffic: %s\n", err.Error())
	}
	defer listener.Close()

	fmt.Println("Listening for TCP traffic on", port)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("error: %s\n", err.Error())
		}
		fmt.Println("Accepted connection from", conn.RemoteAddr())

		httpHeader, err := request.RequestFromReader(conn)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Request line:")
		fmt.Printf("- Method: %s\n", httpHeader.RequestLine.Method)
		fmt.Printf("- Target: %s\n", httpHeader.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", httpHeader.RequestLine.HttpVersion)

		fmt.Println("Headers:")
		for k, v := range httpHeader.Header {
			fmt.Printf("- %s: %s\n", k, v)
		}

		fmt.Println("Body:")
		fmt.Println(string(httpHeader.Body))

		fmt.Println("Connection to ", conn.RemoteAddr(), "closed")
	}
}
