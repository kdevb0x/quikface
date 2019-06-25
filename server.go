// Copyright 2019 kdevb0x Ltd. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause license
// The full license text can be found in the LICENSE file.

package vidchat

import (
	"io"
	"net"
)

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

func (t *stream) Read(p []byte) (n int, err error) {
	return io.ReadFull(t.incomming, p)
}

func (t *stream) Write(p []byte) (n int, err error) {
	return t.outgoing.Write(p)
}
