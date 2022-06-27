package network

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	rewardtypes "github.com/tendermint/spn/x/reward/types"

	"github.com/ignite/cli/ignite/pkg/cliui/icons"
	"github.com/ignite/cli/ignite/pkg/events"
	"github.com/ignite/cli/ignite/services/network/networktypes"
)

// SetReward set a chain reward
func (n Network) SetReward(launchID uint64, lastRewardHeight int64, coins sdk.Coins) error {
	n.ev.Send(events.New(
		events.StatusOngoing,
		fmt.Sprintf("Setting reward %s to the chain %d at height %d",
			coins.String(),
			launchID,
			lastRewardHeight,
		),
	))

	msg := rewardtypes.NewMsgSetRewards(
		n.account.Address(networktypes.SPN),
		launchID,
		lastRewardHeight,
		coins,
	)
	res, err := n.cosmos.BroadcastTx(n.account.Name, msg)
	if err != nil {
		return err
	}

	var setRewardRes rewardtypes.MsgSetRewardsResponse
	if err := res.Decode(&setRewardRes); err != nil {
		return err
	}

	if setRewardRes.PreviousCoins.Empty() {
		n.ev.Send(events.New(
			events.StatusDone,
			"The reward pool was empty",
			events.Icon(icons.Info),
		))
	} else {
		n.ev.Send(events.New(events.StatusDone,
			fmt.Sprintf(
				"Previous reward pool %s at height %d is overwritten",
				coins.String(),
				lastRewardHeight,
			),
			events.Icon(icons.Info),
		))
	}

	if setRewardRes.NewCoins.Empty() {
		n.ev.Send(events.New(events.StatusDone, "The reward pool is removed"))
	} else {
		n.ev.Send(events.New(events.StatusDone, fmt.Sprintf(
			"%s will be distributed to validators at height %d. The chain %d is now an incentivized testnet",
			coins.String(),
			lastRewardHeight,
			launchID,
		)))
	}
	return nil
}
