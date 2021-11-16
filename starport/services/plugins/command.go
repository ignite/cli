package plugins

import "context"

type Command interface {
	ParentCommand() []string
	Name() string
	Usage() string
	ShortDesc() string
	LongDesc() string
	NumArgs() int
	Exec(*cobra.Command, []string) error
}

func ValidateParentCommand(parentCommand []string) {

}
