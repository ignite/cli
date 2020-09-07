package starportcmd

import (
	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/services/scaffolder"
)

func NewModuleImport() *cobra.Command {
	c := &cobra.Command{
		Use:   "import [feature]",
		Short: "Imports a new module to app.",
		Args:  cobra.MinimumNArgs(1),
		RunE:  importModuleHandler,
	}
	return c
}

func importModuleHandler(cmd *cobra.Command, args []string) error {
	sc := scaffolder.New(appPath)
	return sc.AddModule(args[0])
}
