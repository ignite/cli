package ignitecmd

import (
	"github.com/spf13/cobra"

	"github.com/ignite/cli/v29/ignite/pkg/cliui"
	"github.com/ignite/cli/v29/ignite/pkg/cliui/icons"
	"github.com/ignite/cli/v29/ignite/services/chain"
)

func NewGenerateComposables() *cobra.Command {
	c := &cobra.Command{
		Use:   "composables",
		Short: "TypeScript frontend client and Vue 3 composables",
		RunE:  generateComposablesHandler,
	}

	c.Flags().AddFlagSet(flagSetYes())
	c.Flags().StringP(flagOutput, "o", "", "Vue 3 composables output path")

	return c
}

func generateComposablesHandler(cmd *cobra.Command, _ []string) error {
	session := cliui.New(
		cliui.StartSpinnerWithText(statusGenerating),
		cliui.WithoutUserInteraction(getYes(cmd)),
	)
	defer session.End()

	c, err := chain.NewWithHomeFlags(
		cmd,
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

	output, _ := cmd.Flags().GetString(flagOutput)

	var opts []chain.GenerateTarget
	if flagGetEnableProtoVendor(cmd) {
		opts = append(opts, chain.GenerateProtoVendor())
	}

	err = c.Generate(cmd.Context(), cacheStorage, chain.GenerateComposables(output), opts...)
	if err != nil {
		return err
	}

	return session.Println(icons.OK, "Generated Typescript Client and Vue 3 composables")
}
