package starportcmd

import (
	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/services/scaffolder"
)

// NewScaffoldIndexedList returns a new command to scaffold an indexed list.
func NewScaffoldIndexedList() *cobra.Command {
	c := &cobra.Command{
		Use:   "indexed-list NAME [field]...",
		Short: "CRUD for data stored as an indexed array",
		Args:  cobra.MinimumNArgs(1),
		RunE:  scaffoldIndexedListHandler,
	}

	flagSetPath(c)
	c.Flags().AddFlagSet(flagSetScaffoldType())
	c.Flags().StringSlice(FlagIndexes, []string{"index"}, "fields that index the list")

	return c
}

func scaffoldIndexedListHandler(cmd *cobra.Command, args []string) error {
	indexes, err := cmd.Flags().GetStringSlice(FlagIndexes)
	if err != nil {
		return err
	}

	return scaffoldType(cmd, args, scaffolder.IndexedListType(indexes...))
}
