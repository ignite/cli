package entrywriter

import (
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/pkg/errors"

	"github.com/ignite/cli/ignite/pkg/xstrings"
)

const (
	None = "-"
)

var ErrInvalidFormat = errors.New("invalid entry format")

// MustWrite writes into out the tabulated entries and panic if the entry format is invalid.
func MustWrite(out io.Writer, header []string, entries ...[]string) error {
	err := Write(out, header, entries...)
	if errors.Is(err, ErrInvalidFormat) {
		panic(err)
	}
	return err
}

// Write writes into out the tabulated entries.
func Write(out io.Writer, header []string, entries ...[]string) error {
	w := &tabwriter.Writer{}
	w.Init(out, 0, 8, 0, '\t', 0)

	formatLine := func(line []string, title bool) (formatted string) {
		for _, cell := range line {
			if title {
				cell = xstrings.Title(cell)
			}
			formatted += fmt.Sprintf("%s \t", cell)
		}
		return formatted
	}

	if len(header) == 0 {
		return errors.Wrap(ErrInvalidFormat, "empty header")
	}

	// write header
	if _, err := fmt.Fprintln(w, formatLine(header, true)); err != nil {
		return err
	}

	// write entries
	for i, entry := range entries {
		if len(entry) != len(header) {
			return errors.Wrapf(ErrInvalidFormat, "entry %d doesn't match header length", i)
		}
		if _, err := fmt.Fprintf(w, formatLine(entry, false)+"\n"); err != nil {
			return err
		}
	}

	if _, err := fmt.Fprintln(w); err != nil {
		return err
	}
	return w.Flush()
}
