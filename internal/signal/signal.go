// Copyright 2019 kdevb0x Ltd. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause license
// The full license text can be found in the LICENSE file.
//
// Some ideas for signaling in this package have been pulled from collider:
// Copyright (c) 2014 The WebRTC project authors.
// https://github.com/webrtc/apprtc/blob/master/src/collider
// and
// https://github.com/dennwc/dom/blob/master/net/webrtc/signalling.go

// signal is an internal pkg of quikface that implements simple signal server.
package signal

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	gws "github.com/gorilla/websocket"
	"golang.org/x/net/websocket"

	qf "github.com/kdevb0x/quikface"
)

// Allows compressing offer/answer to bypass terminal input limits.
const compress = false

type signalFunc func(ws *websocket.Conn, client *qf.Client, msg Message) error

type Signal struct {
	UserId string // optional user id of client
	Data   []byte // webrtc SDP payload
}

/*
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
*/

type Message interface {
	Cmd() string // returns the signal command as string
	Send(w io.Writer, msg string) error
	String() string // returns the message string
}

// websocket message from the client
type WsClientMessage struct {
	Command  string `json:"cmd"`
	ClientID string `json:"client_id"`
	Msg      string `json:"msg"`
	Type     string `json:"type"`
	SDP      string `json:"sdp,omitempty"`
	Name     string `json:"name,omitempty"`
}

func (cm WsClientMessage) Cmd() string {
	return cm.Command
}

func (cm WsClientMessage) Send(w io.Writer, msg string) error {
	cm.Msg = msg
	return send(w, cm)
}

func (cm WsClientMessage) String() string {
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

func (sig *Signal) Cmd() string {
	if len(sig.UserId) > 0 {
		return sig.UserId
	}
	return ""
}

func (sig *Signal) Send(w io.Writer, msg string) error {
	if msg != "" {
		jmsg, err := json.Marshal(msg)
		if err != nil {
			return err
		}
		sig.Data = append(sig.Data, jmsg...)
		if err := json.NewEncoder(w).Encode(sig.Data); err != nil {
			return err
		}
		return nil
	}
	return errors.New(`error: msg == ""; Unable to send blank message`)
}

func (sig *Signal) String() string {
	var b bytes.Buffer
	if err := json.Unmarshal(sig.Data, &b); err != nil {
		return ""
	}
	return b.String()
}

func Send(w io.Writer, data interface{}) error {
	e := json.NewEncoder(w)
	if err := e.Encode(data); err != nil {
		return err
	}
	return nil
}

// MustReadStdin blocks until input is received from stdin, and panics on error.
func MustReadStdin() string {
	r := bufio.NewReader(os.Stdin)

	var in string
	for {
		var err error
		in, err = r.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil && err != io.EOF {
			return err.Error()
		}
		if len(in) == 0 {
			continue
		}
	}
	in = strings.TrimSpace(in)
	return in
}

// Encode encodes the input in base64
// It can optionally zip the input before encoding
func Encode(obj interface{}, compress bool) string {
	b, err := json.Marshal(obj)
	if err != nil {
		return ""
	}

	if compress {
		b, err = Zip(b)
		if err != nil {
			return err.Error()
		}
	}

	return base64.StdEncoding.EncodeToString(b)
}

// Decode decodes the input from base64
// It can optionally unzip the input after decoding
func Decode(in string, obj interface{}, decompress bool) error {
	b, err := base64.StdEncoding.DecodeString(in)
	if err != nil {
		return err
	}

	if decompress {
		b, err = Unzip(b)
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

func Zip(in []byte) ([]byte, error) {
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

func Unzip(in []byte) ([]byte, error) {
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
		quikface.ErrorLog <- c.Close()
	}()
	for {
		mt, msg, err := c.ReadMessage()
		if err != nil {
			// TODO: handle error

			// errorLog does log.Print* for all errors sent
			errorLog <- err
		}
		wsData := WsClientMessage{}
		if err := json.Unmarshal(msg, &wsData); err != nil {
			err = fmt.Errorf("failed to marshal ws message from signaler: %w\n", err)
			errorLog <- err
		}
		sdp := wsData.SDP
		name := wsData.Name
		if wsData.Type == "publish" {

		}
	}

}
