package ignitecmd

import (
	"github.com/spf13/cobra"

	"github.com/ignite/cli/ignite/pkg/cliui"
	"github.com/ignite/cli/ignite/pkg/cliui/icons"
	"github.com/ignite/cli/ignite/services/chain"
)

func NewGenerateGo() *cobra.Command {
	c := &cobra.Command{
		Use:     "proto-go",
		Short:   "Compile protocol buffer files to Go source code required by Cosmos SDK",
		PreRunE: gitChangesConfirmPreRunHandler,
		RunE:    generateGoHandler,
	}

	c.Flags().AddFlagSet(flagSetYes())

	return c
}

func generateGoHandler(cmd *cobra.Command, _ []string) error {
	session := cliui.New(cliui.StartSpinnerWithText(statusGenerating))
	defer session.End()

	c, err := newChainWithHomeFlags(
		cmd,
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

	if err := c.Generate(cmd.Context(), cacheStorage, chain.GenerateGo()); err != nil {
		return err
	}

	return session.Println(icons.OK, "Generated Go code")
}
