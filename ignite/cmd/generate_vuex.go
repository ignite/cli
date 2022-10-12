package ignitecmd

import (
	"github.com/spf13/cobra"

	"github.com/ignite/cli/ignite/pkg/cliui"
	"github.com/ignite/cli/ignite/services/chain"
)

func NewGenerateVuex() *cobra.Command {
	c := &cobra.Command{
		Use:     "vuex",
		Short:   "Generate Typescript client and Vuex stores for your chain's frontend from your `config.yml` file",
		PreRunE: gitChangesConfirmPreRunHandler,
		RunE:    generateVuexHandler,
	}

	c.Flags().AddFlagSet(flagSetProto3rdParty(""))
	c.Flags().AddFlagSet(flagSetYes())

	return c
}

func generateVuexHandler(cmd *cobra.Command, args []string) error {
	session := cliui.New(cliui.StartSpinner())
	defer session.End()

	session.StartSpinner("Generating...")

	c, err := newChainWithHomeFlags(
		cmd,
		chain.EnableThirdPartyModuleCodegen(),
		chain.WithOutputer(session),
		chain.CollectEvents(session.EventBus()),
	)
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

	return session.Println("⛏️  Generated Typescript Client and Vuex stores")
}
