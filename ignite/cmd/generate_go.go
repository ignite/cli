package ignitecmd

import (
	"github.com/spf13/cobra"

	"github.com/ignite/cli/v29/ignite/pkg/cliui"
	"github.com/ignite/cli/v29/ignite/pkg/cliui/icons"
	"github.com/ignite/cli/v29/ignite/services/chain"
)

func NewGenerateGo() *cobra.Command {
	c := &cobra.Command{
		Use:   "proto-go",
		Short: "Compile protocol buffer files to Go source code required by Cosmos SDK",
		RunE:  generateGoHandler,
	}

	c.Flags().AddFlagSet(flagSetYes())

	return c
}

func generateGoHandler(cmd *cobra.Command, _ []string) error {
	session := cliui.New(
		cliui.StartSpinnerWithText(statusGenerating),
		cliui.WithoutUserInteraction(getYes(cmd)),
	)
	defer session.End()

	c, err := chain.NewWithHomeFlags(
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

	var opts []chain.GenerateTarget
	if flagGetEnableProtoVendor(cmd) {
		opts = append(opts, chain.GenerateProtoVendor())
	}

	err = c.Generate(cmd.Context(), cacheStorage, chain.GenerateGo(), opts...)
	if err != nil {
		return err
	}

	return session.Println(icons.OK, "Generated Go code")
}
