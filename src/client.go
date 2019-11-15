// Copyright 2019 kdevb0x Ltd. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause license
// The full license text can be found in the LICENSE file.

package quikface

import (
	"fmt"
	"io"

	"github.com/google/uuid"
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
	Signal   io.ReadWriteCloser // signalFunc
	MsgQueue MsgQueue
}

type MsgQueue struct {
	Inbound  chan Message // inbound messages from client
	Outbound chan Message
}

func NewClient(displayname ...string) *Client {
	var msgq = MsgQueue{make(chan Message, 10), make(chan Message, 10)}
	if len(displayname) > 0 {
		// cant break if stmnts, so had to flip this test
		if displayname[0] != "" {
			return &Client{ID: uuid.New().ID(), DisplayName: displayname[0], MsgQueue: msgq}

		}
		// here displayname[0] == ""
	}
	return &Client{ID: uuid.New().ID(), DisplayName: randomChatName(), MsgQueue: msgq}

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
		room.Clients[c.ID] = c
		return room, nil

	}
	return nil, fmt.Errorf("error: room %s doesn't exist", name)
}
