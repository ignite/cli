package ignitecmd

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"

	"github.com/ignite/cli/ignite/pkg/chaincmd"
	"github.com/ignite/cli/ignite/services/chain"
)

// NewChainFaucet creates a new faucet command to send coins to accounts.
func NewChainFaucet() *cobra.Command {
	c := &cobra.Command{
		Use:   "faucet [address] [coin<,...>]",
		Short: "Send coins to an account",
		Args:  cobra.ExactArgs(2),
		RunE:  chainFaucetHandler,
	}

	flagSetPath(c)
	c.Flags().AddFlagSet(flagSetHome())
	c.Flags().BoolP("verbose", "v", false, "Verbose output")

	return c
}

func chainFaucetHandler(cmd *cobra.Command, args []string) error {
	var (
		toAddress = args[0]
		coins     = args[1]
	)

	chainOption := []chain.Option{
		chain.LogLevel(logLevel(cmd)),
		chain.KeyringBackend(chaincmd.KeyringBackendTest),
	}

	c, err := newChainWithHomeFlags(cmd, chainOption...)
	if err != nil {
		return err
	}

	faucet, err := c.Faucet(cmd.Context())
	if err != nil {
		return err
	}

	// parse provided coins
	parsedCoins, err := sdk.ParseCoinsNormalized(coins)
	if err != nil {
		return err
	}

	// perform transfer from faucet
	if err := faucet.Transfer(cmd.Context(), toAddress, parsedCoins); err != nil {
		return err
	}

	fmt.Println("📨 Coins sent.")
	return nil
}
