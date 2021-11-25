package entrywriter

import (
	"errors"
	"fmt"
	"io"
	"text/tabwriter"
)

// Write writes into out the tabulated entries
func Write(out io.Writer, header []string, entries ...[]string) error {
	w := &tabwriter.Writer{}
	w.Init(out, 0, 8, 0, '\t', 0)

	formatLine := func(line []string) (formatted string) {
		for _, cell := range line {
			formatted += fmt.Sprintf("%s \t", cell)
		}
		return formatted
	}

	if len(header) == 0 {
		return errors.New("empty header")
	}

	// write header
	if _, err := fmt.Fprintln(w, formatLine(header)); err != nil {
		return err
	}

	// write entries
	for i, entry := range entries {
		if len(entry) != len(header) {
			return fmt.Errorf("entry %d doesn't match header length", i)
		}
		if _, err := fmt.Fprintf(w, formatLine(entry)+"\n"); err != nil {
			return err
		}
	}

	if _, err := fmt.Fprintln(w); err != nil {
		return err
	}
	return w.Flush()
}
