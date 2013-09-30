package wordfilter

import (
	"errors"
	"fmt"
	"io"
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
	repl   *genericReplacer
	replMu sync.RWMutex // repl locker
}

// Returns a new word filter. The word filter is empty by default.
func NewReplacer() *Replacer {
	return &Replacer{
		repl: makeGenericReplacer([]string{}),
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
func (p *Replacer) buildReplacer(words []string) (*genericReplacer, error) {
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

	return makeGenericReplacer(repl), nil
}

// Returns a copy of string v where each word in the text that matches a word
// in the blacklist is replaced by ****.
func (p *Replacer) Replace(v string) string {
	p.replMu.RLock()
	defer p.replMu.RUnlock()
	return p.repl.Replace(v)
}

// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// trieNode is a node in a lookup trie for prioritized key/value pairs. Keys
// and values may be empty. For example, the trie containing keys "ax", "ay",
// "bcbc", "x" and "xy" could have eight nodes:
//
//  n0  -
//  n1  a-
//  n2  .x+
//  n3  .y+
//  n4  b-
//  n5  .cbc+
//  n6  x+
//  n7  .y+
//
// n0 is the root node, and its children are n1, n4 and n6; n1's children are
// n2 and n3; n4's child is n5; n6's child is n7. Nodes n0, n1 and n4 (marked
// with a trailing "-") are partial keys, and nodes n2, n3, n5, n6 and n7
// (marked with a trailing "+") are complete keys.
type trieNode struct {
	// value is the value of the trie node's key/value pair. It is empty if
	// this node is not a complete key.
	value string
	// priority is the priority (higher is more important) of the trie node's
	// key/value pair; keys are not necessarily matched shortest- or longest-
	// first. Priority is positive if this node is a complete key, and zero
	// otherwise. In the example above, positive/zero priorities are marked
	// with a trailing "+" or "-".
	priority int

	// A trie node may have zero, one or more child nodes:
	//  * if the remaining fields are zero, there are no children.
	//  * if prefix and next are non-zero, there is one child in next.
	//  * if table is non-zero, it defines all the children.
	//
	// Prefixes are preferred over tables when there is one child, but the
	// root node always uses a table for lookup efficiency.

	// prefix is the difference in keys between this trie node and the next.
	// In the example above, node n4 has prefix "cbc" and n4's next node is n5.
	// Node n5 has no children and so has zero prefix, next and table fields.
	prefix string
	next   *trieNode

	// table is a lookup table indexed by the next byte in the key, after
	// remapping that byte through genericReplacer.mapping to create a dense
	// index. In the example above, the keys only use 'a', 'b', 'c', 'x' and
	// 'y', which remap to 0, 1, 2, 3 and 4. All other bytes remap to 5, and
	// genericReplacer.tableSize will be 5. Node n0's table will be
	// []*trieNode{ 0:n1, 1:n4, 3:n6 }, where the 0, 1 and 3 are the remapped
	// 'a', 'b' and 'x'.
	table []*trieNode
}

func (t *trieNode) add(key, val string, priority int, r *genericReplacer) {
	if key == "" {
		if t.priority == 0 {
			t.value = val
			t.priority = priority
		}
		return
	}

	if t.prefix != "" {
		// Need to split the prefix among multiple nodes.
		var n int // length of the longest common prefix
		for ; n < len(t.prefix) && n < len(key); n++ {
			if t.prefix[n] != key[n] {
				break
			}
		}
		if n == len(t.prefix) {
			t.next.add(key[n:], val, priority, r)
		} else if n == 0 {
			// First byte differs, start a new lookup table here. Looking up
			// what is currently t.prefix[0] will lead to prefixNode, and
			// looking up key[0] will lead to keyNode.
			var prefixNode *trieNode
			if len(t.prefix) == 1 {
				prefixNode = t.next
			} else {
				prefixNode = &trieNode{
					prefix: t.prefix[1:],
					next:   t.next,
				}
			}
			keyNode := new(trieNode)
			t.table = make([]*trieNode, r.tableSize)
			t.table[r.mapping[t.prefix[0]]] = prefixNode
			t.table[r.mapping[key[0]]] = keyNode
			t.prefix = ""
			t.next = nil
			keyNode.add(key[1:], val, priority, r)
		} else {
			// Insert new node after the common section of the prefix.
			next := &trieNode{
				prefix: t.prefix[n:],
				next:   t.next,
			}
			t.prefix = t.prefix[:n]
			t.next = next
			next.add(key[n:], val, priority, r)
		}
	} else if t.table != nil {
		// Insert into existing table.
		m := r.mapping[key[0]]
		if t.table[m] == nil {
			t.table[m] = new(trieNode)
		}
		t.table[m].add(key[1:], val, priority, r)
	} else {
		t.prefix = key
		t.next = new(trieNode)
		t.next.add("", val, priority, r)
	}
}

func (r *genericReplacer) lookup(s string, ignoreRoot bool) (val string, keylen int, found bool) {
	// Iterate down the trie to the end, and grab the value and keylen with
	// the highest priority.
	bestPriority := 0
	node := &r.root
	n := 0
	for node != nil {
		if node.priority > bestPriority && !(ignoreRoot && node == &r.root) {
			bestPriority = node.priority
			val = node.value
			keylen = n
			found = true
		}

		if s == "" {
			break
		}
		if node.table != nil {
			index := r.mapping[s[0]]
			if int(index) == r.tableSize {
				break
			}
			node = node.table[index]
			s = s[1:]
			n++
		} else if node.prefix != "" && strings.HasPrefix(s, node.prefix) {
			n += len(node.prefix)
			s = s[len(node.prefix):]
			node = node.next
		} else {
			break
		}
	}
	return
}

// genericReplacer is the fully generic algorithm.
// It's used as a fallback when nothing faster can be used.
type genericReplacer struct {
	root trieNode
	// tableSize is the size of a trie node's lookup table. It is the number
	// of unique key bytes.
	tableSize int
	// mapping maps from key bytes to a dense index for trieNode.table.
	mapping [256]byte
}

func makeGenericReplacer(oldnew []string) *genericReplacer {
	r := new(genericReplacer)

	for i := 0; i < len(oldnew); i++ {
		oldnew[i] = strings.ToLower(oldnew[i])
	}

	// Find each byte used, then assign them each an index.
	for i := 0; i < len(oldnew); i += 2 {
		key := oldnew[i]
		for j := 0; j < len(key); j++ {
			r.mapping[key[j]] = 1
		}
	}

	for _, b := range r.mapping {
		r.tableSize += int(b)
	}

	var index byte
	for i, b := range r.mapping {
		if b == 0 {
			r.mapping[i] = byte(r.tableSize)
		} else {
			r.mapping[i] = index
			index++
		}
	}
	// Ensure root node uses a lookup table (for performance).
	r.root.table = make([]*trieNode, r.tableSize)

	for i := 0; i < len(oldnew); i += 2 {
		r.root.add(oldnew[i], oldnew[i+1], len(oldnew)-i, r)
	}
	return r
}

type appendSliceWriter []byte

// Write writes to the buffer to satisfy io.Writer.
func (w *appendSliceWriter) Write(p []byte) (int, error) {
	*w = append(*w, p...)
	return len(p), nil
}

// WriteString writes to the buffer without string->[]byte->string allocations.
func (w *appendSliceWriter) WriteString(s string) (int, error) {
	*w = append(*w, s...)
	return len(s), nil
}

type stringWriterIface interface {
	WriteString(string) (int, error)
}

type stringWriter struct {
	w io.Writer
}

func (w stringWriter) WriteString(s string) (int, error) {
	return w.w.Write([]byte(s))
}

func getStringWriter(w io.Writer) stringWriterIface {
	sw, ok := w.(stringWriterIface)
	if !ok {
		sw = stringWriter{w}
	}
	return sw
}

func (r *genericReplacer) Replace(s string) string {
	buf := make(appendSliceWriter, 0, len(s))
	r.WriteString(&buf, s)
	return string(buf)
}

func (r *genericReplacer) WriteString(w io.Writer, s string) (n int, err error) {
	sw := getStringWriter(w)
	var last, wn int
	var prevMatchEmpty bool
	for i := 0; i <= len(s); {
		// Ignore the empty match iff the previous loop found the empty match.
		val, keylen, match := r.lookup(strings.ToLower(s[i:]), prevMatchEmpty)
		prevMatchEmpty = match && keylen == 0
		if match {
			wn, err = sw.WriteString(s[last:i])
			n += wn
			if err != nil {
				return
			}
			wn, err = sw.WriteString(val)
			n += wn
			if err != nil {
				return
			}
			i += keylen
			last = i
			continue
		}
		i++
	}
	if last != len(s) {
		wn, err = sw.WriteString(s[last:])
		n += wn
	}
	return
}
