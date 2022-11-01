package ignitecmd

import (
	"errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"

	"github.com/ignite/cli/ignite/pkg/cliui"
	"github.com/ignite/cli/ignite/pkg/cosmosutil"
	"github.com/ignite/cli/ignite/services/network"
	"github.com/ignite/cli/ignite/services/network/networkchain"
	"github.com/ignite/cli/ignite/services/network/networktypes"
)

// NewNetworkRequestAddAccount creates a new command to send add account request
func NewNetworkRequestAddAccount() *cobra.Command {
	c := &cobra.Command{
		Use:   "add-account [launch-id] [address] [coins]",
		Short: "Send request to add account",
		RunE:  networkRequestAddAccountHandler,
		Args:  cobra.RangeArgs(2, 3),
	}

	flagSetClearCache(c)
	c.Flags().AddFlagSet(flagNetworkFrom())
	c.Flags().AddFlagSet(flagSetHome())
	c.Flags().AddFlagSet(flagSetKeyringBackend())
	c.Flags().AddFlagSet(flagSetKeyringDir())
	return c
}

func networkRequestAddAccountHandler(cmd *cobra.Command, args []string) error {
	session := cliui.New(cliui.StartSpinner())
	defer session.End()

	nb, err := newNetworkBuilder(cmd, CollectEvents(session.EventBus()))
	if err != nil {
		return err
	}

	// parse launch ID
	launchID, err := network.ParseID(args[0])
	if err != nil {
		return err
	}

	// get the address for the account and change the prefix for Ignite Chain
	address, err := cosmosutil.ChangeAddressPrefix(args[1], networktypes.SPN)
	if err != nil {
		return err
	}

	n, err := nb.Network()
	if err != nil {
		return err
	}

	chainLaunch, err := n.ChainLaunch(cmd.Context(), launchID)
	if err != nil {
		return err
	}

	c, err := nb.Chain(networkchain.SourceLaunch(chainLaunch))
	if err != nil {
		return err
	}

	var balance sdk.Coins
	if c.IsAccountBalanceFixed() {
		balance = c.AccountBalance()
		if len(args) == 3 {
			return fmt.Errorf(
				"balance can't be provided, balance has been set by coordinator to %s",
				balance.String(),
			)
		}
	} else {
		if len(args) < 3 {
			return errors.New("account balance expected")
		}
		balanceStr := args[2]
		balance, err = sdk.ParseCoinsNormalized(balanceStr)
		if err != nil {
			return err
		}
	}

	return n.SendAccountRequest(
		cmd.Context(),
		launchID,
		address,
		balance,
	)
}
