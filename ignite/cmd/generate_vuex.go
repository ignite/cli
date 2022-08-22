package ignitecmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ignite/cli/ignite/pkg/cliui/clispinner"
	"github.com/ignite/cli/ignite/services/chain"
)

func NewGenerateVuex() *cobra.Command {
	c := &cobra.Command{
		Use:   "vuex",
		Short: "Generate Vuex store for you chain's frontend from your config.yml",
		RunE:  generateVuexHandler,
	}
	c.Flags().AddFlagSet(flagSetProto3rdParty(""))
	return c
}

func generateVuexHandler(cmd *cobra.Command, args []string) error {
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

	if err := c.Generate(cmd.Context(), cacheStorage, chain.GenerateVuex()); err != nil {
		return err
	}

	s.Stop()
	fmt.Println("⛏️  Generated vuex stores.")

	return nil
}
