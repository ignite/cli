package network

import (
	"context"
	"errors"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	launchtypes "github.com/tendermint/spn/x/launch/types"
	"github.com/tendermint/starport/starport/pkg/cosmosaccount"
	"github.com/tendermint/starport/starport/services/network/mocks"
	"github.com/tendermint/starport/starport/services/network/networktypes"
	"github.com/tendermint/starport/starport/services/network/testutil"
)

const (
	TestMinRemainingTime = 3600
	TestMaxRemainingTime = 86400
	TestRevertDelay      = 3600
)

func stubNetworkForTriggerLaunch(account cosmosaccount.Account) Network {
	launchQueryMock := new(mocks.LaunchClient)
	launchQueryMock.On("Params", mock.Anything, &launchtypes.QueryParamsRequest{}).
		Return(&launchtypes.QueryParamsResponse{
			Params: launchtypes.NewParams(TestMinRemainingTime, TestMaxRemainingTime, TestRevertDelay, sdk.Coins(nil)),
		}, nil)
	networkClientMock := new(mocks.CosmosClient)
	networkClientMock.On("BroadcastTx", account.Name, &launchtypes.MsgTriggerLaunch{
		Coordinator:   account.Address(networktypes.SPN),
		LaunchID:      testutil.TestLaunchID,
		RemainingTime: TestMaxRemainingTime,
	}).Return(testutil.NewResponse(&launchtypes.MsgTriggerLaunchResponse{}), nil)
	return Network{
		cosmos:      networkClientMock,
		account:     account,
		launchQuery: launchQueryMock,
	}
}

func TestTriggerLaunch(t *testing.T) {
	account, err := testutil.NewTestAccount(testutil.TestAccountName)
	assert.Nil(t, err)

	network := stubNetworkForTriggerLaunch(account)
	err = network.TriggerLaunch(context.Background(), testutil.TestLaunchID, TestMaxRemainingTime*time.Second)
	require.Nil(t, err)
}

func TestTriggerLaunchRemainingTimeLowerThanAllowed(t *testing.T) {
	account, err := testutil.NewTestAccount(testutil.TestAccountName)
	assert.Nil(t, err)

	network := stubNetworkForTriggerLaunch(account)
	err = network.TriggerLaunch(context.Background(), testutil.TestLaunchID, (TestMinRemainingTime-60)*time.Second)
	require.NotNil(t, err)
}

func TestTriggerLaunchRemainingTimeGreaterThanAllowed(t *testing.T) {
	account, err := testutil.NewTestAccount(testutil.TestAccountName)
	assert.Nil(t, err)

	network := stubNetworkForTriggerLaunch(account)
	err = network.TriggerLaunch(context.Background(), testutil.TestLaunchID, (TestMaxRemainingTime+60)*time.Hour)
	require.NotNil(t, err)
}

func TestTriggerLaunchBroadcastFailure(t *testing.T) {
	account, err := testutil.NewTestAccount(testutil.TestAccountName)
	assert.Nil(t, err)

	network := stubNetworkForTriggerLaunch(account)
	networkClientMock := new(mocks.CosmosClient)
	networkClientMock.On("BroadcastTx", account.Name, &launchtypes.MsgTriggerLaunch{
		Coordinator:   account.Address(networktypes.SPN),
		LaunchID:      testutil.TestLaunchID,
		RemainingTime: TestMaxRemainingTime,
	}).Return(testutil.NewResponse(&launchtypes.MsgTriggerLaunch{}), errors.New("Failed to fetch"))
	network.cosmos = networkClientMock
	err = network.TriggerLaunch(context.Background(), testutil.TestLaunchID, TestMaxRemainingTime*time.Second)
	require.NotNil(t, err)
	networkClientMock.AssertCalled(t, "BroadcastTx", account.Name, &launchtypes.MsgTriggerLaunch{
		Coordinator:   account.Address(networktypes.SPN),
		LaunchID:      testutil.TestLaunchID,
		RemainingTime: TestMaxRemainingTime,
	})
}

func TestTriggerLaunchBadTriggerLaunchResponse(t *testing.T) {
	account, err := testutil.NewTestAccount(testutil.TestAccountName)
	assert.Nil(t, err)

	network := stubNetworkForTriggerLaunch(account)
	networkClientMock := new(mocks.CosmosClient)
	networkClientMock.On("BroadcastTx", account.Name, &launchtypes.MsgTriggerLaunch{
		Coordinator:   account.Address(networktypes.SPN),
		LaunchID:      testutil.TestLaunchID,
		RemainingTime: TestMaxRemainingTime,
	}).Return(testutil.NewResponse(&launchtypes.MsgCreateChainResponse{}), errors.New("failed to fetch"))
	network.cosmos = networkClientMock
	err = network.TriggerLaunch(context.Background(), testutil.TestLaunchID, TestMaxRemainingTime*time.Second)
	require.NotNil(t, err)
	networkClientMock.AssertCalled(t, "BroadcastTx", account.Name, &launchtypes.MsgTriggerLaunch{
		Coordinator:   account.Address(networktypes.SPN),
		LaunchID:      testutil.TestLaunchID,
		RemainingTime: TestMaxRemainingTime,
	})
}

func TestTriggerLaunchFailedToQueryChainParams(t *testing.T) {
	account, err := testutil.NewTestAccount(testutil.TestAccountName)
	assert.Nil(t, err)
	network := stubNetworkForTriggerLaunch(account)
	launchQueryMock := new(mocks.LaunchClient)
	launchQueryMock.On("Params", mock.Anything, &launchtypes.QueryParamsRequest{}).
		Return(&launchtypes.QueryParamsResponse{
			Params: launchtypes.NewParams(TestMinRemainingTime, TestMaxRemainingTime, TestRevertDelay, sdk.Coins(nil)),
		}, errors.New("failed to fetch"))
	network.launchQuery = launchQueryMock
	err = network.TriggerLaunch(context.Background(), testutil.TestLaunchID, (TestMaxRemainingTime+60)*time.Second)
	require.NotNil(t, err)
}
