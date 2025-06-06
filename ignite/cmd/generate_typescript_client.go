package ignitecmd

import (
	"github.com/spf13/cobra"

	"github.com/ignite/cli/v29/ignite/pkg/cliui"
	"github.com/ignite/cli/v29/ignite/pkg/cliui/icons"
	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/services/chain"
)

const (
	flagUseCache = "use-cache"
	msgBufAuth   = "Generate ts-client uses a 'buf.build' remote plugin. Buf is begin limiting remote plugin requests from unauthenticated users on 'buf.build'. Intensively using this function will get you rate limited. Authenticate with 'buf registry login' to avoid this (https://buf.build/docs/generate/auth-required)."
)

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
	c.Flags().Bool(flagUseCache, false, "use build cache to speed-up generation")

	return c
}

func generateTSClientHandler(cmd *cobra.Command, _ []string) error {
	session := cliui.New(cliui.StartSpinnerWithText(statusGenerating))
	defer session.End()

	if !getYes(cmd) {
		if err := session.AskConfirm(msgBufAuth); err != nil {
			if errors.Is(err, cliui.ErrAbort) {
				return errors.New("buf not auth")
			}

			return err
		}
	}

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
	useCache, _ := cmd.Flags().GetBool(flagUseCache)

	var opts []chain.GenerateTarget
	if flagGetEnableProtoVendor(cmd) {
		opts = append(opts, chain.GenerateProtoVendor())
	}

	err = c.Generate(cmd.Context(), cacheStorage, chain.GenerateTSClient(output, useCache), opts...)
	if err != nil {
		return err
	}

	return session.Println(icons.OK, "Generated Typescript Client")
}
