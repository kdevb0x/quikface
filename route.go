// Copyright 2019 kdevb0x Ltd. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause license
// The full license text can be found in the LICENSE file.

package quikface

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/binary"
	"encoding/base64"
	"log"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"golang.org/x/crypto/nacl/sign"

	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

// DefaultSessionRouter is the default
var DefaultSessionRouter = NewSessionRouter()



type SessionRouter struct {
	Active *RoomList // globaldir
	// after / recieves a new "create session" req, it creates a new Client,
	// then adds it to the IncommingReq queue.
	IncommingReq chan *Client

	httprouter *mux.Router
}

func (s *SessionRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	lastclienthash, err := r.Cookie("last-session-encoded")
	if err != nil {
		log.Println("no last-session found for " + r.RemoteAddr + " err: " + err.Error())
	}

	var lastSesh PreviousClientSession
	r := bytes.NewReader([]byte(lastclienthash.Value))
	d := base64.NewDecoder(base64.StdEncoding,  r)

	client, err :=
}

// NewSessionRouter creates a new SessionRouter.
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
	var authd authdata
	registry := sessions.GetRegistry(r)
	// create session-cookie using remote addr as the name.

	session, err := registry.Get(nil, r.RemoteAddr)

	if err != nil {
		http.Error(w, err.Error(), http.StatusExpectationFailed)
	}
	// keys[0] is used for authentication, keys[1:] are used for encryption
	sessionKey := securecookie.GenerateRandomKey(32)
	encKey := securecookie.GenerateRandomKey(32)
	store := sessions.NewCookieStore(sessionKey, encKey)
	sessions.NewCookie("auth", b64str)

	if err != nil {

		if err == http.ErrNoCookie {
			var authchan = make(chan authdata)
			// send client to oauth
			authReqHandler(w, r, authchan)
			authd = <-authchan

			publicKey, privateKey, err := sign.GenerateKey(rand.Reader)
			if err != nil {
				throwError(err)
			}
			b64str, err := authd.authobj.SignedBase64(string(privateKey[:]))

		}
	}

	if err := r.ParseForm(); err != nil {
		throwError(ErrHTTPRequestParseError, w)
	}
	DefaultSessionRouter.Active.NewRoom()
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	go authReqHandler(w, r)
}

// WrapHandlerWithContext adds the optional key value pairs in ctxVals to
// r.Context via contex.WithValue then calls the handler with w and r in another
// goroutine.
func WrapHandlerWithContext(w http.ResponseWriter, r *http.Request, handler func(w http.ResponseWriter, r *http.Request), ctxVals ...map[interface{}]interface{}) {

	var curCtx context.Context
	if len(ctxVals) > 0 {
		// pull maps from implicit variatric slice
		for _, m := range ctxVals {
			// now range over the map
			curCtx = r.Context()
			for k, v := range m {
				tmpctx := context.WithValue(curCtx, k, v)
				curCtx = tmpctx
			}
		}
		req, err := http.NewRequestWithContext(curCtx, r.Method, r.URL.String(), r.Body)
		if err != nil {
			throwError(err, w)
		}
		r = req

	}
	go handler(w, r)
}


