// Copyright 2019 kdevb0x Ltd. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause license
// The full license text can be found in the LICENSE file.

package quikface

import (
	"fmt"
	"strconv"
)

type Client struct {
	ID string
	// for signal mesages; known as "rwc" in collider
	// TODO: Maybe change this to io.ReadWriteCloser and have signalFunc
	// implement the interface?
	Signal signalFunc
	// Signal io.ReadWriteCloser

	MsgQueue []Message
}

func NewClient(id string) *Client {
	return &Client{ID: id, MsgQueue: make([]Message, 0, 10)}
}

func (c *Client) Register(signaler signalFunc) error {
	if c.Signal != nil {
		return fmt.Errorf("duplicate registration; %s already has a signal connection registered", c.ID)
	}
}

// JoinRoom return a pointer to the room if the join was successfull, otherwise
// said pointer is nil, and err explains why.
func (c *Client) JoinRoom(name string) (*Room, error) {
	if hash, exists := masterDirectory.Rooms[name]; exists {
		room, err := masterDirectory.GetRoom(hash)
		if err != nil {
			return nil, fmt.Errorf("error finding hash of %s, %w\n", name, err)
		}
		id, err := strconv.ParseUint(c.ID, 2, 64)
		if err != nil {
			return err
		}
		room.Clients[id] = c
		return room, nil

	}
}
