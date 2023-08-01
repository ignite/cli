package ignitecmd

import (
	"github.com/spf13/cobra"

	"github.com/ignite/cli/ignite/pkg/cliui"
	"github.com/ignite/cli/ignite/pkg/placeholder"
	"github.com/ignite/cli/ignite/services/scaffolder"
)

func NewScaffoldWasm() *cobra.Command {
	c := &cobra.Command{
		Use:     "wasm",
		Short:   "Import the wasm module to your app",
		Long:    "Add support for WebAssembly smart contracts to your blockchain",
		Args:    cobra.NoArgs,
		RunE:    scaffoldWasmHandler,
		PreRunE: migrationPreRunHandler,
	}

	flagSetPath(c)

	return c
}

func scaffoldWasmHandler(cmd *cobra.Command, _ []string) error {
	appPath := flagGetPath(cmd)

	session := cliui.New(cliui.StartSpinnerWithText(statusScaffolding))
	defer session.End()

	cacheStorage, err := newCache(cmd)
	if err != nil {
		return err
	}

	sc, err := scaffolder.New(appPath)
	if err != nil {
		return err
	}

	sm, err := sc.ImportModule(cmd.Context(), cacheStorage, placeholder.New(), "wasm")
	if err != nil {
		return err
	}

	modificationsStr, err := sourceModificationToString(sm)
	if err != nil {
		return err
	}

	session.Println(modificationsStr)
	session.Printf("\n🎉 Imported wasm.\n\n")

	return nil
}
