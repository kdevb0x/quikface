// Copyright 2019 kdevb0x Ltd. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause license
// The full license text can be found in the LICENSE file.

// Some ideas for signaling in this package have been pulled from collider:
// Copyright (c) 2014 The WebRTC project authors.
// https://github.com/webrtc/apprtc/blob/master/src/collider

package quikface

import (
	"encoding/json"
	"io"

	_ "github.com/gobwas/ws"
	_ "github.com/pion/webrtc/v2"
	"golang.org/x/net/websocket"
)

type signalFunc func(ws *websocket.Conn, client *Client, msg Message) error

type Message interface {
	Cmd() string // returns the signal command as string
	Send(w io.Writer, msg string) error
	String() string // returns the message string
}

// websocket message from the client
type wsClientMsg struct {
	Command  string `json:"cmd"`
	ClientID string `json:"client_id"`
	Msg      string `json:"msg"`
}

func (cm wsClientMsg) Cmd() string {
	return cm.Command
}

func (cm wsClientMsg) Send(w io.Writer, msg string) error {
	cm.Msg = msg
	return send(w, cm)
}

func (cm wsClientMsg) String() string {
	return cm.Msg
}

type wsServerMsg struct {
	Msg string `json:"msg"`
	Err string `json:"error"`
}

func (sm wsServerMsg) Cmd() string {
	return ""
}

func (sm *wsServerMsg) Send(w io.Writer, msg string) error {
	sm.Msg = msg
	return send(w, sm)
}

func (sm wsServerMsg) String() string {
	return sm.Msg
}

func send(w io.Writer, data interface{}) error {
	e := json.NewEncoder(w)
	if err := e.Encode(data); err != nil {
		return err
	}
	return nil
}
