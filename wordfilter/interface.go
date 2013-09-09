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
