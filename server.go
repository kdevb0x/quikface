// Copyright 2019 kdevb0x Ltd. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause license
// The full license text can be found in the LICENSE file.

package vidchat

import (
	"net"
)

// localhost for dev. FIXME
var addr = "127.0.0.1:8000"

type Server struct {
	HTTPServer                         // net.Conn
	Addr        string                 // host address
	ActiveConns map[string]*clientConn // key == remote address
}

type clientConn struct {
	net.Conn
	videoStream Stream // like *stream
}

func NewServer() *Server {
	return &Server{
		addr:        addr,
		activeConns: make(map[string]*clientConn),
	}
}

func (s *Server) Accept() (net.Conn, error) {
	if s.Conn != nil {
		return s.Conn, nil
	}
	l, err := net.Listen("tcp", s.addr)
	if err != nil {
		return nil, err
	}
	return l.Accept()

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
