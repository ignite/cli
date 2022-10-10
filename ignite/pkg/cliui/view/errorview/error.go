package errorview

import (
	"strings"

	"github.com/ignite/cli/ignite/pkg/cliui/colors"
)

func NewError(err error) Error {
	return Error{err}
}

type Error struct {
	Err error
}

func (e Error) String() string {
	b := strings.Builder{}

	b.WriteString(colors.Error(e.Err.Error()))
	b.WriteRune('\n')
	b.WriteString(colors.Info("Waiting for a fix before retrying...\n"))

	return b.String()
}
