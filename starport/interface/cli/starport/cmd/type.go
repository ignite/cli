package starportcmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/clispinner"
	"github.com/tendermint/starport/starport/services/scaffolder"
)

const (
	moduleFlag string = "module"
	legacyFlag string = "legacy"
)

// NewType command creates a new type command to scaffold types.
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
	c.Flags().Bool(legacyFlag, false, "Scaffold the type without generating MsgServer service")

	return c
}

func typeHandler(cmd *cobra.Command, args []string) error {
	s := clispinner.New().SetText("Scaffolding...")
	defer s.Stop()

	// Get the module to add the type into
	module, err := cmd.Flags().GetString(moduleFlag)
	if err != nil {
		return err
	}
	legacy, err := cmd.Flags().GetBool(legacyFlag)
	if err != nil {
		return err
	}

	sc := scaffolder.New(appPath)
	if err := sc.AddType(legacy, module, args[0], args[1:]...); err != nil {
		return err
	}

	s.Stop()

	fmt.Printf("\n🎉 Created a type `%[1]v`.\n\n", args[0])
	return nil
}
