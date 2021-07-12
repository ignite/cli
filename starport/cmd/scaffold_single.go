package starportcmd

import (
	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/services/scaffolder"
)

// NewScaffoldSingle returns a new command to scaffold a singleton.
func NewScaffoldSingle() *cobra.Command {
	c := &cobra.Command{
		Use:   "single NAME [field]...",
		Short: "CRUD for data stored in a single location",
		Args:  cobra.MinimumNArgs(1),
		RunE:  scaffoldSingleHandler,
	}

	c.Flags().StringVarP(&appPath, "path", "p", "", "path of the app")
	c.Flags().AddFlagSet(flagSetScaffoldType())

	return c
}

func scaffoldSingleHandler(cmd *cobra.Command, args []string) error {
	return scaffoldType(cmd, args, scaffolder.SingletonType())
}
