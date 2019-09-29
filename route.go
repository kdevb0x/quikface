// Copyright 2019 kdevb0x Ltd. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause license
// The full license text can be found in the LICENSE file.

package quikface

import (
	"html/template"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type SessionRouter struct {
	Active *RoomList // globaldir
	// after / recieves a new "create session" req, it creates a new Client,
	// then adds it to the IncommingReq queue.
	IncommingReq chan *Client

	httprouter *mux.Router
}

func NewSessionRouter() *SessionRouter {
	s := &SessionRouter{
		Active:       NewRoomList(),
		IncommingReq: make(chan *Client),
	}
	r := mux.NewRouter()
	r.HandleFunc("/", IndexHandler)

	s.httprouter = r
	return s
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	idxtempl, err := template.ParseFiles("assets/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	idxtempl.Execute(w, nil)

}

func CreateRoomHandler(w http.ResponseWriter, r *http.Request) {
	clientAddr := r.RemoteAddr
	c := NewClient(time.Now().UTC().String())
	if c == nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	c.Addr = clientAddr
}