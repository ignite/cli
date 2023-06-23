package ignitecmd

import (
	"github.com/spf13/cobra"

	"github.com/ignite/cli/ignite/pkg/cliui"
	"github.com/ignite/cli/ignite/pkg/cliui/icons"
	"github.com/ignite/cli/ignite/services/chain"
)

func NewGenerateVuex() *cobra.Command {
	c := &cobra.Command{
		Use:   "vuex",
		Short: "*DEPRECATED* TypeScript frontend client and Vuex stores",
		RunE:  generateVuexHandler,
	}

	c.Flags().AddFlagSet(flagSetYes())
	c.Flags().StringP(flagOutput, "o", "", "Vuex store output path")

	return c
}

func generateVuexHandler(cmd *cobra.Command, _ []string) error {
	session := cliui.New(cliui.StartSpinnerWithText(statusGenerating))
	defer session.End()

	c, err := newChainWithHomeFlags(
		cmd,
		chain.WithOutputer(session),
		chain.CollectEvents(session.EventBus()),
		chain.PrintGeneratedPaths(),
	)
	if err != nil {
		return err
	}

	cacheStorage, err := newCache(cmd)
	if err != nil {
		return err
	}

	output, err := cmd.Flags().GetString(flagOutput)
	if err != nil {
		return err
	}

	if err := c.Generate(cmd.Context(), cacheStorage, chain.GenerateVuex(output)); err != nil {
		return err
	}

	return session.Println(icons.OK, "Generated Typescript Client and Vuex stores")
}
