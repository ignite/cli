package plugins

import (
	"github.com/spf13/cobra"
)

type Command interface {
	ParentCommand() []string
	Name() string
	Usage() string
	ShortDesc() string
	LongDesc() string
	NumArgs() int
	Exec(*cobra.Command, []string) error
}
