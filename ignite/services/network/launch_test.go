package network

import (
	"context"
	"errors"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	launchtypes "github.com/tendermint/spn/x/launch/types"

	"github.com/ignite/cli/ignite/services/network/networktypes"
	"github.com/ignite/cli/ignite/services/network/testutil"
)

const (
	TestMinRemainingTime = time.Second * 3600
	TestMaxRemainingTime = time.Second * 86400
	TestRevertDelay      = time.Second * 3600
)

func TestTriggerLaunch(t *testing.T) {
	t.Run("successfully launch a chain", func(t *testing.T) {
		var (
			account        = testutil.NewTestAccount(t, testutil.TestAccountName)
			suite, network = newSuite(account)
		)

		addr, err := account.Address(networktypes.SPN)
		require.NoError(t, err)

		suite.LaunchQueryMock.
			On("Params", context.Background(), &launchtypes.QueryParamsRequest{}).
			Return(&launchtypes.QueryParamsResponse{
				Params: launchtypes.NewParams(
					TestMinRemainingTime,
					TestMaxRemainingTime,
					TestRevertDelay,
					sdk.Coins(nil),
					sdk.Coins(nil),
				),
			}, nil).
			Once()
		suite.CosmosClientMock.
			On("BroadcastTx",
				context.Background(),
				account,
				&launchtypes.MsgTriggerLaunch{
					Coordinator: addr,
					LaunchID:    testutil.LaunchID,
					LaunchTime:  sampleTime.Add(TestMaxRemainingTime),
				}).
			Return(testutil.NewResponse(&launchtypes.MsgTriggerLaunchResponse{}), nil).
			Once()

		launchError := network.TriggerLaunch(context.Background(), testutil.LaunchID, sampleTime.Add(TestMaxRemainingTime))
		require.NoError(t, launchError)
		suite.AssertAllMocks(t)
	})

	t.Run("failed to launch a chain, remaining time is lower than allowed", func(t *testing.T) {
		var (
			account                       = testutil.NewTestAccount(t, testutil.TestAccountName)
			suite, network                = newSuite(account)
			remainingTimeLowerThanMinimum = sampleTime
		)

		suite.LaunchQueryMock.
			On("Params", context.Background(), &launchtypes.QueryParamsRequest{}).
			Return(&launchtypes.QueryParamsResponse{
				Params: launchtypes.NewParams(
					TestMinRemainingTime,
					TestMaxRemainingTime,
					TestRevertDelay,
					sdk.Coins(nil),
					sdk.Coins(nil),
				),
			}, nil).
			Once()

		launchError := network.TriggerLaunch(context.Background(), testutil.LaunchID, remainingTimeLowerThanMinimum)
		require.Errorf(
			t,
			launchError,
			"remaining time %s lower than minimum %s",
			remainingTimeLowerThanMinimum.String(),
			sampleTime.Add(TestMinRemainingTime).Add(MinLaunchTimeOffset).String(),
		)
		suite.AssertAllMocks(t)
	})

	t.Run("failed to launch a chain, remaining time is greater than allowed", func(t *testing.T) {
		var (
			account                         = testutil.NewTestAccount(t, testutil.TestAccountName)
			suite, network                  = newSuite(account)
			remainingTimeGreaterThanMaximum = sampleTime.Add(TestMaxRemainingTime).Add(time.Second)
		)

		suite.LaunchQueryMock.
			On("Params", context.Background(), &launchtypes.QueryParamsRequest{}).
			Return(&launchtypes.QueryParamsResponse{
				Params: launchtypes.NewParams(
					TestMinRemainingTime,
					TestMaxRemainingTime,
					TestRevertDelay,
					sdk.Coins(nil),
					sdk.Coins(nil),
				),
			}, nil).
			Once()

		launchError := network.TriggerLaunch(context.Background(), testutil.LaunchID, remainingTimeGreaterThanMaximum)
		require.Errorf(
			t,
			launchError,
			"remaining time %s greater than maximum %s",
			remainingTimeGreaterThanMaximum.String(),
			sampleTime.Add(TestMaxRemainingTime).String(),
		)
		suite.AssertAllMocks(t)
	})

	t.Run("failed to launch a chain, failed to broadcast the launch tx", func(t *testing.T) {
		var (
			account        = testutil.NewTestAccount(t, testutil.TestAccountName)
			suite, network = newSuite(account)
			expectedError  = errors.New("Failed to fetch")
		)

		addr, err := account.Address(networktypes.SPN)
		require.NoError(t, err)

		suite.LaunchQueryMock.
			On("Params", context.Background(), &launchtypes.QueryParamsRequest{}).
			Return(&launchtypes.QueryParamsResponse{
				Params: launchtypes.NewParams(
					TestMinRemainingTime,
					TestMaxRemainingTime,
					TestRevertDelay,
					sdk.Coins(nil),
					sdk.Coins(nil),
				),
			}, nil).
			Once()
		suite.CosmosClientMock.
			On("BroadcastTx",
				context.Background(),
				account,
				&launchtypes.MsgTriggerLaunch{
					Coordinator: addr,
					LaunchID:    testutil.LaunchID,
					LaunchTime:  sampleTime.Add(TestMaxRemainingTime),
				}).
			Return(testutil.NewResponse(&launchtypes.MsgTriggerLaunch{}), expectedError).
			Once()

		launchError := network.TriggerLaunch(context.Background(), testutil.LaunchID, sampleTime.Add(TestMaxRemainingTime))
		require.Error(t, launchError)
		require.Equal(t, expectedError, launchError)
		suite.AssertAllMocks(t)
	})

	t.Run("failed to launch a chain, invalid response from chain", func(t *testing.T) {
		var (
			account        = testutil.NewTestAccount(t, testutil.TestAccountName)
			suite, network = newSuite(account)
			expectedError  = errors.New("failed to fetch")
		)

		addr, err := account.Address(networktypes.SPN)
		require.NoError(t, err)

		suite.LaunchQueryMock.
			On("Params", context.Background(), &launchtypes.QueryParamsRequest{}).
			Return(&launchtypes.QueryParamsResponse{
				Params: launchtypes.NewParams(
					TestMinRemainingTime,
					TestMaxRemainingTime,
					TestRevertDelay,
					sdk.Coins(nil),
					sdk.Coins(nil),
				),
			}, nil).
			Once()
		suite.CosmosClientMock.
			On("BroadcastTx",
				context.Background(),
				account,
				&launchtypes.MsgTriggerLaunch{
					Coordinator: addr,
					LaunchID:    testutil.LaunchID,
					LaunchTime:  sampleTime.Add(TestMaxRemainingTime),
				}).
			Return(testutil.NewResponse(&launchtypes.MsgCreateChainResponse{}), expectedError).
			Once()

		launchError := network.TriggerLaunch(context.Background(), testutil.LaunchID, sampleTime.Add(TestMaxRemainingTime))
		require.Error(t, launchError)
		require.Equal(t, expectedError, launchError)
		suite.AssertAllMocks(t)
	})

	t.Run("failed to launch a chain, failed to fetch chain params", func(t *testing.T) {
		var (
			account        = testutil.NewTestAccount(t, testutil.TestAccountName)
			suite, network = newSuite(account)
			expectedError  = errors.New("failed to fetch")
		)

		suite.LaunchQueryMock.
			On("Params", context.Background(), &launchtypes.QueryParamsRequest{}).
			Return(&launchtypes.QueryParamsResponse{
				Params: launchtypes.NewParams(
					TestMinRemainingTime,
					TestMaxRemainingTime,
					TestRevertDelay,
					sdk.Coins(nil),
					sdk.Coins(nil),
				),
			}, expectedError).
			Once()

		launchError := network.TriggerLaunch(context.Background(), testutil.LaunchID, sampleTime.Add(TestMaxRemainingTime))
		require.Error(t, launchError)
		require.Equal(t, expectedError, launchError)
		suite.AssertAllMocks(t)
	})
}

func TestRevertLaunch(t *testing.T) {
	t.Run("successfully revert launch", func(t *testing.T) {
		var (
			account        = testutil.NewTestAccount(t, testutil.TestAccountName)
			suite, network = newSuite(account)
		)

		addr, err := account.Address(networktypes.SPN)
		require.NoError(t, err)

		suite.CosmosClientMock.
			On("BroadcastTx",
				context.Background(),
				account,
				&launchtypes.MsgRevertLaunch{
					Coordinator: addr,
					LaunchID:    testutil.LaunchID,
				}).
			Return(testutil.NewResponse(&launchtypes.MsgRevertLaunchResponse{}), nil).
			Once()

		revertError := network.RevertLaunch(context.Background(), testutil.LaunchID, suite.ChainMock)
		require.NoError(t, revertError)
		suite.AssertAllMocks(t)
	})

	t.Run("failed to revert launch, failed to broadcast revert launch tx", func(t *testing.T) {
		var (
			account        = testutil.NewTestAccount(t, testutil.TestAccountName)
			suite, network = newSuite(account)
			expectedError  = errors.New("failed to revert launch")
		)

		addr, err := account.Address(networktypes.SPN)
		require.NoError(t, err)

		suite.CosmosClientMock.
			On("BroadcastTx",
				context.Background(),
				account,
				&launchtypes.MsgRevertLaunch{
					Coordinator: addr,
					LaunchID:    testutil.LaunchID,
				}).
			Return(
				testutil.NewResponse(&launchtypes.MsgRevertLaunchResponse{}),
				expectedError,
			).
			Once()

		revertError := network.RevertLaunch(context.Background(), testutil.LaunchID, suite.ChainMock)
		require.Error(t, revertError)
		require.Equal(t, expectedError, revertError)
		suite.AssertAllMocks(t)
	})
}
