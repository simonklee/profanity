package server

import (
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/simonz05/profanity/config"
	"github.com/simonz05/profanity/db"
	"github.com/simonz05/profanity/types"
	"github.com/simonz05/profanity/wordfilter"
	"github.com/simonz05/profanity/wordlist"
	"github.com/simonz05/util/handler"
	"github.com/simonz05/util/log"
	"github.com/simonz05/util/sig"
)

var (
	Version       = "0.1.0"
	router        *mux.Router
	filters       *profanityFilters
	dbConn        db.Conn
	newWordfilter func(list wordlist.Wordlist) *wordfilter.Wordfilter
)

func setupServer(conf *config.Config) (err error) {
	dbConn, err = db.Open(conf.Redis.DSN)

	if err != nil {
		return
	}

	newWordfilter = func(list wordlist.Wordlist) *wordfilter.Wordfilter {
		var replacer wordfilter.Replacer

		switch conf.Filter {
		case types.Any:
			replacer = wordfilter.NewStringReplacer()
		case types.Word:
			replacer = wordfilter.NewSetReplacer()
		default:
			replacer = wordfilter.NewSetReplacer()
		}

		return &wordfilter.Wordfilter{
			List:     list,
			Replacer: replacer,
		}
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

func ListenAndServe(conf *config.Config) error {
	setupServer(conf)

	l, err := net.Listen("tcp", conf.Listen)

	if err != nil {
		return err
	}

	log.Printf("Listen on %s", l.Addr())

	sig.TrapCloser(l)
	err = http.Serve(l, nil)
	log.Print("Shutting down ..")
	return err
}
