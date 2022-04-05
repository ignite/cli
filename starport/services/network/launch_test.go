package network

import (
	"context"
	"errors"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	launchtypes "github.com/tendermint/spn/x/launch/types"
	"github.com/tendermint/starport/starport/services/network/mocks"
	"github.com/tendermint/starport/starport/services/network/testdata"
)

const (
	TestMinRemainingTime = 3600
	TestMaxRemainingTime = 86400
	TestRevertDelay      = 3600
)

func stubNetworkForTriggerLaunch() Network {
	launchQueryMock := new(mocks.LaunchClient)
	launchQueryMock.On("Params", mock.Anything, &launchtypes.QueryParamsRequest{}).
		Return(&launchtypes.QueryParamsResponse{
			Params: launchtypes.NewParams(TestMinRemainingTime, TestMaxRemainingTime, TestRevertDelay, sdk.Coins(nil)),
		}, nil)
	networkClientMock := new(mocks.CosmosClient)
	networkClientMock.On("BroadcastTx", testdata.AccountName, &launchtypes.MsgTriggerLaunch{
		Coordinator:   testdata.Address,
		LaunchID:      TestLaunchID,
		RemainingTime: TestMaxRemainingTime,
	}).Return(testdata.NewResponse(&launchtypes.MsgTriggerLaunchResponse{}), nil)
	return Network{
		cosmos:      networkClientMock,
		account:     testdata.GetTestAccount(),
		launchQuery: launchQueryMock,
	}
}

func TestTriggerLaunch(t *testing.T) {
	network := stubNetworkForTriggerLaunch()
	err := network.TriggerLaunch(context.Background(), TestLaunchID, TestMaxRemainingTime*time.Second)
	require.Nil(t, err)
}

func TestTriggerLaunchRemainingTimeLowerThanAllowed(t *testing.T) {
	network := stubNetworkForTriggerLaunch()
	err := network.TriggerLaunch(context.Background(), TestLaunchID, (TestMinRemainingTime-60)*time.Second)
	require.NotNil(t, err)
}

func TestTriggerLaunchRemainingTimeGreaterThanAllowed(t *testing.T) {
	network := stubNetworkForTriggerLaunch()
	err := network.TriggerLaunch(context.Background(), TestLaunchID, (TestMaxRemainingTime+60)*time.Hour)
	require.NotNil(t, err)
}

func TestTriggerLaunchBroadcastFailure(t *testing.T) {
	network := stubNetworkForTriggerLaunch()
	networkClientMock := new(mocks.CosmosClient)
	networkClientMock.On("BroadcastTx", testdata.AccountName, &launchtypes.MsgTriggerLaunch{
		Coordinator:   testdata.Address,
		LaunchID:      TestLaunchID,
		RemainingTime: TestMaxRemainingTime,
	}).Return(testdata.NewResponse(&launchtypes.MsgTriggerLaunch{}), errors.New("Failed to fetch"))
	network.cosmos = networkClientMock
	err := network.TriggerLaunch(context.Background(), TestLaunchID, TestMaxRemainingTime*time.Second)
	require.NotNil(t, err)
	networkClientMock.AssertCalled(t, "BroadcastTx", testdata.AccountName, &launchtypes.MsgTriggerLaunch{
		Coordinator:   testdata.Address,
		LaunchID:      TestLaunchID,
		RemainingTime: TestMaxRemainingTime,
	})
}

func TestTriggerLaunchBadTriggerLaunchResponse(t *testing.T) {
	network := stubNetworkForTriggerLaunch()
	networkClientMock := new(mocks.CosmosClient)
	networkClientMock.On("BroadcastTx", testdata.AccountName, &launchtypes.MsgTriggerLaunch{
		Coordinator:   testdata.Address,
		LaunchID:      TestLaunchID,
		RemainingTime: TestMaxRemainingTime,
	}).Return(testdata.NewResponse(&launchtypes.MsgCreateChainResponse{}), errors.New("failed to fetch"))
	network.cosmos = networkClientMock
	err := network.TriggerLaunch(context.Background(), TestLaunchID, TestMaxRemainingTime*time.Second)
	require.NotNil(t, err)
	networkClientMock.AssertCalled(t, "BroadcastTx", testdata.AccountName, &launchtypes.MsgTriggerLaunch{
		Coordinator:   testdata.Address,
		LaunchID:      TestLaunchID,
		RemainingTime: TestMaxRemainingTime,
	})
}

func TestTriggerLaunchFailedToQueryChainParams(t *testing.T) {
	network := stubNetworkForTriggerLaunch()
	launchQueryMock := new(mocks.LaunchClient)
	launchQueryMock.On("Params", mock.Anything, &launchtypes.QueryParamsRequest{}).
		Return(&launchtypes.QueryParamsResponse{
			Params: launchtypes.NewParams(TestMinRemainingTime, TestMaxRemainingTime, TestRevertDelay, sdk.Coins(nil)),
		}, errors.New("failed to fetch"))
	network.launchQuery = launchQueryMock
	err := network.TriggerLaunch(context.Background(), TestLaunchID, (TestMaxRemainingTime+60)*time.Second)
	require.NotNil(t, err)
}
