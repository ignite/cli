package ignitecmd

import (
	"github.com/spf13/cobra"

	"github.com/ignite/cli/ignite/pkg/cliui"
	"github.com/ignite/cli/ignite/pkg/cliui/icons"
	"github.com/ignite/cli/ignite/services/chain"
)

func NewGeneratePulsar() *cobra.Command {
	c := &cobra.Command{
		Use:   "proto-pulsar",
		Short: "Compile protocol buffer files to Go pulsar source code required by Cosmos SDK",
		RunE:  generatePulsarHandler,
	}

	c.Flags().AddFlagSet(flagSetYes())

	return c
}

func generatePulsarHandler(cmd *cobra.Command, _ []string) error {
	session := cliui.New(cliui.StartSpinnerWithText(statusGenerating))
	defer session.End()

	c, err := newChainWithHomeFlags(
		cmd,
		chain.WithOutputer(session),
		chain.CollectEvents(session.EventBus()),
		chain.CheckCosmosSDKVersion(),
	)
	if err != nil {
		return err
	}

	cacheStorage, err := newCache(cmd)
	if err != nil {
		return err
	}

	if err := c.Generate(cmd.Context(), cacheStorage, chain.GeneratePulsar()); err != nil {
		return err
	}

	return session.Println(icons.OK, "Generated Go pulsar code")
}
