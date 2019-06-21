// Copyright 2019 kdevb0x Ltd. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause license
// The full license text can be found in the LICENSE file.

package vidchat

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
)

// localhost for dev. FIXME
var addr = "127.0.0.1:8000"

var (
	ErrNilConnError = errors.New("error cannot close nil connection")
)

type server struct {
	net.Conn
	addr        string                     // host address
	activeConns map[string]*session_stream // key == remote address
}

type session_stream struct {
	host   *stream
	remote *stream
}

type stream struct {
	outgoing  Stream
	incomming Stream
}

type Stream interface {
	io.ReadWriteCloser
	Open(ctx context.Context, transport *http.Transport) error
}

func NewServer() *server {
	return &server{
		addr:        addr,
		activeConns: make(map[string]*session_stream),
	}
}

func (s *server) Accept() (net.Conn, error) {
	if s.Conn != nil {
		return s.Conn, nil
	}
	l, err := net.Listen("tcp", s.addr)
	if err != nil {
		return nil, err
	}
	return l.Accept()

}

func (s *server) Close() error {
	if s.Conn == nil {
		return ErrNilConnError
	}
	return s.Conn.Close()
}

func (s *server) Addr() net.Addr {
	if s.Conn != nil {
		return s.Conn.LocalAddr()
	}
	return nil
}
