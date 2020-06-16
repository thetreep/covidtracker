/*
	This file is part of covidtracker also known as EviteCovid .

    Copyright 2020 the Treep

    covdtracker is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    covidtracker is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with covidtracker.  If not, see <https://www.gnu.org/licenses/>.
*/

package http

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/thetreep/covidtracker"
	log "github.com/thetreep/covidtracker/logger"
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

	logger covidtracker.Logfer
}

// NewServer returns a new instance of Server.
func NewServer() *Server {
	return &Server{
		Addr:    DefaultAddr,
		Routing: make(map[string]http.Handler),
		logger:  &log.Logger{},
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
	s.logger.Debug(context.Background(), "starting operation api-server listening on %q", s.Addr)

	// Start HTTP server.
	go http.Serve(s.ln, s.Handlers())
	return nil
}

func (s *Server) Handlers() http.Handler {
	return adapt(s.Handler, s.ping(), s.cors(), s.log(), s.auth(), s.routing())
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
