package entrywriter

import (
	"errors"
	"fmt"
	"io"
	"os"
	"text/tabwriter"
)

type Tab string

type Entry []string

// Write writes into out the tabulated entries
func Write(out io.Writer, tabs []Tab, entries ...Entry) error {
	w := &tabwriter.Writer{}
	w.Init(out, 0, 8, 0, '\t', 0)

	formatLine := func (line []string) (formatted string) {
		for _, cell := range line {
			formatted += fmt.Sprintf("%s\t", cell)
		}
		return formatted + "\n"
	}

	if len(tabs) == 0 {
		return errors.New("no tab")
	}

	formatLine(tabs)

	fmt.Fprintln(w, "name\taddress\tpublic key")

	for _, acc := range accounts {
		fmt.Fprintf(w, "%s\t%s\t%s\n",
			acc.Name,
			acc.Address(getAddressPrefix(cmd)),
			acc.PubKey(),
		)
	}

	fmt.Fprintln(w)
	w.Flush()

	return nil
}