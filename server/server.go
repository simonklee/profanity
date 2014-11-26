package server

import (
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/simonz05/profanity/db"
	"github.com/simonz05/util/handler"
	"github.com/simonz05/util/log"
	"github.com/simonz05/util/sig"
)

var (
	Version = "0.1.0"
	router  *mux.Router
	filters *profanityFilters
	dbConn  db.Conn
)

func setupServer(dsn string) (err error) {
	dbConn, err = db.Open(dsn)

	if err != nil {
		return
	}

	filters = newProfanityFilters()

	// HTTP endpoints
	router = mux.NewRouter()
	router.HandleFunc("/v1/profanity/sanitize/", sanitizeHandle).Methods("GET").Name("sanitize")
	router.HandleFunc("/v1/profanity/blacklist/", updateBlacklistHandle).Methods("POST", "PUT").Name("blacklist")
	router.HandleFunc("/v1/profanity/blacklist/remove/", removeBlacklistHandle).Methods("POST", "PUT").Name("blacklist")
	router.HandleFunc("/v1/profanity/blacklist/", getBlacklistHandle).Methods("GET").Name("blacklist")
	router.StrictSlash(false)

	// global middleware
	var middleware []func(http.Handler) http.Handler

	switch log.Severity {
	case log.LevelDebug:
		middleware = append(middleware, handler.LogHandler, handler.DebugHandle, handler.RecoveryHandler)
	default:
		middleware = append(middleware, handler.LogHandler, handler.RecoveryHandler)
	}

	wrapped := handler.Use(router, middleware...)
	http.Handle("/", wrapped)
	return
}

func ListenAndServe(laddr, dsn string) error {
	setupServer(dsn)

	l, err := net.Listen("tcp", laddr)

	if err != nil {
		return err
	}

	log.Printf("Listen on %s", l.Addr())

	sig.TrapCloser(l)
	err = http.Serve(l, nil)
	log.Print("Shutting down ..")
	return err
}
