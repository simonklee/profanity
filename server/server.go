package server

import (
	"net"
	"net/http"
	"os"
	"os/signal"

	"github.com/gorilla/mux"
	"github.com/simonz05/profanity/util"
)

var (
	Version = "0.1.0"
	router  *mux.Router
	filters *profanityServer
)

func sigTrapCloser(l net.Listener) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for _ = range c {
			l.Close()
			util.Logf("Closed listener %s", l.Addr())
		}
	}()
}

func setupServer(filename string) {
	filters = newServer()

	// HTTP endpoints
	router = mux.NewRouter()
	router.HandleFunc("/api/1.0/sanitize/", sanitizeHandle).Methods("GET").Name("sanitize")
	router.HandleFunc("/api/1.0/blacklist/", updateBlacklistHandle).Methods("POST", "PUT").Name("blacklist")
	router.HandleFunc("/api/1.0/blacklist/remove/", removeBlacklistHandle).Methods("POST", "PUT").Name("blacklist")
	router.HandleFunc("/api/1.0/blacklist/", getBlacklistHandle).Methods("GET").Name("blacklist")
	router.StrictSlash(false)
	http.Handle("/", router)

}

func ListenAndServe(laddr, filename string) error {
	setupServer(filename)

	l, err := net.Listen("tcp", laddr)

	if err != nil {
		return err
	}

	util.Logf("Listen on %s", l.Addr())

	sigTrapCloser(l)
	err = http.Serve(l, nil)
	util.Logf("Shutting down ..")
	return err
}
