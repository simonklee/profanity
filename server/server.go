package server

import (
	"net"
	"net/http"
	"os"
	"os/signal"

	"github.com/gorilla/mux"
)

var (
	Version = "0.0.1"
	router  *mux.Router
	pfilter *PFilter
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

func setupServer(filename string) {
	//pfilter = filter.NewFilter()
	pfilter = NewPFilter()

	if filename != "" {
		//loadFromFile(filename)
	}

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

	Logf("Listen on %s", l.Addr())

	sigTrapCloser(l)
	err = http.Serve(l, nil)
	Logf("Shutting down ..")
	return err
}
