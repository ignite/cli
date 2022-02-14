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
	n.ev.Send(events.New(events.StatusOngoing, fmt.Sprintf("Setting reward %s to chain %d", coins.String(), launchID)))

	msg := rewardtypes.NewMsgSetRewards(
		n.account.Address(networktypes.SPN),
		launchID,
		lastRewardHeight,
		coins)
	res, err := n.cosmos.BroadcastTx(n.account.Name, msg)
	if err != nil {
		return err
	}

	var setRewardRes rewardtypes.MsgSetRewardsResponse
	if err := res.Decode(&setRewardRes); err != nil {
		return err
	}

	n.ev.Send(events.New(events.StatusDone,
		fmt.Sprintf("Chain %d reward set: %s", launchID, coins.String()),
	))
	return nil
}
