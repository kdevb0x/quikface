// Copyright 2019 kdevb0x Ltd. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause license
// The full license text can be found in the LICENSE file.

package quikface

import (
	"errors"
	"context"
	"io"
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gobwas/ws"
)

var (
	ErrNilConnError = errors.New("error cannot close nil connection")
)

// localhost for dev. FIXME
var addr = "127.0.0.1:8000"

type Server struct {
	httpserver  *http.Server        // net.Conn
	Addr       string              // host address
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
		Addr:       addr,
		ActiveConns: make(map[string]net.Conn),
	}
}

func (s *Server) Accept() (net.Conn, error) {
	l, err := net.Listen("udp", s.Addr)
	if err != nil {
		return nil, err
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			continue
		}
		raddr := conn.RemoteAddr().String()
		hs. err := ws.Upgrade(conn)
		if err != nil {
			return nil, err
		}

		s.ActiveConns[raddr] = conn
	}
	// BUG: I'm not sure why this works without a return here.
	// s is a listener, so s.Accept should overwrite the promoted method,

}

func (s *Server) Close() error {
	return s.httpserver.Close()
}

func (s *Server) Addr() net.Addr {
	addr, _ := net.ResolveUDPAddr("udp", s.Addr)
	return addr
}

// Tmp catch all paths connection initiator
// TODO: extract out logic and make this pretty.
func (s *Server) ListenAndServeTLS() error {
	r := mux.NewRouter()
	r.Methods("POST").Path("/createAccount").HandlerFunc(CreateAccountHandler)
	s.httpserver.Handler = r
	return http.ServeTLS(s, r, "localhost+2.pem", "cert/localhost+2-key.pem")
}

func (s *Server) ClientFromReq(req *http.Request) *Client {
	ip := req.RemoteAddr
	c := &Client{

	}
}
// genClientSessionKey generates a cryptographic session key for use in calls,
// and stores it in ctx.Value.
func (s *Server) genClientSessionKey(ctx context.Context, c *Client) {

}
