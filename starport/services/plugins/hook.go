package plugins

import "github.com/spf13/cobra"

type Hook interface {
	ParentCommand() []string
	Name() string
	Type() string
	ShortDesc() string

	PreRun(cmd *cobra.Command, args []string) error
	PostRun(cmd *cobra.Command, args []string) error
}
