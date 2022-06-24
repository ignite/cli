package ignitecmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ignite/cli/ignite/pkg/cliui/clispinner"
	"github.com/ignite/cli/ignite/pkg/placeholder"
)

func NewScaffoldWasm() *cobra.Command {
	c := &cobra.Command{
		Use:   "wasm",
		Short: "Import the wasm module to your app",
		Long:  "Add support for WebAssembly smart contracts to your blockchain",
		Args:  cobra.NoArgs,
		RunE:  scaffoldWasmHandler,
	}

	flagSetPath(c)

	return c
}

func scaffoldWasmHandler(cmd *cobra.Command, args []string) error {
	appPath := flagGetPath(cmd)

	s := clispinner.New().SetText("Scaffolding...")
	defer s.Stop()

	cacheStorage, err := newCache(cmd)
	if err != nil {
		return err
	}

	sc, err := newApp(appPath)
	if err != nil {
		return err
	}

	sm, err := sc.ImportModule(cacheStorage, placeholder.New(), "wasm")
	if err != nil {
		return err
	}

	s.Stop()

	modificationsStr, err := sourceModificationToString(sm)
	if err != nil {
		return err
	}

	fmt.Println(modificationsStr)
	fmt.Printf("\n🎉 Imported wasm.\n\n")

	return nil
}
