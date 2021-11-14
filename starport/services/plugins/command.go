package plugins

import "context"

type Command interface {
	ParentCommand() []string
	Name() string
	Usage() string
	ShortDesc() string
	LongDesc() string
	Exec(context.Context, []string) (context.Context, error)
}
