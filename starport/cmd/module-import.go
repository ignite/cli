package starportcmd

import (
	"context"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/clispinner"
	"github.com/tendermint/starport/starport/pkg/validation"
	"github.com/tendermint/starport/starport/services/scaffolder"
)

// NewModuleImport creates a new command to import an sdk module.
func NewModuleImport() *cobra.Command {
	c := &cobra.Command{
		Use:       "import [feature]",
		Short:     "Import a new module to app.",
		Long:      "Add support for WebAssembly smart contracts to your blockchain.",
		Args:      cobra.MinimumNArgs(1),
		ValidArgs: []string{"wasm"},
		RunE:      importModuleHandler,
	}
	RegisterValidationFlags(c)
	return c
}

func importModuleHandler(cmd *cobra.Command, args []string) error {
	s := clispinner.New().SetText("Scaffolding...")
	defer s.Stop()

	name := args[0]
	sc, err := scaffolder.New(appPath)
	if err != nil {
		return err
	}
	ctx := WithValidation(context.Background(), cmd)
	if err := sc.ImportModule(ctx, name); err != nil {
		var valerr validation.Error
		if errors.As(err, &valerr) {
			return errors.New(valerr.ValidationInfo())
		}
		return err
	}

	s.Stop()

	fmt.Printf("\n🎉 Imported module `%s`.\n\n", name)
	return nil
}
