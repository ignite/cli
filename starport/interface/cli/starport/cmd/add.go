package starportcmd

import (
	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/services/scaffolder"
)

func NewAdd() *cobra.Command {
	c := &cobra.Command{
		Use:   "add [feature]",
		Short: "Adds a feature to a project.",
		Args:  cobra.MinimumNArgs(1),
		RunE:  addHandler,
	}
	return c
}

func addHandler(cmd *cobra.Command, args []string) error {
	sc := scaffolder.New(appPath)
	return sc.AddModule(args[0])
}
