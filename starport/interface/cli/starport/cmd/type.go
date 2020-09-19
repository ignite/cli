package starportcmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/services/scaffolder"
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
	return c
}

func typeHandler(cmd *cobra.Command, args []string) error {
	sc := scaffolder.New(appPath)
	if err := sc.AddType(args[0], args[1:]...); err != nil {
		return err
	}
	fmt.Printf("\n🎉 Created a type `%[1]v`.\n\n", args[0])
	return nil
}
