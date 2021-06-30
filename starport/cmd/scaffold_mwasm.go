package starportcmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/clispinner"
	"github.com/tendermint/starport/starport/pkg/placeholder"
	"github.com/tendermint/starport/starport/services/scaffolder"
)

func NewScaffoldWasm() *cobra.Command {
	c := &cobra.Command{
		Use:   "wasm",
		Short: "Import a wasm module to app.",
		Long:  "Add support for WebAssembly smart contracts to your blockchain.",
		Args:  cobra.NoArgs,
		RunE:  scaffoldWasmHandler,
	}
	return c
}

func scaffoldWasmHandler(cmd *cobra.Command, args []string) error {
	s := clispinner.New().SetText("Scaffolding...")
	defer s.Stop()

	sc, err := scaffolder.New(appPath)
	if err != nil {
		return err
	}
	sm, err := sc.ImportModule(placeholder.New(), "wasm")
	if err != nil {
		return err
	}

	s.Stop()

	fmt.Println(sourceModificationToString(sm))
	fmt.Printf("\n🎉 Imported wasm.\n\n")
	return nil
}
