package ignitecmd

import (
	"fmt"

	"github.com/ignite-hq/cli/ignite/pkg/cliui/clispinner"
	"github.com/ignite-hq/cli/ignite/services/chain"
	"github.com/spf13/cobra"
)

func NewGenerateOpenAPI() *cobra.Command {
	return &cobra.Command{
		Use:   "openapi",
		Short: "Generate generates an OpenAPI spec for your chain from your config.yml",
		RunE:  generateOpenAPIHandler,
	}
}

func generateOpenAPIHandler(cmd *cobra.Command, args []string) error {
	s := clispinner.New().SetText("Generating...")
	defer s.Stop()

	c, err := newChainWithHomeFlags(cmd)
	if err != nil {
		return err
	}

	if err := c.Generate(cmd.Context(), chain.GenerateOpenAPI()); err != nil {
		return err
	}

	s.Stop()
	fmt.Println("⛏️  Generated OpenAPI spec.")

	return nil
}
