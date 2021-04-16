package starportcmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tendermint/starport/starport/pkg/chaincmd"
	"github.com/tendermint/starport/starport/pkg/cosmoscoin"
	"github.com/tendermint/starport/starport/services/chain"
)

// NewFaucet creates a new faucet command to send coins to accounts.
func NewFaucet() *cobra.Command {
	c := &cobra.Command{
		Use:   "faucet [address] [coin<,...>]",
		Short: "Send coins to an account",
		Args:  cobra.ExactArgs(2),
		RunE:  faucetHandler,
	}
	c.Flags().AddFlagSet(flagSetHome())
	c.Flags().StringVarP(&appPath, "path", "p", "", "path of the app")
	c.Flags().BoolP("verbose", "v", false, "Verbose output")
	return c
}

func faucetHandler(cmd *cobra.Command, args []string) error {
	var (
		toAddress = args[0]
		coins     = args[1]
	)

	chainOption := []chain.Option{
		chain.LogLevel(logLevel(cmd)),
		chain.KeyringBackend(chaincmd.KeyringBackendTest),
	}

	c, err := newChainWithHomeFlags(cmd, appPath, chainOption...)
	if err != nil {
		return err
	}

	faucet, err := c.Faucet(cmd.Context())
	if err != nil {
		return err
	}

	for _, coin := range strings.Split(coins, ",") {
		amount, denom, err := cosmoscoin.Parse(coin)
		if err != nil {
			return fmt.Errorf("%s: %s", err, coin)
		}
		if err := faucet.Transfer(cmd.Context(), toAddress, amount, denom); err != nil {
			return err
		}
	}

	fmt.Println("📨 Coins sent.")
	return nil
}
