// Copyright 2019 kdevb0x Ltd. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause license
// The full license text can be found in the LICENSE file.

// +build wasm,js

package quikface // import "github.com/kdevb0x/quikface"

import (
	"io"
	gojs "syscall/js"

	"github.com/dennwc/dom"
	"github.com/dennwc/dom/js"
	"github.com/dennwc/dom/net/webrtc"

	rtc "github.com/pion/webrtc/v2"
)

var (
	_ webrtc.Local
)

type MediaStream interface {
	// returns the instances unique 36 char string
	Id() string
	GetAudioTracks []MediaStreamTrack
	GetTrackById(id string) (MediaStreamTrack, error)
	GetVideoTracks() []MediaStreamTrack

}

// JSAudioTrack represents an audio track from the browsers MediaStream.
type JSAudioTrack struct {
	js.Value
}

// JSVideoTrack represents a video track from the browsers MediaStream.
type JSVideoTrack struct {
	js.Value
}

func (au *JSAudioTrack) Id() string {

}

// InitBrowserCam instantiates the webcam through the browsers
// navigator.MediaDevices.getUserMedia API.
func InitBrowserCam() (MediaStream, error) {
	var localVideo = dom.GetDocument().QuerySelector("localVideo")
	goconstraints := map[string]interface{}{"audio": true, "video": map[string]interface{}{"facingMode": "user"}}
	constraints := js.ValueOf(goconstraints)
	mediaDevices := gojs.Global().Get("navigator").Get("mediaDevices")

	getUserMediaPromise := mediaDevices.Call("getUserMedia", constraints)

	js.AsyncCallbackOf(func(streams []js.Value) {
		if len(streams) > 0 {
			for _, stream := range streams {
				var vidTrack = stream.Call("getVideoTracks")
			}
			localVideo.SetAttribute("srcObject", stream)
		}

	}
	// var mediaPromise = js.NewPromise(func()([]interface{}, error) {
	// 	return mediaDevices.Call("getUserMedia", constraints), nil
	// })

	var rtcConfig = rtc.Configuration{
		ICEServers: []rtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}
	peerconn, err := rtc.NewPeerConnection(rtcConfig)
	if err != nil {
		return nil, err
	}
	peerconn.
}
