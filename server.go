// server.go

// Copyright (C) 2019 Kdevb0x Ltd.
// This software may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.

package vidchat

import (
	"context"

	"github.com/gorilla/schema"
	"github.com/gorilla/websocket"
)

var _ websocket.Upgrader
var _ schema.Decoder

type ChatClient struct {
	Name           string `schema:"display_name"`
	PhoneNumber    string `schema:"phone_number"`
	InitializeCall func(ctx context.Context, partner *ChatClient) (*CallSession, error)
}

type CallSession struct {
	LocalClient  *ChatClient
	RemoteClient *ChatClient
}
