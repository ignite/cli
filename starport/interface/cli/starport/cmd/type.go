package starportcmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/services/scaffolder"
)

const (
	moduleFlag string = "module"
)

func NewType() *cobra.Command {
	c := &cobra.Command{
		Use:   "type [typeName] [field1] [field2] ...",
		Short: "Generates CRUD actions for type",
		Args:  cobra.MinimumNArgs(1),
		RunE:  typeHandler,
	}
	c.Flags().StringVarP(&appPath, "path", "p", "", "path of the app")
	addSdkVersionFlag(c)

	c.Flags().String(moduleFlag, "", "Module to add the type into. Default: app's main module")

	return c
}

func typeHandler(cmd *cobra.Command, args []string) error {
	// Get the module to add the type into
	module, _ := cmd.Flags().GetString(moduleFlag)

	sc := scaffolder.New(appPath)
	if err := sc.AddType(module, args[0], args[1:]...); err != nil {
		return err
	}
	fmt.Printf("\n🎉 Created a type `%[1]v`.\n\n", args[0])
	return nil
}
