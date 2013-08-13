package filter

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
)

// A thread-safe profanity filter
type Filter struct {
	blacklist []string
	blackMu   sync.RWMutex // blacklist locker
	repl      *strings.Replacer
	replMu    sync.RWMutex // repl locker
}

// starmap used to draw N stars in place of a blacklisted word. 
var starmap [16]string

func init() {
	for i := 1; i < 16; i++ {
		starmap[i] = fmt.Sprintf("%s*", starmap[i-1])
	}
}

func imin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func imax(a, b int) int {
	if a < b {
		return b
	}
	return a
}

// Returns a new word filter. The word filter is empty by default.
func NewFilter() *Filter {
	return &Filter{
		blacklist: []string{},
		repl:      strings.NewReplacer(),
	}
}

func printSlice(slice []string) {
	for i, x := range slice {
		fmt.Printf("%d#: %s\n", i, x)
	}
}

// Merge data into an asc sorted slice. Extend the slice if necessary.
func merge(slice, data []string) []string {
	l := len(slice)
	//println(len(data), len(slice), cap(slice))

	if l+len(data) > cap(slice) { // reallocate
		newSlice := make([]string, l, (l + len(data)*2))
		copy(newSlice, slice)
		slice = newSlice
	}

	n := len(data)

	for _, c := range data {
		i := sort.SearchStrings(slice, c)

		// if the word already exists in slice
		if i < len(slice) && slice[i] == c {
			continue
		}

		// if it's the last elem simply append
		if i == n {
			slice = slice[0 : len(slice)+1]
			slice[i] = c
		} else {
			// insert into the slice
			slice = slice[0 : len(slice)+1]
			copy(slice[i+1:], slice[i:])
			slice[i] = c
		}
	}

	return slice
}

// Replace blacklist with the slice input 
func (p *Filter) Replace(blacklist []string) error {
	sort.Strings(blacklist)
	return p.reload(blacklist)
}

// Update blacklist with words from slice
func (p *Filter) Update(blacklist []string) error {
	p.blackMu.RLock()
	n := len(p.blacklist)
	newBlacklist := make([]string, n, n+len(blacklist))
	copy(newBlacklist, p.blacklist)
	p.blackMu.RUnlock()

	if !sort.StringsAreSorted(newBlacklist) {
		sort.Strings(newBlacklist)
	}

	blacklist = merge(newBlacklist, blacklist)
	return p.reload(blacklist)
}

// Delete words in blacklist which match words in input.
func (p *Filter) Remove(blacklist []string) error {
	if len(blacklist) == 0 {
		return nil
	}

	slice := p.Blacklist()
	oldLen := len(slice)

	for _, c := range blacklist {
		i := sort.SearchStrings(slice, c)

		if i == -1 {
			continue
		}

		slice = append(slice[:i], slice[i+1:]...)
	}

	// nothing changed
	if oldLen == len(slice) {
		return nil
	}

	return p.reload(slice)
}

// reload blacklist
func (p *Filter) reload(blacklist []string) error {
	repl, err := p.buildReplacer(blacklist)

	if err != nil {
		return err
	}

	p.blackMu.Lock()
	p.blacklist = blacklist
	p.blackMu.Unlock()
	p.replMu.Lock()
	p.repl = repl
	p.replMu.Unlock()
	return nil
}

// Build string replacer from blacklist
func (p *Filter) buildReplacer(blacklist []string) (*strings.Replacer, error) {
	var starindex int
	n := len(blacklist) * 2

	if n == 0 {
		return nil, errors.New("Got empty blacklist")
	}

	repl := make([]string, n)

	for i, w := range blacklist {
		repl[i*2] = w
		starindex = imin(len(w), len(starmap)-1)
		repl[i*2+1] = starmap[starindex]
	}

	return strings.NewReplacer(repl...), nil
}

// Returns a copy of string v where each word in the text that matches a word
// in the blacklist is replaced by ****.
func (p *Filter) Sanitize(v string) string {
	p.replMu.RLock()
	defer p.replMu.RUnlock()
	return p.repl.Replace(v)
}

// Returns a copy of the current blacklist.
func (p *Filter) Blacklist() []string {
	p.blackMu.RLock()
	defer p.blackMu.RUnlock()
	newSlice := make([]string, len(p.blacklist))
	copy(newSlice, p.blacklist)
	return newSlice
}

// Returns the blacklist lenght.
func (p *Filter) BlacklistLen() int {
	p.blackMu.RLock()
	defer p.blackMu.RUnlock()
	return len(p.blacklist)
}
