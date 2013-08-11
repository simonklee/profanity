package server

import (
	"net"
	"net/http"
	"os"
	"os/signal"
	"io/ioutil"
	"strings"

	"github.com/gorilla/mux"
	"github.com/simonz05/profanity/filter"
)

var (
	Version = "0.0.1"
	profanity *filter.Filter
	router  *mux.Router
)

func sigTrapCloser(l net.Listener) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for _ = range c {
			l.Close()
			Logf("Closed listener %s", l.Addr())
		}
	}()
}

func ListenAndServe(laddr, filename string) error {
	profanity = filter.NewFilter()

	if filename != "" {
		loadFromFile(filename)
	}

	// HTTP endpoints
	router = mux.NewRouter()
	router.HandleFunc("/api/1.0/sanitize/", sanitizeHandle).Methods("GET").Name("sanitize")
	router.HandleFunc("/api/1.0/blacklist/", postBlacklistHandle).Methods("POST", "PUT").Name("blacklist")
	router.HandleFunc("/api/1.0/blacklist/", getBlacklistHandle).Methods("GET").Name("blacklist")
	router.StrictSlash(false)
	http.Handle("/", router)

	l, err := net.Listen("tcp", laddr)

	if err != nil {
		return err
	}

	Logf("Listen on %s", l.Addr())

	sigTrapCloser(l)
	err = http.Serve(l, nil)
	Logf("Shutting down ..")
	return err
}

func loadFromFile(filename string) error {
	content, err := ioutil.ReadFile(filename)

	if err != nil {
		return err
	}

	words := strings.Split(strings.TrimSpace(string(content)), "\n")
	profanity.Reload(words, true)
	return nil
}
