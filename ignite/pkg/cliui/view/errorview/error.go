package errorview

import (
	"strings"

	"github.com/ignite/cli/v29/ignite/pkg/cliui/colors"
	"github.com/muesli/reflow/wordwrap"
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
	_, _ = w.Write([]byte(s))
	_ = w.Close()

	return colors.Error(w.String())
}
