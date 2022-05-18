package errorview

import "github.com/ignite-hq/cli/ignite/pkg/cliui/colors"

type Error struct {
	Err error
}

func NewError(err error) Error {
	return Error{Err: err}
}

func (e Error) String() string {
	return colors.Error(e.Err.Error())
}
