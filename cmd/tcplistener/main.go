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
		fmt.Println("Request line: ")
		fmt.Printf("+%v", httpHeader.RequestLine)
		for key, value := range httpHeader.Header {
			fmt.Printf("%vkey: %v\n", key, value)
		}

		fmt.Println("Connection to ", conn.RemoteAddr(), "closed")
	}
}
