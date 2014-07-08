package wordfilter

import (
	"errors"
	"strings"
	"sync"

	"github.com/simonz05/util/math"
)

// A thread-safe word filter
type SetReplacer struct {
	repl   map[string]string
	replMu sync.RWMutex // repl locker
}

// Returns a new word filter. The word filter is empty by default.
func NewSetReplacer() *SetReplacer {
	return &SetReplacer{
		repl: make(map[string]string),
	}
}

// reload wordlist
func (p *SetReplacer) Reload(words []string) error {
	repl, err := p.buildReplacer(words)

	if err != nil {
		return err
	}

	p.replMu.Lock()
	p.repl = repl
	p.replMu.Unlock()
	return nil
}

// Build string replacer from blacklist
func (p *SetReplacer) buildReplacer(words []string) (map[string]string, error) {
	var starindex int
	n := len(words)

	if n == 0 {
		return nil, errors.New("Got empty blacklist")
	}

	repl := make(map[string]string, n)

	for _, w := range words {
		starindex = math.IntMin(len(w), len(starmap)-1)
		repl[strings.ToLower(w)] = starmap[starindex]
	}

	return repl, nil
}

// Returns a copy of string v where each word in the text that matches a word
// in the blacklist is replaced by ****.
func (p *SetReplacer) Replace(v string) string {
	p.replMu.RLock()
	defer p.replMu.RUnlock()
	buf := make(appendSliceWriter, 0, len(v))
	p.WriteString(&buf, v)
	return string(buf)
}

func (p *SetReplacer) WriteString(buf *appendSliceWriter, s string) {
	sepCr := "\r\n"
	sepNl := "\n"
	sepSpace := " "

	start := 0
	var sep string

	for i := 0; i+2 < len(s); i++ {
		if s[i] == sepNl[0] {
			sep = "\n"
		} else if s[i] == sepSpace[0] {
			sep = " "
		} else if s[i:i+2] == sepCr {
			sep = sepCr
		} else {
			continue
		}

		buf.WriteString(p.replace(s[start:i]))
		buf.WriteString(sep)
		start = i + len(sep)
		i += len(sep) - 1
	}

	if strings.HasSuffix(s, sepNl) {
		sep = sepNl
	} else if strings.HasSuffix(s, sepSpace) {
		sep = sepSpace
	} else if strings.HasSuffix(s, sepCr) {
		sep = sepSpace
	} else {
		sep = ""
	}

	buf.WriteString(p.replace(s[start : len(s)-len(sep)]))
	buf.WriteString(sep)
}

func (p *SetReplacer) replace(word string) string {
	if stars, ok := p.repl[strings.ToLower(word)]; ok {
		return stars
	}

	return word
}
