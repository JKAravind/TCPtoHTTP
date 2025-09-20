package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/JKAravind/TCPtoHTTP/internal/request"
	"github.com/JKAravind/TCPtoHTTP/internal/response"
	"github.com/JKAravind/TCPtoHTTP/internal/server"
)

const port = 42069

var Html200 = `
<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>
`

// 400 Bad Request page
var Html400 = `
<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>
`

// 500 Internal Server Error page
var Html500 = `
<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>
`

func handler(w *response.Writer, req *request.Request) {
	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		w.WriteStatusLine(response.StatusWrong)
		w.WriteHeaders(response.GetDefaultHeaders(len(Html200)))
		w.WriteBody([]byte(Html200))
	case "/myproblem":
		w.WriteStatusLine(response.StatusServerFail)
		w.WriteHeaders(response.GetDefaultHeaders(len(Html400)))
		w.WriteBody([]byte(Html400))
	default:
		w.WriteStatusLine(response.StatusSuccess)
		w.WriteHeaders(response.GetDefaultHeaders(len(Html500)))
		w.WriteBody([]byte(Html500))
	}
}

func main() {
	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
