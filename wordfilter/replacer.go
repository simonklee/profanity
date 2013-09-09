package wordfilter

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/simonz05/profanity/util"
)

// starmap used to draw N stars in place of a blacklisted word.
var starmap [16]string

func init() {
	for i := 1; i < 16; i++ {
		starmap[i] = fmt.Sprintf("%s*", starmap[i-1])
	}
}

// A thread-safe word filter
type Replacer struct {
	repl   *strings.Replacer
	replMu sync.RWMutex // repl locker
}

// Returns a new word filter. The word filter is empty by default.
func NewReplacer() *Replacer {
	return &Replacer{
		repl: strings.NewReplacer(),
	}
}

// reload wordlist
func (p *Replacer) Reload(words []string) error {
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
func (p *Replacer) buildReplacer(words []string) (*strings.Replacer, error) {
	var starindex int
	n := len(words) * 2

	if n == 0 {
		return nil, errors.New("Got empty blacklist")
	}

	repl := make([]string, n)

	for i, w := range words {
		repl[i*2] = w
		starindex = util.IntMin(len(w), len(starmap)-1)
		repl[i*2+1] = starmap[starindex]
	}

	return strings.NewReplacer(repl...), nil
}

// Returns a copy of string v where each word in the text that matches a word
// in the blacklist is replaced by ****.
func (p *Replacer) Sanitize(v string) string {
	p.replMu.RLock()
	defer p.replMu.RUnlock()
	return p.repl.Replace(v)
}
