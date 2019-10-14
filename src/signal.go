// Copyright 2019 kdevb0x Ltd. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause license
// The full license text can be found in the LICENSE file.

// Some ideas for signaling in this package have been pulled from collider:
// Copyright (c) 2014 The WebRTC project authors.
// https://github.com/webrtc/apprtc/blob/master/src/collider
// and
// https://github.com/dennwc/dom/blob/master/net/webrtc/signalling.go

package quikface

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	gws "github.com/gorilla/websocket"
	"golang.org/x/net/websocket"
)

type signalFunc func(ws *websocket.Conn, client *Client, msg Message) error

type Signal struct {
	UserId string // optional user id of client
	Data   []byte // webrtc SDP payload
}

type Signaller interface {
	Broadcast(s Signal) (AnswerStream, error)
	Listen(uid string) (OfferStream, error)
}

type AnswerStream interface {
	Next() (Signal, error)
	Close() error
}

type Offer interface {
	Info() Signal
	Answer(s Signal) error
}

type OfferStream interface {
	Next() (Offer, error)
	Close() error
}

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
	Type     string `json:"type"`
	SDP      string `json:"sdp,omitempty"`
	Name     string `json:"name,omitempty"`
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

// Allows compressing offer/answer to bypass terminal input limits.
const compress = false

// MustReadStdin blocks until input is received from stdin, and panics on error.
func MustReadStdin() string {
	r := bufio.NewReader(os.Stdin)

	var in string
	for {
		var err error
		in, err = r.ReadString('\n')
		if err != io.EOF {
			if err != nil {
				panic(err)
			}
		}
		in = strings.TrimSpace(in)
		if len(in) > 0 {
			break
		}
	}

	return in
}

// Encode encodes the input in base64
// It can optionally zip the input before encoding
func encode(obj interface{}) string {
	b, err := json.Marshal(obj)
	if err != nil {
		return ""
	}

	if compress {
		b, err = zip(b)
		if err != nil {
			return ""
		}
	}

	return base64.StdEncoding.EncodeToString(b)
}

// Decode decodes the input from base64
// It can optionally unzip the input after decoding
func decode(in string, obj interface{}) error {
	b, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		return err
	}

	if compress {
		b, err = unzip(b)
		if err != nil {
			return err
		}
	}

	err = json.Unmarshal(b, obj)
	if err != nil {
		return err
	}
	return nil
}

func zip(in []byte) ([]byte, error) {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	_, err := gz.Write(in)
	if err != nil {
		return nil, err
	}
	err = gz.Flush()
	if err != nil {
		return nil, err
	}
	err = gz.Close()
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func unzip(in []byte) ([]byte, error) {
	var b bytes.Buffer
	_, err := b.Write(in)
	if err != nil {
		return nil, err
	}
	r, err := gzip.NewReader(&b)
	if err != nil {
		return nil, err
	}
	res, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func WSHandler(w http.ResponseWriter, r *http.Request) {
	upgrader := gws.Upgrader{}
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "unable to upgrade ws conn", http.StatusInternalServerError)
	}
	defer func() {
		errorLog <- c.Close()
	}()
	for {
		mt, msg, err := c.ReadMessage()
		if err != nil {
			// TODO: handle error
			errorLog <- err
		}
		wsData := wsClientMsg{}
		if err := json.Unmarshal(msg, &wsData); err != nil {
			errorLog <- err
		}
		sdp := wsData.SDP
		name := wsData.Name
		if wsData.Type == "publish" {

		}
	}
}
