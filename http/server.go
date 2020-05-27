package http

import (
	"fmt"
	"log"
	"net"
	"net/http"
)

// DefaultAddr is the default bind address.
const DefaultAddr = ":3456"

// Server represents an HTTP server.
type Server struct {
	ln net.Listener

	Handler http.Handler

	// Handlers to serve by map
	Routing map[string]http.Handler

	// Bind address to open.
	Addr string
}

// NewServer returns a new instance of Server.
func NewServer() *Server {
	return &Server{
		Addr:    DefaultAddr,
		Routing: make(map[string]http.Handler),
	}
}

func (s *Server) AddHandler(h http.Handler, path string) error {
	if _, exist := s.Routing[path]; exist {
		return fmt.Errorf("route %q has already an handler", path)
	}
	s.Routing[path] = h
	return nil
}

// Open opens a socket and serves the HTTP server.
func (s *Server) Open() error {
	// Open socket.
	ln, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}
	s.ln = ln
	log.Printf("starting operation api-server listening on %q", s.Addr)

	// Start HTTP server.
	go http.Serve(s.ln, adapt(s.Handler, s.routing(), s.cors(), s.log(), s.auth(), s.ping()))
	return nil
}

// Close closes the socket.
func (s *Server) Close() error {
	if s.ln != nil {
		s.ln.Close()
	}
	return nil
}

// Port returns the port that the server is open on. Only valid after open.
func (s *Server) Port() int {
	return s.ln.Addr().(*net.TCPAddr).Port
}
