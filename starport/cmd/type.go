package starportcmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/clispinner"
	"github.com/tendermint/starport/starport/pkg/placeholder"
	"github.com/tendermint/starport/starport/services/scaffolder"
)

const (
	flagModule    = "module"
	flagIndexed   = "indexed"
	flagNoMessage = "no-message"
)

// NewType command creates a new type command to scaffold types.
func NewType() *cobra.Command {
	c := &cobra.Command{
		Use:   "type [typeName] [field1] [field2] ...",
		Short: "Scaffold a type with CRUD operations",
		Long:  "Scaffold a type with create, read, update and delete operations",
		Args:  cobra.MinimumNArgs(1),
		RunE:  typeHandler,
	}
	c.Flags().StringVarP(&appPath, "path", "p", "", "path of the app")
	c.Flags().String(flagModule, "", "Module to add the type into. Default: app's main module")
	c.Flags().Bool(flagIndexed, false, "Scaffold an indexed type")
	c.Flags().Bool(flagNoMessage, false, "Disable CRUD interaction messages scaffolding")

	return c
}

func typeHandler(cmd *cobra.Command, args []string) error {
	s := clispinner.New().SetText("Scaffolding...")
	defer s.Stop()

	// Get the module to add the type into
	module, _ := cmd.Flags().GetString(flagModule)

	// Add type options
	var opts scaffolder.AddTypeOption
	opts.Indexed, _ = cmd.Flags().GetBool(flagIndexed)
	opts.NoMessage, _ = cmd.Flags().GetBool(flagNoMessage)

	sc, err := scaffolder.New(appPath)
	if err != nil {
		return err
	}
	sm, err := sc.AddType(placeholder.New(), opts, module, args[0], args[1:]...)
	if err != nil {
		return err
	}

	s.Stop()

	fmt.Print(sourceModificationToString(sm))
	fmt.Printf("\n🎉 Created a type `%[1]v`.\n\n", args[0])
	return nil
}
