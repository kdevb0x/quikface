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

type Stream interface {
	io.ReadWriteCloser
	// Open(ctx context.Context, transport *http.Transport) error
}

type stream struct {
	outbound io.WriteCloser // to server
	inbound  io.ReadCloser  // from server
}

// Read reads up to len(p) bytes from the stream into p, returning the length
// of written data (n <= len(p)), and any errors encountered.
// Implements io.Reader interface.
func (r *stream) Read(p []byte) (n int, err error) {
	return r.inbound.Read(p)
}

// Write writes p to the stream, returning the length written, and any errors encountered.
// Implements io.Writer interface.
func (r *stream) Write(p []byte) (n int, err error) {
	return r.outbound.Write(p)
}

// Close closes the streams and returns any errors encountered.
func (r *stream) Close() error {
	if err := r.inbound.Close(); err != nil {
		return err
	}
	if err := r.outbound.Close(); err != nil {
		return err
	}
	return nil
}
