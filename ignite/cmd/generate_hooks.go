package ignitecmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ignite/cli/ignite/pkg/cliui/clispinner"
	"github.com/ignite/cli/ignite/services/chain"
)

func NewGenerateHooks() *cobra.Command {
	c := &cobra.Command{
		Use:   "hooks",
		Short: "Generate Typescript client and React composables for your chain's frontend from your `config.yml` file",
		RunE:  generateHooksHandler,
	}
	c.Flags().AddFlagSet(flagSetProto3rdParty(""))
	return c
}

func generateHooksHandler(cmd *cobra.Command, args []string) error {
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

	if err := c.Generate(cmd.Context(), cacheStorage, chain.GenerateHooks()); err != nil {
		return err
	}

	s.Stop()
	fmt.Println("⛏️  Generated Typescript Client and React hooks")

	return nil
}
