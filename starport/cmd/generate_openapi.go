package starportcmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/clispinner"
	"github.com/tendermint/starport/starport/services/chain"
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

	chainOption := []chain.Option{}

	c, err := newChainWithHomeFlags(cmd, appPath, chainOption...)
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
