package ignitecmd

import (
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"

<<<<<<< HEAD
	"github.com/ignite/cli/v28/ignite/pkg/cliui"
	"github.com/ignite/cli/v28/ignite/pkg/cliui/icons"
	"github.com/ignite/cli/v28/ignite/services/chain"
=======
	"github.com/ignite/cli/v29/ignite/pkg/cliui"
	"github.com/ignite/cli/v29/ignite/pkg/cliui/icons"
	"github.com/ignite/cli/v29/ignite/pkg/errors"
	"github.com/ignite/cli/v29/ignite/services/chain"
>>>>>>> 7ff4b5d1 (feat: create a message for authenticate buf for generate ts-client (#4322))
)

const (
	flagUseCache = "use-cache"
	msgBufAuth   = "Generate ts-client depends on a 'buf.build' remote plugin, and as of August 1, 2024, Buf will begin limiting remote plugin requests from unauthenticated users on 'buf.build'. If you send more than ten unauthenticated requests per hour using remote plugins, you’ll start to see rate limit errors. Please authenticate before running ts-client command using 'buf registry login' command and follow the instructions. For more info, check https://buf.build/docs/generate/auth-required."
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

	if err := session.AskConfirm(msgBufAuth); err != nil {
		if errors.Is(err, promptui.ErrAbort) {
			return errors.New("buf not auth")
		}
		return err
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
