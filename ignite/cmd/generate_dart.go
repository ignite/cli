package ignitecmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ignite/cli/ignite/pkg/cliui/clispinner"
	"github.com/ignite/cli/ignite/services/chain"
)

func NewGenerateDart() *cobra.Command {
	c := &cobra.Command{
		Use:   "dart",
		Short: "Generate a Dart client",
		RunE:  generateDartHandler,
	}
	return c
}

func generateDartHandler(cmd *cobra.Command, args []string) error {
	s := clispinner.New().SetText("Generating...")
	defer s.Stop()

	c, err := newChainWithHomeFlags(cmd, chain.EnableThirdPartyModuleCodegen())
	if err != nil {
		return err
	}

	cacheStorage, err := newCache(cmd)
	if err != nil {
		return err
	}

	if err := c.Generate(cmd.Context(), cacheStorage, chain.GenerateDart()); err != nil {
		return err
	}

	s.Stop()
	fmt.Println("⛏️  Generated Dart client.")

	return nil
}
