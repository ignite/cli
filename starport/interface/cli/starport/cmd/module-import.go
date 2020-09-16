package starportcmd

import (
	"fmt"

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
	name := args[0]
	sc := scaffolder.New(appPath)
	if err := sc.AddModule(name); err != nil {
		return err
	}
	fmt.Printf("\n🎉 Imported module `%s`.\n\n", name)
	return nil
}
