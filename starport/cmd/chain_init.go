package starportcmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/chaincmd"
	"github.com/tendermint/starport/starport/services/chain"
)

func NewChainInit() *cobra.Command {
	c := &cobra.Command{
		Use:   "init",
		Short: "Initialize your chain",
		Args:  cobra.ExactArgs(0),
		RunE:  chainInitHandler,
	}

	c.Flags().AddFlagSet(flagSetHome())
	c.Flags().StringVarP(&appPath, "path", "p", "", "Path of the app")

	return c
}

func chainInitHandler(cmd *cobra.Command, args []string) error {
	chainOption := []chain.Option{
		chain.LogLevel(logLevel(cmd)),
		chain.KeyringBackend(chaincmd.KeyringBackendTest),
	}

	c, err := newChainWithHomeFlags(cmd, appPath, chainOption...)
	if err != nil {
		return err
	}

	if _, err := c.Build(cmd.Context()); err != nil {
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
