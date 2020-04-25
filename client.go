// Copyright 2019 kdevb0x Ltd. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause license
// The full license text can be found in the LICENSE file.

package quikface

import (
	"crypto"
	"crypto/ed25519"
	"crypto/rand"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/pion/webrtc/v2"

	"github.com/kdevb0x/quikface/internal/signal"
)

type Client struct {
	// uuid
	ID          uint32
	DisplayName string
	// Client network address
	Addr string
	// for signal mesages; known as "rwc" in collider
	// TODO: Maybe change this to io.ReadWriteCloser and have signalFunc
	// implement the interface?
	Signal io.ReadWriteCloser // signalFunc
	// MsgQueue MsgQueue

	// PrivateKey is a ed25519 PrivateKey for signing and authentication.
	PrivateKey crypto.PrivateKey

	// PublicKey is a ed25519 PublicKey for signing and authentication.
	PublicKey crypto.PublicKey
}

type MsgQueue struct {
	Inbound  chan Message // inbound messages from client
	Outbound chan Message
}

// NewClient represents the client side of connection during api interactions.
func NewClient(displayname ...string) *Client {
	var c *Client
	c.ID = uuid.New().ID()

	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		// TODO: make this a non-lethal error so it only terminates the
		// connection and not the whole server.
		panic(err)
	}
	c.PublicKey = pub
	c.PrivateKey = priv

	// var msgq = MsgQueue{make(chan Message, 10), make(chan Message, 10)}
	if len(displayname) > 0 {
		// cant break if stmnts, so had to flip this test
		if displayname[0] != "" {
			c.DisplayName = displayname[0]
		}
		// here displayname[0] == ""
	}
	c.DisplayName = randomChatName()
	return c

}

func (c *Client) Register(signaler io.ReadWriteCloser) error {
	if c.Signal != nil {
		return fmt.Errorf("duplicate registration; %s already has a signal connection registered", c.ID)
	}
	c.Signal = signaler
	return nil
}

// JoinRoom return a pointer to the room if the join was successfull, otherwise
// said pointer is nil, and err explains why.
func (c *Client) JoinRoom(name string, masterDirectory *RoomList) (*Room, error) {
	if hash, exists := masterDirectory.Rooms[name]; exists {
		room, err := masterDirectory.GetRoom(hash)
		if err != nil {
			return nil, fmt.Errorf("error finding hash of %s, %w\n", name, err)
		}
		room.ClientList[c.ID] = c
		if len(room.OfferQueue) > 0 {
			go func(r *Room) {
				c.initRTCSession(room)
			}(room)
		}
		return room, nil

	}
	return nil, fmt.Errorf("error: room %s doesn't exist", name)
}

func (c *Client) initRTCSession(r *http.Request, room *Room, recvr ...*Client) error {

	// Prepare the configuration
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
		SDPSemantics: webrtc.SDPSemanticsUnifiedPlanWithFallback,
	}

	// Create a MediaEngine object to configure the supported codec
	m := webrtc.MediaEngine{}

	// if browser is saffari, it only supports H264

	m.RegisterCodec(webrtc.NewRTPVP8Codec(webrtc.DefaultPayloadTypeVP8, 90000))
	m.RegisterCodec(webrtc.NewRTPOpusCodec(webrtc.DefaultPayloadTypeOpus, 48000))

	// Create the API object with the MediaEngine
	api := webrtc.NewAPI(webrtc.WithMediaEngine(m))

	// Create a new RTCPeerConnection
	peerConnection, err := api.NewPeerConnection(config)
	if err != nil {
		return err
	}

	_, err = peerConnection.AddTransceiver(webrtc.RTPCodecTypeAudio)
	if err != nil {
		return err
	}

	_, err = peerConnection.AddTransceiver(webrtc.RTPCodecTypeVideo)
	if err != nil {
		return err
	}
	cert, err := webrtc.GenerateCertificate(c.PrivateKey)
	if err != nil {
		return err
	}
	peerConnection.
		// Set a handler for when a new remote track starts, this handler copies inbound RTP packets,
		// replaces the SSRC and sends them back
		peerConnection.OnTrack(func(track *webrtc.Track, receiver *webrtc.RTPReceiver) {

		fmt.Printf("Track has started, of type %d: %s \n", track.PayloadType(), track.Codec().Name)
		for {
			// Read RTP packets being sent to Pion
			rtp, readErr := track.ReadRTP()
			if readErr != nil {
				if readErr == io.EOF {
					return
				}
				throwError(err)
			}
			switch track.Kind() {
			case webrtc.RTPCodecTypeAudio:
				saver.PushOpus(rtp)
			case webrtc.RTPCodecTypeVideo:
				saver.PushVP8(rtp)
			}
		}
	})
	// Set the handler for ICE connection state
	// This will notify you when the peer has connected/disconnected
	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		fmt.Printf("Connection State has changed %s \n", connectionState.String())
	})

	if len(room.OfferQueue) == 0 {
		// we are the first to initiate session

		// create out offer
		offer, err := peerConnection.CreateOffer(nil)
		if err != nil {
			return err
		}

		// set the local description
		err = peerConnection.SetLocalDescription(offer)
		if err != nil {
			return err
		}
		signal.Deco

	}
	remoteOffer := <-room.OfferQueue
}
