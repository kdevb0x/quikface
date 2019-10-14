// Copyright 2019 kdevb0x Ltd. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause license
// The full license text can be found in the LICENSE file.

package quikface

import (
	"html/template"
	"net/http"

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
	r.HandleFunc("/session/{create:(?:create)}", CreateRoomHandler).Methods("POST").Name("create")
	r.HandleFunc("/session/{join:(?:join)}", JoinRoomHandler).Methods("POST").Queries("roomname")
	r.Handle("/login", &templateHandler{filename: "login.html"})
	r.HandleFunc("/auth/{action}/{provider}", LoginHandler)
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

func JoinRoomHandler(w http.ResponseWriter, r *http.Request) {
	var c *Client
	clientAddr := r.RemoteAddr

	// this is really unsafe without sanitization, removing for now.
	/*
		if displayname, err := r.Cookie("DisplayName"); err == nil {
			c = NewClient(displayname.Value)
		}
	*/

	c = NewClient()
	if c == nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	c.Addr = clientAddr
	if _, exists := mux.Vars(r)["join"]; exists {

		http.Error(w, "the requested room doesn't exist.", http.StatusNotFound)
	}
}

func CreateRoomHandler(w http.ResponseWriter, r *http.Request) {
	var authchan = make(chan authdata)
	// send client to oauth
	authReqHandler(w, r, authchan)
	authd := <-authchan

}

func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	go authReqHandler(w, r)
}
