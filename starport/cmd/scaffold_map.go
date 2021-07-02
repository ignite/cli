package starportcmd

import (
	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/services/scaffolder"
)

const (
	flagModule    = "module"
	flagNoMessage = "no-message"
)

// NewScaffoldMap returns a new command to scaffold a map.
func NewScaffoldMap() *cobra.Command {
	c := &cobra.Command{
		Use:   "map NAME [field]...",
		Short: "Scaffold a map",
		Args:  cobra.MinimumNArgs(1),
		RunE:  scaffoldMapHandler,
	}

	c.Flags().StringVarP(&appPath, "path", "p", "", "path of the app")
	c.Flags().AddFlagSet(flagSetScaffoldType())

	return c
}

func scaffoldMapHandler(cmd *cobra.Command, args []string) error {
	opts := scaffolder.AddTypeOption{
		Indexed:   true,
		NoMessage: flagGetNoMessage(cmd),
	}

	return scaffoldType("map", flagGetModule(cmd), args[0], args[1:], opts)
}
