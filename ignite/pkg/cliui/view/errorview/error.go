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
	return colors.Error(strings.TrimSpace(e.Err.Error()))
}
