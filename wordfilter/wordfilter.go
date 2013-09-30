package wordfilter

import (
	"github.com/simonz05/profanity/wordlist"
)

// ProfanityFilter is a implements a word filter. It takes a list of words
// which are used to sanitize text. The sanitizer will replace all words which
// match a word in the list with **** (stars).
type ProfanityFilter interface {
	wordlist.Wordlist
	Sanitize(v string) string
}

// Wordfilter implements the ProfanityFilter interface.
type Wordfilter struct {
	List     wordlist.Wordlist
	Replacer *Replacer
}

func NewWordfilter(list wordlist.Wordlist) *Wordfilter {
	return &Wordfilter{
		List:     list,
		Replacer: NewReplacer(),
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

	return w.Replacer.Reload(strings)
}

// Reset the wordlist
func (w *Wordfilter) Empty() error {
	if err := w.List.Empty(); err != nil {
		return err
	}

	return w.Replacer.Reload([]string{})
}

// Reset the wordlist
func (w *Wordfilter) Sanitize(v string) string {
	return w.Replacer.Replace(v)
}
