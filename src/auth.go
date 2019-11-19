// Copyright 2019 kdevb0x Ltd. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause license
// The full license text can be found in the LICENSE file.

package quikface

import (
	"crypto/rand"
	"html/template"
	"net/http"
	"path/filepath"
	"sync"

	"golang.org/x/crypto/nacl/sign"

	"github.com/gorilla/mux"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
)

type authdata struct {
	client *Client

	// auth data from provider
	authobj   objx.Map
	hasCookie bool
}

// templ represents a single template
type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

// ServeHTTP handles the HTTP request.
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("assets", t.filename)))
	})
	t.templ.Execute(w, nil)
}

type authHandler struct {
	next http.Handler
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if _, err := r.Cookie("auth"); err == http.ErrNoCookie {
		// not authenticated
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else if err != nil {
		// some other error
	} else {
		// successful auth
		h.next.ServeHTTP(w, r)
	}
}

func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}

func initOmniauth() {
	// gomniauth.SetSecurityKey(/* TODO: add base64 encoded crypto key */)
	publicKey, privateKey, err := sign.GenerateKey(rand.Reader)
	if err != nil {

	}
	gomniauth.SetSecurityKey(string(privateKey[:]))
	gomniauth.WithProviders(
		facebook.New("key", "secret",
			"http://localhost:8080/auth/callback/facebook"),
		github.New("key", "secret",
			"http://localhost:8080/auth/callback/github"),
		google.New("key", "secret",
			"http://localhost:8080/auth/callback/google"),
	)
}

func authReqHandler(w http.ResponseWriter, r *http.Request, authchan ...chan authdata) {
	// format auth/{action}/{provider}

	vars := mux.Vars(r)
	action := vars["action"]
	provider := vars["provider"]
	switch action {
	case "login":
		pvdr, err := gomniauth.Provider(provider)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		loginURL, err := pvdr.GetBeginAuthURL(nil, nil)
		w.Header().Set("Location", loginURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
	case "callback":
		pvdr, err := gomniauth.Provider(provider)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		creds, err := pvdr.CompleteAuth(objx.MustFromURLQuery(r.URL.RawQuery))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		user, err := pvdr.GetUser(creds)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		authcookie := objx.New(map[string]interface{}{
			"name": user.Name(),
		})
		authobj := authdata{
			authobj: authcookie,
		}
		authCookieVal := authcookie.MustBase64()
		http.SetCookie(w, &http.Cookie{
			Name:  "auth",
			Value: authCookieVal,
			Path:  "/",
		})
		if len(authchan) > 0 {
			authchan[0] <- authobj
		}
		w.Header()["Location"] = []string{"/session"}
		w.WriteHeader(http.StatusTemporaryRedirect)

	}

}
