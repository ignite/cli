package errorview

import (
	"strings"

	"github.com/muesli/reflow/wordwrap"

	"github.com/ignite/cli/ignite/pkg/cliui/colors"
)

func NewError(err error) Error {
	return Error{err}
}

type Error struct {
	Err error
}

func (e Error) String() string {
	s := strings.TrimSpace(e.Err.Error())

	w := wordwrap.NewWriter(80)
	w.Breakpoints = []rune{' '}
	w.Write([]byte(s))
	w.Close()

	return colors.Error(w.String())
}
