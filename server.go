// Copyright 2019 kdevb0x Ltd. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause license
// The full license text can be found in the LICENSE file.

package quikface

import (
	"errors"
	"io"
	"net"
	"net/http"
)

var (
	ErrNilConnError = errors.New("error cannot close nil connection")
)

// localhost for dev. FIXME
var addr = "127.0.0.1:8000"

type Server struct {
	http.Server                        // net.Conn
	Addr        string                 // host address
	ActiveConns map[string]*ClientConn // key == remote address
}

type ClientConn struct {
	net.Conn
	VideoStream VideoStreamer
	Addr        net.Addr
	outbound    io.WriteCloser // to server
	inbound     io.ReadCloser  // from server
}

func (c *ClientConn) Dial(address string) (net.Conn, error) {
	return net.Dial("tcp", address)
}

// Read reads up to len(p) bytes from the ClientConn into p, returning the length
// of written data (n <= len(p)), and any errors encountered.
// Implements io.Reader interface.
func (r *ClientConn) Read(p []byte) (n int, err error) {
	return r.inbound.Read(p)
}

// Write writes p to the ClientConn, returning the length written, and any errors encountered.
// Implements io.Writer interface.
func (r *ClientConn) Write(p []byte) (n int, err error) {
	return r.outbound.Write(p)
}

// Close closes the ClientConns and returns any errors encountered.
func (r *ClientConn) Close() error {
	if err := r.inbound.Close(); err != nil {
		return err
	}
	if err := r.outbound.Close(); err != nil {
		return err
	}
	return nil
}

func NewServer() *Server {
	return &Server{
		Addr:        addr,
		ActiveConns: make(map[string]*ClientConn),
	}
}
func (s *Server) Accept() (net.Conn, error) {
	l, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return nil, err
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			continue
		}
		raddr := conn.RemoteAddr().String()
		s.ActiveConns[raddr] = &ClientConn{}
	}
}

func (s *Server) Close() error {
	if s.Conn == nil {
		return ErrNilConnError
	}
	return s.Conn.Close()
}

func (s *Server) Addr() net.Addr {
	if s.Conn != nil {
		return s.Conn.LocalAddr()
	}
	return nil
}
