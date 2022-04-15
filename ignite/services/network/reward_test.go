package network

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ignite-hq/cli/ignite/services/network/networktypes"
	"github.com/ignite-hq/cli/ignite/services/network/testutil"
	rewardtypes "github.com/tendermint/spn/x/reward/types"
)

func TestSetReward(t *testing.T) {
	t.Run("successfully set reward", func(t *testing.T) {
		var (
			account         = testutil.NewTestAccount(t, testutil.TestAccountName)
			suite, network  = newSuite(account)
			coins           = sdk.NewCoins(sdk.NewCoin(TestDenom, sdk.NewInt(TestAmountInt)))
			lastRewarHeight = int64(10)
		)
		suite.CosmosClientMock.
			On(
				"BroadcastTx",
				account.Name,
				&rewardtypes.MsgSetRewards{
					Provider:         account.Address(networktypes.SPN),
					LaunchID:         testutil.LaunchID,
					Coins:            coins,
					LastRewardHeight: lastRewarHeight,
				},
			).
			Return(testutil.NewResponse(&rewardtypes.MsgSetRewardsResponse{
				PreviousCoins:            nil,
				PreviousLastRewardHeight: lastRewarHeight - 1,
				NewCoins:                 coins,
				NewLastRewardHeight:      lastRewarHeight,
			}), nil).
			Once()
		setRewardError := network.SetReward(testutil.LaunchID, lastRewarHeight, coins)
		require.NoError(t, setRewardError)
		suite.AssertAllMocks(t)
	})
	t.Run("failed to set reward, failed to broadcast set reward tx", func(t *testing.T) {
		var (
			account         = testutil.NewTestAccount(t, testutil.TestAccountName)
			suite, network  = newSuite(account)
			coins           = sdk.NewCoins(sdk.NewCoin(TestDenom, sdk.NewInt(TestAmountInt)))
			lastRewarHeight = int64(10)
		)
		suite.CosmosClientMock.
			On(
				"BroadcastTx",
				account.Name,
				&rewardtypes.MsgSetRewards{
					Provider:         account.Address(networktypes.SPN),
					LaunchID:         testutil.LaunchID,
					Coins:            coins,
					LastRewardHeight: lastRewarHeight,
				},
			).
			Return(testutil.NewResponse(&rewardtypes.MsgSetRewardsResponse{}), errors.New("failed to set reward")).
			Once()
		setRewardError := network.SetReward(testutil.LaunchID, lastRewarHeight, coins)
		require.Error(t, setRewardError)
		suite.AssertAllMocks(t)
	})
}
