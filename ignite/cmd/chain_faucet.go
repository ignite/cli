package ignitecmd

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"

	"github.com/ignite-hq/cli/ignite/pkg/chaincmd"
	"github.com/ignite-hq/cli/ignite/pkg/cliui"
	"github.com/ignite-hq/cli/ignite/services/chain"
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
		session   = cliui.New(cliui.WithVerbosity(logLevel(cmd)))
	)
	defer session.Cleanup()

	chainOption := []chain.Option{
		chain.KeyringBackend(chaincmd.KeyringBackendTest),
		chain.CollectEvents(session.EventBus()),
		chain.WithLogStreamer(session),
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

	session.StopSpinner()
	return session.Println("📨 Coins sent.")
}
