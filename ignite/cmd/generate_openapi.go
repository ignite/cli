package ignitecmd

import (
	"github.com/spf13/cobra"

	"github.com/ignite/cli/v29/ignite/pkg/cliui"
	"github.com/ignite/cli/v29/ignite/pkg/cliui/icons"
	"github.com/ignite/cli/v29/ignite/services/chain"
)

var excludeFlag = "exclude"

func NewGenerateOpenAPI() *cobra.Command {
	c := &cobra.Command{
		Use:   "openapi",
		Short: "OpenAPI spec for your chain",
		RunE:  generateOpenAPIHandler,
	}

	c.Flags().AddFlagSet(flagSetYes())
	c.Flags().StringSlice(excludeFlag, []string{}, "List of proto files or directories to exclude from the OpenAPI spec generation")

	return c
}

func generateOpenAPIHandler(cmd *cobra.Command, _ []string) error {
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

	var opts []chain.GenerateTarget
	if flagGetEnableProtoVendor(cmd) {
		opts = append(opts, chain.GenerateProtoVendor())
	}

	excludeList, _ := cmd.Flags().GetStringArray(excludeFlag)

	err = c.Generate(cmd.Context(), cacheStorage, chain.GenerateOpenAPI(excludeList), opts...)
	if err != nil {
		return err
	}

	return session.Println(icons.OK, "Generated OpenAPI spec")
}
