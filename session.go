// Copyright 2019 kdevb0x Ltd. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause license
// The full license text can be found in the LICENSE file.

package quikface

import (
	"fmt"

	rtc "github.com/pion/webrtc/v2"
)

func NewDefaultPeerConnection() (*rtc.PeerConnection, error) {
	config := rtc.Configuration{
		ICEServers: []rtc.ICEServer{
			{
				URLs: []string{"stun:stun.1.google.com:19302"},
			},
		},
	}

	peerConn, err := rtc.NewPeerConnection(config)
	if err != nil {
		return nil, err
	}
	return peerConn, nil

}

type dataChannel struct {
	dc       *rtc.DataChannel
	inbound  chan rtc.DataChannelMessage // TODO: refine this, drop interface{}
	outbound chan rtc.DataChannelMessage
}

func NewDataChannel(label string, peerconn *rtc.PeerConnection) (*dataChannel, error) {
	dc, err := peerconn.CreateDataChannel(label, nil)
	if err != nil {
		return nil, err
	}
	peerconn.OnDataChannel(func(dataChannel *rtc.DataChannel) {
		fmt.Printf("New Data Channel %s %d \n", dataChannel.Label, dataChannel.ID())
	})

	dcs := &dataChannel{dc: dc}
	dcs.dc.OnOpen(func() {
		var ib = make(chan rtc.DataChannelMessage)
		var ob = make(chan rtc.DataChannelMessage)
		dcs.inbound = ib
		dcs.outbound = ob

	})
	dcs.dc.OnMessage(func(msg rtc.DataChannelMessage) {
		go func() {
			dcs.inbound <- msg
		}()
	})

	return dcs, nil
}

func (dc *dataChannel) SendTxt(msg string) error {
	return dc.dc.SendText(msg)
}
