// Copyright 2019 kdevb0x Ltd. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause license
// The full license text can be found in the LICENSE file.

package quikface

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/gorilla/mux"
)

var (
	_ = mux.NewRouter()
	_ = ws.NewMask()
	_ = new(wsutil.ControlHandler)
)

type SessionKey struct {
	Key        []byte
	Expiration time.Time
}

type SessionConfig struct {
	Peers map[string]SessionKey // remote ip addresses of peers on call
}

type WSMuxer interface {
	HandleWS(context.Context, net.Conn)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {

}
