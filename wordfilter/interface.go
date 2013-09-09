package wordfilter

import (
	"github.com/simonz05/profanity/wordlist"
)

// ProfanityFilter is an interface.
type ProfanityFilter interface {
	wordlist.Wordlist
	Sanitize(v string) string 
}
