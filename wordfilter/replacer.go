package wordfilter

import (
	"fmt"
	"io"
)

// starmap used to draw N stars in place of a blacklisted word.
var starmap [16]string

func init() {
	for i := 1; i < 16; i++ {
		starmap[i] = fmt.Sprintf("%s*", starmap[i-1])
	}
}

type Replacer interface {
	Replace(v string) string
	Reload(words []string) error
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
