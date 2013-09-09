package wordfilter

import (
	"github.com/simonz05/profanity/wordlist"
)

// Wordfilter implements the ProfanityFilter interface.
type Wordfilter struct {
	List   wordlist.Wordlist
	Filter *Replacer
}

func NewWordfilter(list wordlist.Wordlist) *Wordfilter {
	return &Wordfilter{
		List:   list,
		Filter: NewReplacer(),
	}
}

// Return a count of total entries in wordlist
func (w *Wordfilter) Count() (int, error) {
	return w.List.Count()
}

// Return `count` entries from `offset`
func (w *Wordfilter) Get(count, offset int) ([]string, error) {
	return w.List.Get(count, offset)
}

// Add or overwrite words
func (w *Wordfilter) Set(words []string) error {
	if err := w.List.Set(words); err != nil {
		return err
	}

	return w.reload()
}

// Delete words
func (w *Wordfilter) Delete(words []string) error {
	if err := w.List.Delete(words); err != nil {
		return err
	}

	return w.reload()
}

// Replace wordlist with `words`
func (w *Wordfilter) Replace(words []string) error {
	if err := w.List.Replace(words); err != nil {
		return err
	}

	return w.reload()
}

func (w *Wordfilter) reload() error {
	cnt, err := w.List.Count()

	if err != nil {
		return err
	}

	strings, err := w.List.Get(cnt, 0)

	if err != nil {
		return err
	}

	return w.Filter.Reload(strings)
}

// Reset the wordlist
func (w *Wordfilter) Empty() error {
	if err := w.List.Empty(); err != nil {
		return err
	}

	return w.Filter.Reload([]string{})
}

// Reset the wordlist
func (w *Wordfilter) Sanitize(v string) string {
	return w.Filter.Sanitize(v) 
}
