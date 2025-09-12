package main

import (
	"fmt"
	"io"
	"log"
	"net"
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

		linesChan := getLinesChannel(conn)

		for line := range linesChan {
			fmt.Println(line)
		}
		fmt.Println("Connection to ", conn.RemoteAddr(), "closed")
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	lines := make(chan string)
	go func() {
		defer f.Close()
		defer close(lines)
		currentLineContents := ""
		for {
			flag := false
			b := make([]byte, 8)
			n, err := f.Read(b)
			if err != nil {
				if currentLineContents != "" {
					currentLineContents += string(b[:n])
					lines <- currentLineContents
					return
				}
			}
			for index, element := range b {
				if element == '\n' {
					currentLineContents += string(b[:index])
					lines <- currentLineContents
					currentLineContents = ""
					currentLineContents += string(b[index+1 : n])
					flag = true
					break
				}
			}
			if !flag {
				currentLineContents += string(b[:n])
			}
		}
	}()
	return lines
}
