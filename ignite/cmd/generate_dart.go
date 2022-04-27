package ignitecmd

import (
	"fmt"

	"github.com/ignite-hq/cli/ignite/pkg/cliui/clispinner"
	"github.com/ignite-hq/cli/ignite/services/chain"
	"github.com/spf13/cobra"
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

	if err := c.Generate(cmd.Context(), chain.GenerateDart()); err != nil {
		return err
	}

	s.Stop()
	fmt.Println("⛏️  Generated Dart client.")

	return nil
}
