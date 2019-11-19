// Copyright 2019 kdevb0x Ltd. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause license
// The full license text can be found in the LICENSE file.

package quikface

import (
	"crypto/hmac"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"golang.org/x/crypto/sha3"

	"github.com/pion/webrtc/v2"
	"github.com/pion/rtp"
	"github.com/pion/rtp/codecs"
	"github.com/pion/rtcp"
	"github.com/pion/sdp/v2"
)

// maxRoomCap is the maximum number of clients allowed to a room.
const maxRoomCap = 2

// NewRoomList creates the global room directory, which lists all Rooms in existence.
func NewRoomList() *RoomList {
	return &RoomList{
		Rooms:      make(map[string]string),
		roomhashes: make(map[string]*Room),
	}
}

type RoomList struct {
	mu sync.Mutex
	// Rooms is the global list of existing rooms, mapping their
	// human-readable names to thier id hashes, which can be used by
	// GetRoom() to get a pointer to the actual room.
	Rooms map[string]string

	roomhashes map[string]*Room
}

type Room struct {
	// directory points back to the master list of all existing rooms.
	directory *RoomList

	Name string

	// id is a hash to uniquely identify the Room.
	id string

	// ClientList is a map of the Clients present, keyed by their id's.
	ClientList map[uint32]*Client

	// Client registration time limit; ClientList must register before timeout.
	RegistrationTimeout time.Duration

	// url.URL.String() of the room
	RoomURL string

	// OfferQueue is a shared queue for client session descriptions
	// This queue takes the place of a signaling server for client msgs.
	OfferQueue chan sdp.SessionDescription // unbuffered
}

// GetRoom searches by hash in existing rooms, returning a pointer if it exists,
// and an error if not.
func (rl *RoomList) GetRoom(hash string) (*Room, error) {
	r, exists := rl.roomhashes[hash]
	if !exists {
		return nil, fmt.Errorf("error: this roomhash [%s] doesn't exist!", hash)
	}
	return r, nil
}

// NewRoom creates a new room, adding it to rl. Returns non-nil err if it
// already exists.
// If > 1 registrationTimeout is provided, all but the first one are ignored.
func (rl *RoomList) NewRoom(name string, url string, registrationTimeout ...time.Duration) (*Room, error) {
	if _, exists := rl.Rooms[name]; exists {
		return nil, fmt.Errorf("error: room name [%s] already exists!", name)
	}
	r := &Room{
		directory: rl,
		Name:      name,
		RoomURL:   url,
	}
	if len(registrationTimeout) > 0 {
		r.RegistrationTimeout = registrationTimeout[0]
	}

	rid, err := rl.HashRoom(r.Name)
	if err != nil {
		return nil, fmt.Errorf("can't create room named %s, [internal error]: %w\n", name, err)
	}
	r.ClientList = make(map[uint32]*Client)
	r.id = rid
	rl.Rooms[r.Name] = r.id
	rl.roomhashes[r.id] = r
	return r, nil
}

func (rl *RoomList) HashRoom(name string) (string, error) {
	var ut = time.Now().UnixNano()
	key := strconv.AppendInt(make([]byte, strconv.IntSize), ut, 2)

	mac := hmac.New(sha3.New512, key)
	_, err := mac.Write([]byte(name))
	if err != nil {
		if err != nil {
			return "", err
		}

	}
	return hex.EncodeToString(mac.Sum(nil)), nil
}

// Client returns the client if present, otherwise a new Client is created
// inside the Room r, (if maxRoomCap iis not met) then returned.
func (r *Room) Client(clientID uint32) (*Client, error) {
	if c, exists := r.ClientList[clientID]; exists {
		return c, nil
	}
	if len(r.ClientList) >= maxRoomCap {
		return nil, errors.New("Room at max capacity, cannot add client")

	}
	r.ClientList[clientID] = NewClient()
	return r.ClientList[clientID], nil
}

func (r *Room) InitMediaSession(c1, c2 *Client) error {
	if _, found := r.ClientList[c1.ID]; !found {
		return fmt.Errorf("client id not in found in %s's client list. ClientList must be in the rooms ClientList to participate", r.Name)
	}
	if _, found := r.ClientList[c2.ID]; !found {
		return fmt.Errorf("client id not in found in %s client list. ClientList must be in the rooms ClientList to participate", r.Name)
	}
	err := c1.initRTCSession(r)
	if err != nil {
		return err
	}
}
