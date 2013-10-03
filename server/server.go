package server

import (
	"net"
	"net/http"
	"os"
	"os/signal"

	"github.com/gorilla/mux"
	"github.com/simonz05/profanity/db"
	"github.com/simonz05/util/log"
)

var (
	Version = "0.1.0"
	router  *mux.Router
	filters *profanityFilters
	dbConn  db.Conn
)

func sigTrapCloser(l net.Listener) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for _ = range c {
			l.Close()
			log.Printf("Closed listener %s", l.Addr())
		}
	}()
}

func setupServer(dsn string) (err error) {
	dbConn, err = db.Open(dsn)

	if err != nil {
		return
	}

	filters = newProfanityFilters()

	// HTTP endpoints
	router = mux.NewRouter()
	router.HandleFunc("/api/1.0/sanitize/", sanitizeHandle).Methods("GET").Name("sanitize")
	router.HandleFunc("/api/1.0/blacklist/", updateBlacklistHandle).Methods("POST", "PUT").Name("blacklist")
	router.HandleFunc("/api/1.0/blacklist/remove/", removeBlacklistHandle).Methods("POST", "PUT").Name("blacklist")
	router.HandleFunc("/api/1.0/blacklist/", getBlacklistHandle).Methods("GET").Name("blacklist")
	router.StrictSlash(false)
	http.Handle("/", router)
	return
}

func ListenAndServe(laddr, dsn string) error {
	setupServer(dsn)

	l, err := net.Listen("tcp", laddr)

	if err != nil {
		return err
	}

	log.Printf("Listen on %s", l.Addr())

	sigTrapCloser(l)
	err = http.Serve(l, nil)
	log.Print("Shutting down ..")
	return err
}
