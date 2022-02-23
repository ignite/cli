package network

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
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

	doneMsg := fmt.Sprintf("%s will be distributed to validators at height %d.", coins.String(), lastRewardHeight)
	if setRewardRes.PreviousCoins.Empty() {
		doneMsg = fmt.Sprintf("New reward pool was created.\n%s", doneMsg)
	} else {
		doneMsg = fmt.Sprintf("Previous reward pool was overwritten.\n%s", doneMsg)
	}

	if setRewardRes.NewCoins.Empty() {
		n.ev.Send(events.New(events.StatusDone, "The reward pool was removed."))
	} else {
		n.ev.Send(events.New(events.StatusDone, fmt.Sprintf(
			"%s.\nThe chain %d is now an incentivized testnet",
			doneMsg,
			launchID,
		)))
	}
	return nil
}
