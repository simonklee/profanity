package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/simonz05/profanity/wordfilter"
	"github.com/simonz05/profanity/wordlist"
	"github.com/simonz05/util/log"
)

type profanityFilters struct {
	lang *map[string]wordfilter.ProfanityFilter
	mu   sync.RWMutex
}

func newProfanityFilters() *profanityFilters {
	m := new(map[string]wordfilter.ProfanityFilter)
	return &profanityFilters{
		lang: m,
	}
}

func (s *profanityFilters) addLang(lang string) wordfilter.ProfanityFilter {
	s.mu.Lock()
	m := make(map[string]wordfilter.ProfanityFilter, len(*s.lang))

	for k, v := range *(s.lang) {
		m[k] = v
	}

	list := wordlist.NewRedisWordlist(dbConn, lang)
	f := wordfilter.NewWordfilter(list)
	m[lang] = f
	s.lang = &m
	s.mu.Unlock()
	f.Reload()
	return f
}

func (s *profanityFilters) get(lang string) wordfilter.ProfanityFilter {
	f, ok := (*s.lang)[lang]

	if !ok {
		f = s.addLang(lang)
	}

	return f
}

func jsonError(w http.ResponseWriter, error string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	fmt.Fprintf(w, `{error: %s}`, error)
}

type errorResponse struct {
	Error string `json:error`
	Code  int    `json:code`
}

type sanitizeResponse struct {
	Text string `json:"text"`
}

type blacklistResponse struct {
	Blacklist []string `json:"blacklist"`
	Total     int      `json:"total"`
}

func sanitizeHandle(w http.ResponseWriter, r *http.Request) {
	lang := r.FormValue("lang")
	if lang == "" {
		jsonError(w, "Invalid lang", 400)
		return
	}

	text := r.FormValue("text")
	sanitized := filters.get(lang).Sanitize(text)
	log.Printf("lang: %s, text: %s, sanitized: %s", lang, text, sanitized)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(&sanitizeResponse{Text: sanitized})
}

func updateBlacklistHandle(w http.ResponseWriter, r *http.Request) {
	log.Printf("update blacklist")
	lang := r.FormValue("lang")
	if lang == "" {
		jsonError(w, "Invalid lang", 400)
		return
	}
	blacklist, ok := r.Form["blacklist"]

	if !ok || len(blacklist) == 0 {
		jsonError(w, "Expected `blacklist` key", 400)
		return
	}

	switch r.Method {
	case "PUT":
		filters.get(lang).Set(blacklist)
		w.WriteHeader(200)
	case "POST":
		filters.get(lang).Replace(blacklist)
		w.WriteHeader(201)
	default:
		panic("should not reach")
	}
}

func removeBlacklistHandle(w http.ResponseWriter, r *http.Request) {
	lang := r.FormValue("lang")
	if lang == "" {
		jsonError(w, "Invalid lang", 400)
		return
	}
	blacklist, ok := r.Form["blacklist"]

	if !ok || len(blacklist) == 0 {
		jsonError(w, "Expected `blacklist` key", 400)
		return
	}

	filters.get(lang).Delete(blacklist)
	w.WriteHeader(200)
}

func getBlacklistHandle(w http.ResponseWriter, r *http.Request) {
	log.Printf("GET blacklist")
	lang := r.FormValue("lang")
	if lang == "" {
		jsonError(w, "Invalid lang", 400)
		return
	}

	count, err := strconv.Atoi(r.FormValue("count"))
	if err != nil {
		count = 20
	}

	offset, err := strconv.Atoi(r.FormValue("offset"))
	if err != nil {
		offset = 0
	}

	log.Printf("lang: %s, count: %d, offset: %d", lang, count, offset)
	// TODO: handle err
	filter := filters.get(lang)
	list, _ := filter.Get(count, offset)
	cnt, _ := filter.Count()

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	resp := &blacklistResponse{
		Blacklist: list,
		Total:     cnt,
	}

	json.NewEncoder(w).Encode(resp)
}
