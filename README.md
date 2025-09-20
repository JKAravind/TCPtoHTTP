# TCP → HTTP Server in Go

This project is a **learning exploration** of how HTTP works under the hood.
Instead of using `net/http` or frameworks like Express, we build an HTTP/1.1 server **from scratch**, starting from raw TCP connections.

## 🧩 What It Does

* Accepts TCP connections on a given port
* Parses HTTP request lines, headers, and body manually
* Implements a `response.Writer` to write status lines, headers, and body
* Handles simple routing:

  * `/yourproblem` → 400 Bad Request
  * `/myproblem` → 500 Internal Server Error
  * everything else → 200 OK

## ⚡ Learning Outcomes

* Understand the **low-level mechanics of HTTP**:

  * How TCP streams become HTTP requests
  * Why headers and body are separated by `\r\n`
  * Why `Content-Length` or chunked encoding is needed
* See the difference between **high-level frameworks** (Express, net/http) and **raw protocol handling**
* Learn how HTTP status codes and responses are structured

## 🚀 Run It

```bash
go run ./cmd
```

Test with curl:

```bash
curl http://localhost:42069
curl http://localhost:42069/yourproblem
curl http://localhost:42069/myproblem
```

## 📖 Why Build This?

High-level frameworks are convenient, but they **abstract away the real engine**.
By writing your own server from scratch, you learn **how HTTP really works** — a foundational skill for backend engineers.
