package ignitecmd

import (
	"github.com/spf13/cobra"

	"github.com/ignite/cli/v29/ignite/pkg/cliui"
	"github.com/ignite/cli/v29/ignite/pkg/cliui/icons"
	"github.com/ignite/cli/v29/ignite/services/chain"
)

const flagDisableCache = "disable-cache"

func NewGenerateTSClient() *cobra.Command {
	c := &cobra.Command{
		Use:   "ts-client",
		Short: "TypeScript frontend client",
		Long: `Generate a framework agnostic TypeScript client for your blockchain project.

By default the TypeScript client is generated in the "ts-client/" directory. You
can customize the output directory in config.yml:

	client:
	  typescript:
	    path: new-path

Output can also be customized by using a flag:

	ignite generate ts-client --output new-path

TypeScript client code can be automatically regenerated on reset or source code
changes when the blockchain is started with a flag:

	ignite chain serve --generate-clients
`,
		RunE: generateTSClientHandler,
	}

	c.Flags().AddFlagSet(flagSetYes())
	c.Flags().StringP(flagOutput, "o", "", "TypeScript client output path")
	c.Flags().Bool(flagDisableCache, false, "disable build cache")

	return c
}

func generateTSClientHandler(cmd *cobra.Command, _ []string) error {
	session := cliui.New(
		cliui.StartSpinnerWithText(statusGenerating),
		cliui.WithoutUserInteraction(getYes(cmd)),
	)
	defer session.End()

	c, err := chain.NewWithHomeFlags(
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

	output, _ := cmd.Flags().GetString(flagOutput)
	disableCache, _ := cmd.Flags().GetBool(flagDisableCache)

	var opts []chain.GenerateTarget
	if flagGetEnableProtoVendor(cmd) {
		opts = append(opts, chain.GenerateProtoVendor())
	}

	err = c.Generate(cmd.Context(), cacheStorage, chain.GenerateTSClient(output, !disableCache), opts...)
	if err != nil {
		return err
	}

	return session.Println(icons.OK, "Generated Typescript Client")
}
