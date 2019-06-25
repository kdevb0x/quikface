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
	host   Stream
	remote Stream
}

type stream struct {
	outgoing  io.Writer
	incomming io.Reader
}

type Stream interface {
	io.ReadWriteCloser
	Open(ctx context.Context, transport *http.Transport) error
}
