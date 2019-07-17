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
	http.Server                     // net.Conn
	SAddr       string              // host address
	ActiveConns map[string]net.Conn // key == remote address
}

type Client struct {
	Conn        net.Conn
	VideoStream VideoStreamer
	Addr        net.Addr

	// FIXME: inbound and outbound together are basically just a net.Conn,
	// so extact this into a type that implements the interface.
	outbound io.WriteCloser // to server
	inbound  io.ReadCloser  // from server
}

func (c *Client) Dial(address string) (net.Conn, error) {
	return net.Dial("tcp", address)
}

// Read reads up to len(p) bytes from the Client into p, returning the length
// of written data (n <= len(p)), and any errors encountered.
// Implements io.Reader interface.
func (r *Client) Read(p []byte) (n int, err error) {
	return r.inbound.Read(p)
}

// Write writes p to the Client, returning the length written, and any errors encountered.
// Implements io.Writer interface.
func (r *Client) Write(p []byte) (n int, err error) {
	return r.outbound.Write(p)
}

// Close closes the Clients and returns any errors encountered.
func (r *Client) Close() error {
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
		SAddr:       addr,
		ActiveConns: make(map[string]net.Conn),
	}
}

func (s *Server) Accept() (net.Conn, error) {
	l, err := net.Listen("udp", s.SAddr)
	if err != nil {
		return nil, err
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			continue
		}
		raddr := conn.RemoteAddr().String()
		s.ActiveConns[raddr] = conn
	}
}

func (s *Server) Close() error {
	return s.Server.Close()
}

func (s *Server) Addr() net.Addr {
	addr, _ := net.ResolveUDPAddr("udp", s.SAddr)
	return addr
}

func (s *Server) ListenAndServeTLS() {
	go s.ServeTLS(s, "localhost+2.pem", "cert/localhost+2-key.pem")
}
