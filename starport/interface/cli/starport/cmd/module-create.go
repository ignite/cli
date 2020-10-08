package starportcmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/services/scaffolder"
)

func NewModuleCreate() *cobra.Command {
	c := &cobra.Command{
		Use:   "create [name]",
		Short: "Creates a new empty module to app.",
		Long:  "Use starport module create to create a new empty module to your blockchain.",
		Args:  cobra.MinimumNArgs(1),
		RunE:  createModuleHandler,
	}
	return c
}

func createModuleHandler(cmd *cobra.Command, args []string) error {
	name := args[0]
	sc := scaffolder.New(appPath)
	if err := sc.CreateModule(name); err != nil {
		return err
	}
	fmt.Printf("\n🎉 Module created %s.\n\n", name)
	return nil
}
