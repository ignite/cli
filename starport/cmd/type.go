package starportcmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/clispinner"
	"github.com/tendermint/starport/starport/services/scaffolder"
)

const (
	flagModule  string = "module"
	flagLegacy  string = "legacy"
	flagIndexed string = "indexed"
	flagNoMessage = "no-message"
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

	c.Flags().String(flagModule, "", "Module to add the type into. Default: app's main module")
	c.Flags().Bool(flagLegacy, false, "Scaffold the type without generating MsgServer service")
	c.Flags().Bool(flagIndexed, false, "Scaffold an indexed type")

	return c
}

func typeHandler(cmd *cobra.Command, args []string) error {
	s := clispinner.New().SetText("Scaffolding...")
	defer s.Stop()

	// Get the module to add the type into
	module, err := cmd.Flags().GetString(flagModule)
	if err != nil {
		return err
	}

	// Add type options
	var opts scaffolder.AddTypeOption
	opts.Legacy, err = cmd.Flags().GetBool(flagLegacy)
	if err != nil {
		return err
	}
	opts.Indexed, err = cmd.Flags().GetBool(flagIndexed)
	if err != nil {
		return err
	}

	sc := scaffolder.New(appPath)
	if err := sc.AddType(opts, module, args[0], args[1:]...); err != nil {
		return err
	}

	s.Stop()

	fmt.Printf("\n🎉 Created a type `%[1]v`.\n\n", args[0])
	return nil
}
