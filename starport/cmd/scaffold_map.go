package starportcmd

import (
	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/services/scaffolder"
)

const (
	FlagIndexes = "index"
)

// NewScaffoldMap returns a new command to scaffold a map.
func NewScaffoldMap() *cobra.Command {
	c := &cobra.Command{
		Use:   "map NAME [field]...",
		Short: "CRUD for data stored as key-value pairs",
		Args:  cobra.MinimumNArgs(1),
		RunE:  scaffoldMapHandler,
	}

	c.Flags().StringVarP(&appPath, "path", "p", "", "path of the app")
	c.Flags().AddFlagSet(flagSetScaffoldType())
	c.Flags().StringSlice(FlagIndexes, []string{"index"}, "fields that index the value")

	return c
}

func scaffoldMapHandler(cmd *cobra.Command, args []string) error {
	indexes, err := cmd.Flags().GetStringSlice(FlagIndexes)
	if err != nil {
		return err
	}

	return scaffoldType(cmd, args, scaffolder.MapType(indexes...))
}
