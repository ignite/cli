package ignitecmd

import (
	"github.com/spf13/cobra"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ignite/cli/v29/ignite/pkg/chaincmd"
	"github.com/ignite/cli/v29/ignite/pkg/cliui"
	"github.com/ignite/cli/v29/ignite/services/chain"
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
	c.Flags().BoolP("verbose", "v", false, "verbose output")

	return c
}

func chainFaucetHandler(cmd *cobra.Command, args []string) error {
	var (
		toAddress = args[0]
		coins     = args[1]
		session   = cliui.New(cliui.StartSpinner())
	)
	defer session.End()

	chainOption := []chain.Option{
		chain.KeyringBackend(chaincmd.KeyringBackendTest),
		chain.WithOutputer(session),
		chain.CollectEvents(session.EventBus()),
	}

	config, _ := cmd.Flags().GetString(flagConfig)
	if config != "" {
		chainOption = append(chainOption, chain.ConfigFile(config))
	}

	c, err := chain.NewWithHomeFlags(cmd, chainOption...)
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
	hash, err := faucet.Transfer(cmd.Context(), toAddress, parsedCoins)
	if err != nil {
		return err
	}

	_ = session.Println("📨 Coins sent.")
	return session.Printf("Transaction Hash: %s\n", hash)
}
