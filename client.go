// Copyright 2019 kdevb0x Ltd. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause license
// The full license text can be found in the LICENSE file.

package vidchat

import (
	"errors"
	"io"
	"net"
)

var (
	ErrNilConnError = errors.New("error cannot close nil connection")
)

// Client is a single participant dialing into chat other client.
type Client struct {
	Number string `json:"contact_number"` // telephone number
	Addr   net.Addr
	video  VideoStreamer
}

func (c *Client) Dial(address string) (net.Conn, error) {

}

type clientStream struct {
	local  Stream // own stream
	remote Stream // other callers stream
}

type stream struct {
	outgoing  io.Writer
	incomming io.Reader
}

type Stream interface {
	io.ReadWriteCloser
	// Open(ctx context.Context, transport *http.Transport) error
}
