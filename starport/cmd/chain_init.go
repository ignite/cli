package starportcmd

import (
	"fmt"

	"github.com/ignite-hq/cli/starport/pkg/chaincmd"
	"github.com/ignite-hq/cli/starport/services/chain"
	"github.com/spf13/cobra"
)

func NewChainInit() *cobra.Command {
	c := &cobra.Command{
		Use:   "init",
		Short: "Initialize your chain",
		Args:  cobra.NoArgs,
		RunE:  chainInitHandler,
	}

	flagSetPath(c)
	c.Flags().AddFlagSet(flagSetHome())

	return c
}

func chainInitHandler(cmd *cobra.Command, args []string) error {
	chainOption := []chain.Option{
		chain.LogLevel(logLevel(cmd)),
		chain.KeyringBackend(chaincmd.KeyringBackendTest),
	}

	c, err := newChainWithHomeFlags(cmd, chainOption...)
	if err != nil {
		return err
	}

	if _, err := c.Build(cmd.Context(), ""); err != nil {
		return err
	}

	if err := c.Init(cmd.Context(), true); err != nil {
		return err
	}

	home, err := c.Home()
	if err != nil {
		return err
	}

	fmt.Printf("🗃  Initialized. Checkout your chain's home (data) directory: %s\n", infoColor(home))

	return nil
}
