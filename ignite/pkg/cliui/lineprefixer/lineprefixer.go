// Package lineprefixer is a helpers to add prefixes to new lines.
package lineprefixer

import (
	"bytes"
	"io"
)

// Writer is a prefixed line writer.
type Writer struct {
	prefix       func() string
	w            io.Writer
	shouldPrefix bool
}

// NewWriter returns a new Writer that adds prefixes to each line
// written. It then writes prefixed data stream into w.
func NewWriter(w io.Writer, prefix func() string) *Writer {
	return &Writer{
		prefix:       prefix,
		w:            w,
		shouldPrefix: true,
	}
}

// Write implements io.Writer.
func (p *Writer) Write(b []byte) (n int, err error) {
	var (
		numBytes     = len(b)
		lastChar     = b[numBytes-1]
		newLine      = byte('\n')
		snewLine     = []byte{newLine}
		replaceCount = bytes.Count(b, snewLine)
		prefix       = []byte(p.prefix())
	)
	if lastChar == newLine {
		replaceCount--
	}
	b = bytes.Replace(b, snewLine, append(snewLine, prefix...), replaceCount)
	if p.shouldPrefix {
		b = append(prefix, b...)
	}
	p.shouldPrefix = lastChar == newLine
	if _, err := p.w.Write(b); err != nil {
		return 0, err
	}
	return numBytes, nil
}
