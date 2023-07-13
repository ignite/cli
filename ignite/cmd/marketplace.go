package ignitecmd

import (
	"context"
	"os"
	"time"

	"github.com/google/go-github/github"
	"github.com/ignite/cli/ignite/pkg/cliui"
	mpService "github.com/ignite/cli/ignite/services/marketplace"
	"github.com/spf13/cobra"
)

func NewMarketplace() *cobra.Command {
	c := &cobra.Command{
		Use:   "marketplace [command] [args]",
		Short: "Installing plugins from marketplace",
		Long: `Instal a plugin from GitHub.

Ignite CLI has a feature called plugins. It allows you
to extend the functionality of the CLI without having
to touch the core codebase. This command allows you to
install a plugin from GitHub.
    `,
		Aliases: []string{"m"},
		Args:    cobra.ExactArgs(1),
	}

	flagSetPath(c)
	flagSetClearCache(c)
	c.AddCommand(NewMarketplaceList())
	c.AddCommand(NewMarketplaceInfo())
	c.AddCommand(NewMarketplaceAdd())

	return c
}

func NewMarketplaceAdd() *cobra.Command {
	c := &cobra.Command{
		Use:   "add",
		Short: "For adding plugins from marketplace",
		Args:  cobra.ExactArgs(1),
		RunE:  marketplaceAddHandler,
	}

	c.Flags().AddFlagSet(flagSetYes())

	return c
}

func NewMarketplaceInfo() *cobra.Command {
	c := &cobra.Command{
		Use:   "info",
		Short: "For getting info about plugins from marketplace",
		RunE:  marketplaceInfoHandler,
		Args:  cobra.ExactArgs(1),
	}

	c.Flags().AddFlagSet(flagSetYes())

	return c
}

func NewMarketplaceList() *cobra.Command {
	c := &cobra.Command{
		Use:   "list",
		Short: "For listing plugins from marketplace",
		RunE:  listMarketplaceHandler,
	}

	return c
}

func listMarketplaceHandler(_ *cobra.Command, _ []string) error {
	session := cliui.New(cliui.StartSpinnerWithText(statusQuerying))
	defer session.End()

	var (
		ctx              = context.Background()
		c                = mpService.NewClient(ctx, getAccessToken())
		queryCtx, cancel = context.WithTimeout(ctx, 5*time.Second)

		opts = &github.SearchOptions{
			Sort:  "stars",
			Order: "desc",
		}
	)

	defer cancel()

	return mpService.ListPlugins(queryCtx, c, opts)
}

func marketplaceInfoHandler(_ *cobra.Command, args []string) error {
	session := cliui.New(cliui.StartSpinnerWithText(statusQuerying))
	defer session.End()

	var (
		ctx              = context.Background()
		c                = mpService.NewClient(ctx, getAccessToken())
		queryCtx, cancel = context.WithTimeout(ctx, 5*time.Second)
	)

	defer cancel()

	return mpService.InfoPlugin(queryCtx, c, args[0])
}

func marketplaceAddHandler(_ *cobra.Command, args []string) error {
	session := cliui.New(cliui.WithStdout(os.Stdout))
	defer session.End()

	return mpService.AddPlugin(args[0])
}

func getAccessToken() string {
	return os.Getenv("GITHUB_TOKEN")
}
