package network

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gookit/color"
	rewardtypes "github.com/tendermint/spn/x/reward/types"

	"github.com/tendermint/starport/starport/pkg/events"
	"github.com/tendermint/starport/starport/services/network/networktypes"
)

// SetReward set a chain reward
func (n Network) SetReward(launchID, lastRewardHeight uint64, coins sdk.Coins) error {
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
			events.StatusInfo,
			"The reward pool was removed.",
			events.Color(color.Yellow),
		))
	} else {
		n.ev.Send(events.New(events.StatusInfo,
			fmt.Sprintf(
				"Previous reward pool %s at height %d was overwritten.",
				coins.String(),
				lastRewardHeight,
			),
			events.Color(color.Yellow),
		))
	}

	if setRewardRes.NewCoins.Empty() {
		n.ev.Send(events.New(events.StatusDone, "The reward pool was removed."))
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
