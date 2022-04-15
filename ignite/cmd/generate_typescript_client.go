package ignitecmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ignite-hq/cli/ignite/pkg/clispinner"
	"github.com/ignite-hq/cli/ignite/services/chain"
)

func NewGenerateTSClient() *cobra.Command {
	c := &cobra.Command{
		Use:   "typescript",
		Short: "Generate Typescript Client for you chain's frontend from your config.yml",
		RunE:  generateTSClientHandler,
	}
	c.Flags().AddFlagSet(flagSetProto3rdParty(""))
	return c
}

func generateTSClientHandler(cmd *cobra.Command, args []string) error {
	s := clispinner.New().SetText("Generating...")
	defer s.Stop()

	c, err := newChainWithHomeFlags(cmd, chain.EnableThirdPartyModuleCodegen())
	if err != nil {
		return err
	}

	if err := c.Generate(cmd.Context(), chain.GenerateTSClient()); err != nil {
		return err
	}

	s.Stop()
	fmt.Println("⛏️  Generated Typescript Client")

	return nil
}
