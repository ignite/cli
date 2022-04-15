package network

import (
	"context"
	"errors"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ignite-hq/cli/ignite/services/network/networktypes"
	"github.com/ignite-hq/cli/ignite/services/network/testutil"
	"github.com/stretchr/testify/require"
	launchtypes "github.com/tendermint/spn/x/launch/types"
)

const (
	TestMinRemainingTime = 3600
	TestMaxRemainingTime = 86400
	TestRevertDelay      = 3600
)

func TestTriggerLaunch(t *testing.T) {
	t.Run("successfully launch a chain", func(t *testing.T) {
		var (
			account        = testutil.NewTestAccount(t, testutil.TestAccountName)
			suite, network = newSuite(account)
		)

		suite.LaunchQueryMock.
			On("Params", context.Background(), &launchtypes.QueryParamsRequest{}).
			Return(&launchtypes.QueryParamsResponse{
				Params: launchtypes.NewParams(TestMinRemainingTime, TestMaxRemainingTime, TestRevertDelay, sdk.Coins(nil)),
			}, nil).
			Once()
		suite.CosmosClientMock.
			On("BroadcastTx", account.Name, &launchtypes.MsgTriggerLaunch{
				Coordinator:   account.Address(networktypes.SPN),
				LaunchID:      testutil.LaunchID,
				RemainingTime: TestMaxRemainingTime,
			}).
			Return(testutil.NewResponse(&launchtypes.MsgTriggerLaunchResponse{}), nil).
			Once()

		err := network.TriggerLaunch(context.Background(), testutil.LaunchID, TestMaxRemainingTime*time.Second)
		require.NoError(t, err)
		suite.AssertAllMocks(t)
	})

	t.Run("failed to launch a chain, remaining time is lower than allowed", func(t *testing.T) {
		var (
			account        = testutil.NewTestAccount(t, testutil.TestAccountName)
			suite, network = newSuite(account)
		)

		suite.LaunchQueryMock.
			On("Params", context.Background(), &launchtypes.QueryParamsRequest{}).
			Return(&launchtypes.QueryParamsResponse{
				Params: launchtypes.NewParams(TestMinRemainingTime, TestMaxRemainingTime, TestRevertDelay, sdk.Coins(nil)),
			}, nil).
			Once()

		err := network.TriggerLaunch(context.Background(), testutil.LaunchID, (TestMinRemainingTime-60)*time.Second)
		require.Error(t, err)
		suite.AssertAllMocks(t)
	})

	t.Run("failed to launch a chain, remaining time is greater than allowed", func(t *testing.T) {
		var (
			account        = testutil.NewTestAccount(t, testutil.TestAccountName)
			suite, network = newSuite(account)
		)

		suite.LaunchQueryMock.
			On("Params", context.Background(), &launchtypes.QueryParamsRequest{}).
			Return(&launchtypes.QueryParamsResponse{
				Params: launchtypes.NewParams(TestMinRemainingTime, TestMaxRemainingTime, TestRevertDelay, sdk.Coins(nil)),
			}, nil).
			Once()

		err := network.TriggerLaunch(context.Background(), testutil.LaunchID, (TestMaxRemainingTime+60)*time.Hour)
		require.Error(t, err)
		suite.AssertAllMocks(t)
	})

	t.Run("failed to launch a chain, failed to broadcast the launch tx", func(t *testing.T) {
		var (
			account        = testutil.NewTestAccount(t, testutil.TestAccountName)
			suite, network = newSuite(account)
		)

		suite.LaunchQueryMock.
			On("Params", context.Background(), &launchtypes.QueryParamsRequest{}).
			Return(&launchtypes.QueryParamsResponse{
				Params: launchtypes.NewParams(TestMinRemainingTime, TestMaxRemainingTime, TestRevertDelay, sdk.Coins(nil)),
			}, nil).
			Once()
		suite.CosmosClientMock.
			On("BroadcastTx", account.Name, &launchtypes.MsgTriggerLaunch{
				Coordinator:   account.Address(networktypes.SPN),
				LaunchID:      testutil.LaunchID,
				RemainingTime: TestMaxRemainingTime,
			}).
			Return(testutil.NewResponse(&launchtypes.MsgTriggerLaunch{}), errors.New("Failed to fetch")).
			Once()

		err := network.TriggerLaunch(context.Background(), testutil.LaunchID, TestMaxRemainingTime*time.Second)
		require.Error(t, err)
		suite.AssertAllMocks(t)
	})

	t.Run("failed to launch a chain, invalid response from chain", func(t *testing.T) {
		var (
			account        = testutil.NewTestAccount(t, testutil.TestAccountName)
			suite, network = newSuite(account)
		)

		suite.LaunchQueryMock.
			On("Params", context.Background(), &launchtypes.QueryParamsRequest{}).
			Return(&launchtypes.QueryParamsResponse{
				Params: launchtypes.NewParams(TestMinRemainingTime, TestMaxRemainingTime, TestRevertDelay, sdk.Coins(nil)),
			}, nil).
			Once()
		suite.CosmosClientMock.
			On("BroadcastTx", account.Name, &launchtypes.MsgTriggerLaunch{
				Coordinator:   account.Address(networktypes.SPN),
				LaunchID:      testutil.LaunchID,
				RemainingTime: TestMaxRemainingTime,
			}).
			Return(testutil.NewResponse(&launchtypes.MsgCreateChainResponse{}), errors.New("failed to fetch")).
			Once()

		err := network.TriggerLaunch(context.Background(), testutil.LaunchID, TestMaxRemainingTime*time.Second)
		require.Error(t, err)
		suite.AssertAllMocks(t)
	})

	t.Run("failed to launch a chain, failed to fetch chain params", func(t *testing.T) {
		var (
			account        = testutil.NewTestAccount(t, testutil.TestAccountName)
			suite, network = newSuite(account)
		)

		suite.LaunchQueryMock.
			On("Params", context.Background(), &launchtypes.QueryParamsRequest{}).
			Return(&launchtypes.QueryParamsResponse{
				Params: launchtypes.NewParams(TestMinRemainingTime, TestMaxRemainingTime, TestRevertDelay, sdk.Coins(nil)),
			}, errors.New("failed to fetch")).
			Once()

		err := network.TriggerLaunch(context.Background(), testutil.LaunchID, (TestMaxRemainingTime+60)*time.Second)
		require.Error(t, err)
		suite.AssertAllMocks(t)
	})
}
