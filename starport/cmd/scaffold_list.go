package starportcmd

import (
	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/services/scaffolder"
)

// NewScaffoldList returns a new command to scaffold a list.
func NewScaffoldList() *cobra.Command {
	c := &cobra.Command{
		Use:   "list NAME [field]...",
		Short: "Scaffold a list",
		Args:  cobra.MinimumNArgs(1),
		RunE:  scaffoldListHandler,
	}

	c.Flags().StringVarP(&appPath, "path", "p", "", "path of the app")
	c.Flags().AddFlagSet(flagSetScaffoldType())

	return c
}

func scaffoldListHandler(cmd *cobra.Command, args []string) error {
	opts := scaffolder.AddTypeOption{
		NoMessage: flagGetNoMessage(cmd),
	}

	return scaffoldType("list", flagGetModule(cmd), args[0], args[1:], opts)
}
