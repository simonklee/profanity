package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	//"strconv"
	"github.com/simonz05/profanity/filter"
	"sync"
)

type PFilter struct {
	f  *map[string]*filter.Filter
	mu sync.RWMutex
}

func NewPFilter() *PFilter {
	m := make(map[string]*filter.Filter, 10)
	return &PFilter{
		f: &m,
	}
}

func (pf *PFilter) addLang(lang string) *filter.Filter {
	pf.mu.Lock()
	m := make(map[string]*filter.Filter, len(*pf.f))
	for k, v := range *(pf.f) {
		m[k] = v
	}
	f := filter.NewFilter()
	m[lang] = f
	pf.f = &m
	pf.mu.Unlock()
	return f
}

func (pf *PFilter) Get(lang string) *filter.Filter {
	f, ok := (*pf.f)[lang]

	if !ok {
		f = pf.addLang(lang)
	}
	return f
}

func JsonError(w http.ResponseWriter, error string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	fmt.Fprintf(w, `{error: %s}`, error)
}

type ErrorResponse struct {
	Error string `json:error`
	Code  int    `json:code`
}

type SanitizeResponse struct {
	Text string `json:"text"`
}

type BlacklistResponse struct {
	Blacklist []string `json:"blacklist"`
	Total     int      `json:"total"`
}

func sanitizeHandle(w http.ResponseWriter, r *http.Request) {
	lang := r.FormValue("lang")
	if lang == "" {
		JsonError(w, "Invalid lang", 400)
		return
	}

	text := r.FormValue("text")
	sanitized := pfilter.Get(lang).Sanitize(text)
	Logf("lang: %s, text: %s, sanitized: %s", lang, text, sanitized)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(&SanitizeResponse{Text: sanitized})
}

func updateBlacklistHandle(w http.ResponseWriter, r *http.Request) {
	lang := r.FormValue("lang")
	if lang == "" {
		JsonError(w, "Invalid lang", 400)
		return
	}
	blacklist, ok := r.Form["blacklist"]

	if !ok || len(blacklist) == 0 {
		JsonError(w, "Expected `blacklist` key", 400)
		return
	}

	switch r.Method {
	case "PUT":
		pfilter.Get(lang).Update(blacklist)
		w.WriteHeader(200)
	case "POST":
		pfilter.Get(lang).Replace(blacklist)
		w.WriteHeader(201)
	default:
		panic("should not reach")
	}
}

func removeBlacklistHandle(w http.ResponseWriter, r *http.Request) {
	lang := r.FormValue("lang")
	if lang == "" {
		JsonError(w, "Invalid lang", 400)
		return
	}
	blacklist, ok := r.Form["blacklist"]

	if !ok || len(blacklist) == 0 {
		JsonError(w, "Expected `blacklist` key", 400)
		return
	}

	pfilter.Get(lang).Remove(blacklist)
	w.WriteHeader(200)
}

func getBlacklistHandle(w http.ResponseWriter, r *http.Request) {
	lang := r.FormValue("lang")
	if lang == "" {
		JsonError(w, "Invalid lang", 400)
		return
	}

	//count, err := strconv.Atoi(r.FormValue("count"))
	//if err != nil {
	//	count = 20
	//}
	//offset, err := strconv.Atoi(r.FormValue("offset"))
	//if err != nil {
	//	offset = 0
	//}

	blacklist := pfilter.Get(lang).Blacklist()
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	//println(offset, count)

	Logf("lang: %s, blacklist: %v", lang, blacklist)

	resp := &BlacklistResponse{
		Blacklist: blacklist,
		Total:     len(blacklist),
	}

	json.NewEncoder(w).Encode(resp)
}
