package ignitecmd

import (
	"github.com/spf13/cobra"

	"github.com/ignite/cli/ignite/pkg/cliui"
	"github.com/ignite/cli/ignite/pkg/cliui/icons"
	"github.com/ignite/cli/ignite/services/chain"
)

func NewGenerateTSClient() *cobra.Command {
	c := &cobra.Command{
		Use:     "ts-client",
		Short:   "Generate Typescript client for your chain's frontend",
		PreRunE: gitChangesConfirmPreRunHandler,
		RunE:    generateTSClientHandler,
	}

	c.Flags().AddFlagSet(flagSetYes())
	c.Flags().StringP(flagOutput, "o", "", "typescript client output path")

	return c
}

func generateTSClientHandler(cmd *cobra.Command, _ []string) error {
	session := cliui.New(cliui.StartSpinnerWithText(statusGenerating))
	defer session.End()

	c, err := NewChainWithHomeFlags(
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

	err = c.Generate(cmd.Context(), cacheStorage, chain.GenerateTSClient(output))
	if err != nil {
		return err
	}

	return session.Println(icons.OK, "Generated Typescript Client")
}
