package wordlist

// Wordlist is an index of words.
type Wordlist interface {
	// Return a count of total entries in wordlist
	Count() (int, error)

	// Return `count` entries from `offset`
	Get(count, offset int) ([]string, error)

	// Add or overwrite words
	Set(words []string) error

	// Delete words
	Delete(words []string) error

	// Replace wordlist with `words`
	Replace(words []string) error

	// Reset the wordlist
	Empty() error
}
