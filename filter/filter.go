package filter

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"sort"
)

// a thread-safe profanity filter
type Filter struct {
	blacklist []string
	blackMu    sync.RWMutex // blacklist locker
	repl      *strings.Replacer
	replMu    sync.RWMutex // repl locker
}

// used to optimize reloading
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

func NewFilter() *Filter {
	return &Filter{
		blacklist: []string{},
		repl: strings.NewReplacer(),
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

    if l + len(data) > cap(slice) {  // reallocate
        newSlice := make([]string, l, (l+len(data)*2))
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
			slice = slice[0:len(slice)+1]
			slice[i] = c
		} else {
			// insert into the slice
			slice = slice[0:len(slice)+1]
			copy(slice[i+1:], slice[i:])
			slice[i] = c
		}
	}

    return slice
}

func (p *Filter) Replace(blacklist []string) error {
	return p.reload(blacklist, false)
}

// replace the current blacklist with a new one
func (p *Filter) Update(blacklist []string) error {
	return p.reload(blacklist, true)
}

func (p *Filter) reload(blacklist []string, update bool) error {
	p.blackMu.Lock()
	defer p.blackMu.Unlock()

	if update {
		n := len(p.blacklist)
		newBlacklist := make([]string, n, n+len(blacklist))
		copy(newBlacklist, p.blacklist)
		// we dont want to hold lock while merging blacklists

		if !sort.StringsAreSorted(newBlacklist) {
			sort.Strings(newBlacklist)
		}

		blacklist = merge(newBlacklist, blacklist)
	}

	repl, err := p.buildReplacer(blacklist)

	if err != nil {
		return err
	}

	p.replMu.Lock()
	defer p.replMu.Unlock()
	p.blacklist = blacklist
	p.repl = repl
	return nil
}

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

func (p *Filter) Sanitize(v string) string {
	p.replMu.RLock()
	defer p.replMu.RUnlock()
	return p.repl.Replace(v)
}

func (p *Filter) Blacklist() []string {
	p.blackMu.RLock()
	defer p.blackMu.RUnlock()
	newSlice := make([]string, len(p.blacklist))
	copy(newSlice, p.blacklist)
	return newSlice
}

func (p *Filter) BlacklistLen() int {
	p.blackMu.RLock()
	defer p.blackMu.RUnlock()
	return len(p.blacklist)
}
