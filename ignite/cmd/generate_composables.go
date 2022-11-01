package ignitecmd

import (
	"github.com/spf13/cobra"

	"github.com/ignite/cli/ignite/pkg/cliui"
	"github.com/ignite/cli/ignite/pkg/cliui/icons"
	"github.com/ignite/cli/ignite/services/chain"
)

func NewGenerateComposables() *cobra.Command {
	c := &cobra.Command{
		Use:   "composables",
		Short: "Generate Typescript client and Vue 3 composables for your chain's frontend from your `config.yml` file",
		RunE:  generateComposablesHandler,
	}
	c.Flags().AddFlagSet(flagSetProto3rdParty(""))
	return c
}

func generateComposablesHandler(cmd *cobra.Command, args []string) error {
	session := cliui.New(cliui.StartSpinnerWithText(statusGenerating))
	defer session.End()

	c, err := NewChainWithHomeFlags(cmd, chain.EnableThirdPartyModuleCodegen(),
		chain.WithOutputer(session),
		chain.CollectEvents(session.EventBus()),
		chain.PrintGeneratedPaths())
	if err != nil {
		return err
	}

	cacheStorage, err := newCache(cmd)
	if err != nil {
		return err
	}

	if err := c.Generate(cmd.Context(), cacheStorage, chain.GenerateComposables()); err != nil {
		return err
	}

	return session.Println(icons.OK, "Generated Typescript Client and Vue 3 composables")
}
