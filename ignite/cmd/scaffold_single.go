package ignitecmd

import (
	"github.com/spf13/cobra"

	"github.com/ignite/cli/v29/ignite/services/scaffolder"
)

// NewScaffoldSingle returns a new command to scaffold a singleton.
func NewScaffoldSingle() *cobra.Command {
	c := &cobra.Command{
		Use:   "single NAME [field:type]...",
		Short: "CRUD for data stored in a single location",
		Long: `CRUD for data stored in a single location.
		
For detailed type information use ignite scaffold type --help.`,
		Example: "  ignite scaffold single todo-single title:string done:bool",
		Args:    cobra.MinimumNArgs(1),
		PreRunE: migrationPreRunHandler,
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
