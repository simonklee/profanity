// Copyright 2014 Simon Zimmermann. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

/*
Package util/handler is a collection of HTTP handlers for net/http package
*/

package handler

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/gorilla/context"
	"github.com/simonz05/util/log"
	"github.com/simonz05/util/session"
	"github.com/tideland/goas/v2/monitoring"
)

// Creates a stack of HTTP handlers. Each HTTP handler is responsible for
// calling the next. The handlers are executed in reverse order, the last is
// called first.
func Use(handler http.Handler, middleware ...func(http.Handler) http.Handler) http.Handler {
	for _, m := range middleware {
		handler = m(handler)
	}
	return handler
}

// LogHandler adds logging to http requests
func LogHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("%s %s", r.Method, r.URL.Path)
		h.ServeHTTP(w, r)
		//log.Println(r.Method, r.URL, r.Form, time.Since(start))
		log.Printf("%s %s in %v\n", r.Method, r.URL.Path, time.Since(start))
	})
}

// DebugHandle
func DebugHandle(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := httputil.DumpRequest(r, true)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		log.Debug(string(data))
		h.ServeHTTP(w, r)
	})
}

// MeasureHandler adds measuring to http requests
func MeasureHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m := monitoring.BeginMeasuring(r.URL.Path)
		h.ServeHTTP(w, r)
		m.EndMeasuring()
	})
}

func RecoveryHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Errorln("Recovered:", rec)
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()
		h.ServeHTTP(w, r)
	})
}

// LogHandler adds logging to http requests
func NewCORSHandler(domains ...string) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// TODO: compare origin header to a list of accepted domains
			// instead of sending all OK domains.
			for _, domain := range domains {
				w.Header().Add("Access-Control-Allow-Origin", domain)
			}
			h.ServeHTTP(w, r)
		})
	}
}

type contextKey int

const (
	sessionKey contextKey = iota
	tokenKey
)

const (
	SessionHeader = "Authorization-Session"
	TokenHeader   = "Authorization-Token"
)

type authHandler struct {
	handler    http.Handler
	backend    session.Storage
	contextKey contextKey
	headerKey  string
	mustAuth   bool
}

func (a *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := a.loadSession(r); err != nil {
		log.Println("handler error: ", err)

		if a.mustAuth {
			log.Println("Status Unauthorized")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}
	a.handler.ServeHTTP(w, r)
	//rec := httptest.NewRecorder()
	//a.handler.ServeHTTP(rec, r)

	//for k, v := range rec.Header() {
	//	w.Header()[k] = v
	//}

	//w.Header().Set("X-Authenticated-By", "kogama")
	//w.WriteHeader(rec.Code)
	//rec.Body.WriteTo(w)
}

// Tries to load a session from the auth backend for the given authentication
func (a *authHandler) loadSession(r *http.Request) error {
	value := r.Header[a.headerKey]
	var id string

	if len(value) > 0 {
		id = value[0]
	} else {
		id = r.URL.Query().Get("session")
	}

	if id == "" {
		return fmt.Errorf("Header %s and session was empty", a.headerKey)
	}

	ses, err := a.backend.Read(id)

	if err != nil {
		return err
	}

	context.Set(r, a.contextKey, ses)
	return nil
}

// NewAuthSessionHandler creates a AuthSessionHandler for the specified
// backend. It loads a session from a session key.
func NewAuthSessionHandler(sessionStorage session.Storage, must bool) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return &authHandler{
			handler:    h,
			backend:    sessionStorage,
			contextKey: sessionKey,
			headerKey:  SessionHeader,
			mustAuth:   must,
		}
	}
}

// NewAuthTokenHandler creates a AuthTokenHandler for the specified
// It loads a session from a token key.
func NewAuthTokenHandler(tokenStorage session.Storage) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return &authHandler{
			handler:    h,
			backend:    tokenStorage,
			contextKey: tokenKey,
			headerKey:  TokenHeader,
			mustAuth:   true,
		}
	}
}

// CurrentSession returns the matched session for the current request, if any.
func CurrentSession(r *http.Request) *session.Session {
	if rv := context.Get(r, sessionKey); rv != nil {
		return rv.(*session.Session)
	}

	return nil
}

// CurrentToken returns the matched token for the current request, if any.
func CurrentToken(r *http.Request) *session.Session {
	if rv := context.Get(r, tokenKey); rv != nil {
		return rv.(*session.Session)
	}

	return nil
}
