package ignitecmd

import (
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"

	"github.com/ignite/cli/ignite/pkg/cliui"
	"github.com/ignite/cli/ignite/services/network"
)

// NewNetworkRewardSet creates a new chain reward set command to
// add the chain reward to the network as a coordinator.
func NewNetworkRewardSet() *cobra.Command {
	c := &cobra.Command{
		Use:   "set [launch-id] [last-reward-height] [coins]",
		Short: "set a network chain reward",
		Args:  cobra.ExactArgs(3),
		RunE:  networkChainRewardSetHandler,
	}
	c.Flags().AddFlagSet(flagSetKeyringBackend())
	c.Flags().AddFlagSet(flagNetworkFrom())
	c.Flags().AddFlagSet(flagSetHome())
	return c
}

func networkChainRewardSetHandler(cmd *cobra.Command, args []string) error {
	session := cliui.New()
	defer session.Cleanup()

	nb, err := newNetworkBuilder(cmd, CollectEvents(session.EventBus()))
	if err != nil {
		return err
	}

	// parse launch ID
	launchID, err := network.ParseID(args[0])
	if err != nil {
		return err
	}

	// parse the last reward height
	lastRewardHeight, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return err
	}

	coins, err := sdk.ParseCoinsNormalized(args[2])
	if err != nil {
		return fmt.Errorf("failed to parse coins: %w", err)
	}

	n, err := nb.Network()
	if err != nil {
		return err
	}

	return n.SetReward(launchID, lastRewardHeight, coins)
}
