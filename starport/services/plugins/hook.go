package plugins

type Hook interface {
	ParentCommand() []string
	Name() string
	Usage() string
	ShortDesc() string
	LongDesc() string

	PreRun() error
	PostRun() error
}
