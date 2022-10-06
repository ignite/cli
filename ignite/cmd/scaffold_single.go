package ignitecmd

import (
	"github.com/spf13/cobra"

	"github.com/ignite/cli/ignite/services/scaffolder"
)

// NewScaffoldSingle returns a new command to scaffold a singleton.
func NewScaffoldSingle() *cobra.Command {
	c := &cobra.Command{
		Use:     "single NAME [field]...",
		Short:   "CRUD for data stored in a single location",
		Args:    cobra.MinimumNArgs(1),
		PreRunE: gitChangesConfirmPreRunHandler,
		RunE:    scaffoldSingleHandler,
	}

	flagSetPath(c)
	flagSetClearCache(c)

	c.Flags().AddFlagSet(flagSetYes())
	c.Flags().AddFlagSet(flagSetScaffoldType())

	return c
}

func scaffoldSingleHandler(cmd *cobra.Command, args []string) error {
	return scaffoldType(cmd, args, scaffolder.SingletonType())
}
