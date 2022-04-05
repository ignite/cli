package starportcmd

import (
	"github.com/ignite-hq/cli/starport/services/scaffolder"
	"github.com/spf13/cobra"
)

// NewScaffoldList returns a new command to scaffold a list.
func NewScaffoldList() *cobra.Command {
	c := &cobra.Command{
		Use:   "list NAME [field]...",
		Short: "CRUD for data stored as an array",
		Args:  cobra.MinimumNArgs(1),
		RunE:  scaffoldListHandler,
	}

	flagSetPath(c)
	c.Flags().AddFlagSet(flagSetScaffoldType())

	return c
}

func scaffoldListHandler(cmd *cobra.Command, args []string) error {
	return scaffoldType(cmd, args, scaffolder.ListType())
}
