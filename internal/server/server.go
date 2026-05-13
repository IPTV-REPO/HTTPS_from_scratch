package server

import (
	"fmt"
	"net"
	"sync/atomic"
)

type Server struct {
    closed   atomic.Bool
    listener net.Listener
}

// Serve initializes the listener and starts the background loop.
func Serve(port int) (*Server, error) {
    l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
    if err != nil {
        return nil, err
    }

    srv := &Server{
        listener: l,
    }

    // Start the background listener
    go srv.listen()

    return srv, nil
}

// Close triggers a graceful shutdown by flagging the state and killing the listener.
func (s *Server) Close() error {
    s.closed.Store(true)
    return s.listener.Close()
}

// listen handles the accept loop.
func (s *Server) listen() {
    // Ensure the listener is cleaned up when this goroutine eventually exits.
    defer s.listener.Close()

    for {
        conn, err := s.listener.Accept()
        if err != nil {
            // If the error happened because we intentionally closed the server, 
            // exit the loop quietly.
            if s.closed.Load() {
                break
            }
            // Otherwise, it's a real network error; log it and keep trying.
            fmt.Println("Error accepting connection:", err)
            continue
        }

        // Handle the connection in a separate goroutine.
        go s.handle(conn)
    }
}

// handle writes the specific HTTP response required and closes the connection.
func (s *Server) handle(conn net.Conn) {
    defer conn.Close()

    // The raw HTTP response with the mandatory double CRLF (\r\n\r\n)
    response := "HTTP/1.1 200 OK\r\n" +
        "Content-Type: text/plain\r\n" +
        "Content-Length: 13\r\n" +
        "\r\n" +
        "Hello World!\n"

    _, err := conn.Write([]byte(response))
    if err != nil {
        fmt.Println("Error writing response:", err)
    }
}