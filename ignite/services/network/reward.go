package network

import (
	"context"
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

// RewardsInfo Fetches the consensus state with the validator set,
// the unbounding time, and the last block height from chain rewards.
func (n Network) RewardsInfo(
	ctx context.Context,
	launchID uint64,
	height int64,
) (
	rewardsInfo networktypes.Reward,
	lastRewardHeight int64,
	unboundingTime int64,
	err error,
) {
	rewardsInfo, err = n.node.consensus(ctx, n.cosmos, height)
	if err != nil {
		return rewardsInfo, 0, 0, err
	}

	stakingParams, err := n.node.stakingParams(ctx)
	if err != nil {
		return rewardsInfo, 0, 0, err
	}
	unboundingTime = int64(stakingParams.UnbondingTime.Seconds())

	chainReward, err := n.ChainReward(ctx, launchID)
	if err == ErrObjectNotFound {
		return rewardsInfo, 1, unboundingTime, nil
	} else if err != nil {
		return rewardsInfo, 0, 0, err
	}
	lastRewardHeight = chainReward.LastRewardHeight

	return
}
