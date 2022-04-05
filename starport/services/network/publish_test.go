package network

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	campaigntypes "github.com/tendermint/spn/x/campaign/types"
	launchtypes "github.com/tendermint/spn/x/launch/types"
	profiletypes "github.com/tendermint/spn/x/profile/types"
	"github.com/tendermint/starport/starport/pkg/cosmoserror"
	"github.com/tendermint/starport/starport/pkg/cosmosutil"
	"github.com/tendermint/starport/starport/services/network/mocks"
	"github.com/tendermint/starport/starport/services/network/testdata"
)

const (
	TestChainSourceHash = "testhash"
	TestChainSourceURL  = "http://example.com/test"
	TestChainName       = "test"
	TestChainChainID    = "test-1"

	TestLaunchID   = uint64(1)
	TestCampaignID = uint64(1)
	TestMainnetID  = uint64(1)
)

func stubNetworkForPublish() Network {
	profileQueryMock := new(mocks.ProfileClient)
	profileQueryMock.On("CoordinatorByAddress", mock.Anything, &profiletypes.QueryGetCoordinatorByAddressRequest{
		Address: testdata.Address,
	}).Return(nil, nil)
	campaignQueryMock := new(mocks.CampaignClient)
	campaignQueryMock.On("Campaign", mock.Anything, &campaigntypes.QueryGetCampaignRequest{
		CampaignID: TestCampaignID,
	}).Return(nil, nil)
	networkClientMock := new(mocks.CosmosClient)
	networkClientMock.On("BroadcastTx",
		testdata.AccountName,
		mock.AnythingOfType("*types.MsgCreateChain"),
	).Return(testdata.NewResponse(&launchtypes.MsgCreateChainResponse{
		LaunchID: TestLaunchID,
	}), nil)
	networkClientMock.On("BroadcastTx",
		testdata.AccountName,
		mock.AnythingOfType("*types.MsgCreateCampaign"),
	).Return(testdata.NewResponse(&campaigntypes.MsgCreateCampaignResponse{
		CampaignID: TestCampaignID,
	}), nil)
	networkClientMock.On("BroadcastTx",
		testdata.AccountName,
		mock.AnythingOfType("*types.MsgCreateCoordinator"),
	).Return(testdata.NewResponse(&profiletypes.MsgCreateCoordinatorResponse{}), nil)
	return Network{
		cosmos:        networkClientMock,
		account:       testdata.GetTestAccount(),
		profileQuery:  profileQueryMock,
		campaignQuery: campaignQueryMock,
	}
}

func newChainMockForPublish(cacheBinaryFail, chainIDFail bool) *mocks.Chain {
	chainMock := new(mocks.Chain)
	chainMock.On("SourceHash").Return(TestChainSourceHash)
	chainMock.On("SourceURL").Return(TestChainSourceURL)
	chainMock.On("Name").Return(TestChainName)
	if cacheBinaryFail {
		chainMock.On("CacheBinary", TestLaunchID).Return(errors.New("failed to cache binary"))
	} else {
		chainMock.On("CacheBinary", TestLaunchID).Return(nil)
	}
	if chainIDFail {
		chainMock.On("ChainID").Return("", errors.New("failed to get chainID"))
	} else {
		chainMock.On("ChainID").Return(TestChainChainID, nil)
	}
	return chainMock
}

func startGenesisTestServer(genesis cosmosutil.ChainGenesis) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		encodedGenesis, _ := json.Marshal(genesis)
		w.Write(encodedGenesis)
	}))
}

func startNotFoundServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	}))
}

func TestPublish(t *testing.T) {
	network := stubNetworkForPublish()
	chainMock := newChainMockForPublish(false, false)
	launchID, campaignID, _, err := network.Publish(context.Background(), chainMock)
	require.Nil(t, err)
	require.Equal(t, TestLaunchID, launchID)
	require.Equal(t, TestCampaignID, campaignID)
	chainMock.AssertNumberOfCalls(t, "ChainID", 1)
	network.profileQuery.(*mocks.ProfileClient).AssertNumberOfCalls(t, "CoordinatorByAddress", 1)
	network.cosmos.(*mocks.CosmosClient).AssertNumberOfCalls(t, "BroadcastTx", 2)
	chainMock.AssertNumberOfCalls(t, "CacheBinary", 1)
}

func TestPublishWithCampaignID(t *testing.T) {
	network := stubNetworkForPublish()
	chainMock := newChainMockForPublish(false, false)
	launchID, campaignID, _, err := network.Publish(context.Background(), chainMock, WithCampaign(TestCampaignID))
	require.Nil(t, err)
	require.Equal(t, TestLaunchID, launchID)
	require.Equal(t, TestCampaignID, campaignID)
	chainMock.AssertNumberOfCalls(t, "ChainID", 1)
	network.profileQuery.(*mocks.ProfileClient).AssertNumberOfCalls(t, "CoordinatorByAddress", 1)
	network.campaignQuery.(*mocks.CampaignClient).AssertNumberOfCalls(t, "Campaign", 1)
	network.cosmos.(*mocks.CosmosClient).AssertNumberOfCalls(t, "BroadcastTx", 1)
	chainMock.AssertNumberOfCalls(t, "CacheBinary", 1)
}

func TestPublishWithCustomGenesis(t *testing.T) {
	var customGenesisChainID string = "test-custom-1"
	network := stubNetworkForPublish()
	chainMock := newChainMockForPublish(false, false)
	gts := startGenesisTestServer(cosmosutil.ChainGenesis{ChainID: customGenesisChainID})
	defer gts.Close()
	launchID, campaignID, _, err := network.Publish(context.Background(), chainMock, WithCustomGenesis(gts.URL))
	require.Nil(t, err)
	require.Equal(t, TestLaunchID, launchID)
	require.Equal(t, TestCampaignID, campaignID)
	chainMock.AssertNumberOfCalls(t, "ChainID", 0)
	network.campaignQuery.(*mocks.CampaignClient).AssertNumberOfCalls(t, "Campaign", 0)
	network.profileQuery.(*mocks.ProfileClient).AssertNumberOfCalls(t, "CoordinatorByAddress", 1)
	network.cosmos.(*mocks.CosmosClient).AssertNumberOfCalls(t, "BroadcastTx", 2)
	chainMock.AssertNumberOfCalls(t, "CacheBinary", 1)
}

func TestPublishWithCustomChainID(t *testing.T) {
	var customGenesisChainID string = "test-custom-1"
	network := stubNetworkForPublish()
	chainMock := newChainMockForPublish(false, false)
	gts := startGenesisTestServer(cosmosutil.ChainGenesis{ChainID: customGenesisChainID})
	defer gts.Close()
	launchID, campaignID, _, err := network.Publish(context.Background(), chainMock, WithChainID(TestChainChainID))
	require.Nil(t, err)
	require.Equal(t, TestLaunchID, launchID)
	require.Equal(t, TestCampaignID, campaignID)
	chainMock.AssertNumberOfCalls(t, "ChainID", 0)
	network.profileQuery.(*mocks.ProfileClient).AssertNumberOfCalls(t, "CoordinatorByAddress", 1)
	network.cosmos.(*mocks.CosmosClient).AssertNumberOfCalls(t, "BroadcastTx", 2)
	chainMock.AssertNumberOfCalls(t, "CacheBinary", 1)
}

func TestPublishMainnet(t *testing.T) {
	network := stubNetworkForPublish()
	network.cosmos.(*mocks.CosmosClient).On(
		"BroadcastTx",
		testdata.AccountName,
		&campaigntypes.MsgInitializeMainnet{
			Coordinator:    testdata.Address,
			CampaignID:     TestCampaignID,
			SourceURL:      TestChainSourceURL,
			SourceHash:     TestChainSourceHash,
			MainnetChainID: TestChainChainID,
		},
	).Return(testdata.NewResponse(&campaigntypes.MsgInitializeMainnetResponse{
		MainnetID: TestMainnetID,
	}), nil)
	chainMock := newChainMockForPublish(false, false)
	launchID, campaignID, mainnetID, err := network.Publish(context.Background(), chainMock, Mainnet())
	require.Nil(t, err)
	require.Equal(t, TestLaunchID, launchID)
	require.Equal(t, TestCampaignID, campaignID)
	require.Equal(t, TestMainnetID, mainnetID)
	chainMock.AssertNumberOfCalls(t, "ChainID", 1)
	network.profileQuery.(*mocks.ProfileClient).AssertNumberOfCalls(t, "CoordinatorByAddress", 1)
	network.cosmos.(*mocks.CosmosClient).AssertNumberOfCalls(t, "BroadcastTx", 3)
	chainMock.AssertNumberOfCalls(t, "CacheBinary", 1)
}

func TestPublishWithCustomGenesisAndFailedToFetchIt(t *testing.T) {
	network := stubNetworkForPublish()
	chainMock := newChainMockForPublish(false, false)
	gts := startNotFoundServer()
	defer gts.Close()
	_, _, _, err := network.Publish(context.Background(), chainMock, WithCustomGenesis(gts.URL))
	require.NotNil(t, err)
	chainMock.AssertNumberOfCalls(t, "ChainID", 0)
	network.campaignQuery.(*mocks.CampaignClient).AssertNumberOfCalls(t, "Campaign", 0)
	network.profileQuery.(*mocks.ProfileClient).AssertNumberOfCalls(t, "CoordinatorByAddress", 0)
	network.cosmos.(*mocks.CosmosClient).AssertNumberOfCalls(t, "BroadcastTx", 0)
	chainMock.AssertNumberOfCalls(t, "CacheBinary", 0)
}

func TestPublishFailedToInitializeMainnet(t *testing.T) {
	network := stubNetworkForPublish()
	network.cosmos.(*mocks.CosmosClient).On(
		"BroadcastTx",
		testdata.AccountName,
		&campaigntypes.MsgInitializeMainnet{
			Coordinator:    testdata.Address,
			CampaignID:     TestCampaignID,
			SourceURL:      TestChainSourceURL,
			SourceHash:     TestChainSourceHash,
			MainnetChainID: TestChainChainID,
		},
	).Return(testdata.NewResponse(&campaigntypes.MsgInitializeMainnetResponse{
		MainnetID: TestMainnetID,
	}), errors.New("failed to initialize mainnet"))
	chainMock := newChainMockForPublish(false, false)
	_, _, _, err := network.Publish(context.Background(), chainMock, Mainnet())
	require.NotNil(t, err)
	chainMock.AssertNumberOfCalls(t, "ChainID", 1)
	network.profileQuery.(*mocks.ProfileClient).AssertNumberOfCalls(t, "CoordinatorByAddress", 1)
	network.cosmos.(*mocks.CosmosClient).AssertNumberOfCalls(t, "BroadcastTx", 3)
	chainMock.AssertNumberOfCalls(t, "CacheBinary", 1)
}

func TestPublishWithCustomGenesisAndFailedToParseIt(t *testing.T) {
	var customGenesisChainID string = "test-custom-1"
	network := stubNetworkForPublish()
	chainMock := newChainMockForPublish(false, false)
	gts := startGenesisTestServer(cosmosutil.ChainGenesis{ChainID: customGenesisChainID})
	defer gts.Close()
	launchID, campaignID, _, err := network.Publish(context.Background(), chainMock, WithCustomGenesis(gts.URL))
	require.Nil(t, err)
	require.Equal(t, TestLaunchID, launchID)
	require.Equal(t, TestCampaignID, campaignID)
	chainMock.AssertNumberOfCalls(t, "ChainID", 0)
	network.campaignQuery.(*mocks.CampaignClient).AssertNumberOfCalls(t, "Campaign", 0)
	network.profileQuery.(*mocks.ProfileClient).AssertNumberOfCalls(t, "CoordinatorByAddress", 1)
	network.cosmos.(*mocks.CosmosClient).AssertNumberOfCalls(t, "BroadcastTx", 2)
	chainMock.AssertNumberOfCalls(t, "CacheBinary", 1)
}

func TestPublishWithCoordinatorCreation(t *testing.T) {
	network := stubNetworkForPublish()
	profileQueryMock := new(mocks.ProfileClient)
	profileQueryMock.On("CoordinatorByAddress", mock.Anything, &profiletypes.QueryGetCoordinatorByAddressRequest{
		Address: testdata.Address,
	}).Return(nil, cosmoserror.ErrNotFound)
	network.profileQuery = profileQueryMock
	chainMock := newChainMockForPublish(false, false)
	launchID, campaignID, _, err := network.Publish(context.Background(), chainMock)
	require.Nil(t, err)
	require.Equal(t, TestLaunchID, launchID)
	require.Equal(t, TestCampaignID, campaignID)
	chainMock.AssertNumberOfCalls(t, "ChainID", 1)
	network.profileQuery.(*mocks.ProfileClient).AssertNumberOfCalls(t, "CoordinatorByAddress", 1)
	network.cosmos.(*mocks.CosmosClient).AssertNumberOfCalls(t, "BroadcastTx", 3)
	chainMock.AssertNumberOfCalls(t, "CacheBinary", 1)
}

func TestPublishFailedToFetchCoordinator(t *testing.T) {
	network := stubNetworkForPublish()
	profileQueryMock := new(mocks.ProfileClient)
	profileQueryMock.On("CoordinatorByAddress", mock.Anything, &profiletypes.QueryGetCoordinatorByAddressRequest{
		Address: testdata.Address,
	}).Return(nil, cosmoserror.ErrInternal)
	network.profileQuery = profileQueryMock
	chainMock := newChainMockForPublish(false, false)
	_, _, _, err := network.Publish(context.Background(), chainMock)
	require.NotNil(t, err)
	chainMock.AssertNumberOfCalls(t, "ChainID", 1)
	network.profileQuery.(*mocks.ProfileClient).AssertNumberOfCalls(t, "CoordinatorByAddress", 1)
	network.cosmos.(*mocks.CosmosClient).AssertNumberOfCalls(t, "BroadcastTx", 0)
	chainMock.AssertNumberOfCalls(t, "CacheBinary", 0)
}

func TestPublishFailedToReadChainID(t *testing.T) {
	network := stubNetworkForPublish()
	chainMock := newChainMockForPublish(false, true)
	_, _, _, err := network.Publish(context.Background(), chainMock)
	require.NotNil(t, err)
	chainMock.AssertNumberOfCalls(t, "ChainID", 1)
	network.profileQuery.(*mocks.ProfileClient).AssertNumberOfCalls(t, "CoordinatorByAddress", 0)
	network.cosmos.(*mocks.CosmosClient).AssertNumberOfCalls(t, "BroadcastTx", 0)
	chainMock.AssertNumberOfCalls(t, "CacheBinary", 0)
}

func TestPublishFailedToQueryCampaign(t *testing.T) {
	network := stubNetworkForPublish()
	campaignQueryMock := new(mocks.CampaignClient)
	campaignQueryMock.On("Campaign", mock.Anything, &campaigntypes.QueryGetCampaignRequest{
		CampaignID: TestCampaignID,
	}).Return(nil, cosmoserror.ErrNotFound)
	network.campaignQuery = campaignQueryMock
	chainMock := newChainMockForPublish(false, false)
	_, _, _, err := network.Publish(context.Background(), chainMock, WithCampaign(TestCampaignID))
	require.NotNil(t, err)
	chainMock.AssertNumberOfCalls(t, "ChainID", 1)
	network.profileQuery.(*mocks.ProfileClient).AssertNumberOfCalls(t, "CoordinatorByAddress", 1)
	network.cosmos.(*mocks.CosmosClient).AssertNumberOfCalls(t, "BroadcastTx", 0)
	chainMock.AssertNumberOfCalls(t, "CacheBinary", 0)
}

func TestPublishFailedToCreateCampaign(t *testing.T) {
	network := stubNetworkForPublish()
	networkClientMock := new(mocks.CosmosClient)
	networkClientMock.On("BroadcastTx",
		testdata.AccountName,
		mock.AnythingOfType("*types.MsgCreateCampaign"),
	).Return(testdata.NewResponse(&campaigntypes.MsgCreateCampaignResponse{
		CampaignID: TestCampaignID,
	}), errors.New("failed to create"))
	network.cosmos = networkClientMock
	chainMock := newChainMockForPublish(false, false)
	_, _, _, err := network.Publish(context.Background(), chainMock)
	require.NotNil(t, err)
	chainMock.AssertNumberOfCalls(t, "ChainID", 1)
	network.profileQuery.(*mocks.ProfileClient).AssertNumberOfCalls(t, "CoordinatorByAddress", 1)
	network.cosmos.(*mocks.CosmosClient).AssertNumberOfCalls(t, "BroadcastTx", 1)
	chainMock.AssertNumberOfCalls(t, "CacheBinary", 0)
}

func TestPublishFailedToCreateChain(t *testing.T) {
	network := stubNetworkForPublish()
	networkClientMock := new(mocks.CosmosClient)
	networkClientMock.On("BroadcastTx",
		testdata.AccountName,
		mock.AnythingOfType("*types.MsgCreateChain"),
	).Return(testdata.NewResponse(&launchtypes.MsgCreateChainResponse{
		LaunchID: TestLaunchID,
	}), errors.New("failed to create chain"))
	networkClientMock.On("BroadcastTx",
		testdata.AccountName,
		mock.AnythingOfType("*types.MsgCreateCampaign"),
	).Return(testdata.NewResponse(&campaigntypes.MsgCreateCampaignResponse{
		CampaignID: TestCampaignID,
	}), nil)
	network.cosmos = networkClientMock
	chainMock := newChainMockForPublish(false, false)
	_, _, _, err := network.Publish(context.Background(), chainMock)
	require.NotNil(t, err)
	chainMock.AssertNumberOfCalls(t, "ChainID", 1)
	network.profileQuery.(*mocks.ProfileClient).AssertNumberOfCalls(t, "CoordinatorByAddress", 1)
	network.cosmos.(*mocks.CosmosClient).AssertNumberOfCalls(t, "BroadcastTx", 2)
	chainMock.AssertNumberOfCalls(t, "CacheBinary", 0)
}

func TestPublishFailedToCacheBinary(t *testing.T) {
	network := stubNetworkForPublish()
	chainMock := newChainMockForPublish(true, false)
	_, _, _, err := network.Publish(context.Background(), chainMock)
	require.NotNil(t, err)
	chainMock.AssertNumberOfCalls(t, "ChainID", 1)
	network.profileQuery.(*mocks.ProfileClient).AssertNumberOfCalls(t, "CoordinatorByAddress", 1)
	network.cosmos.(*mocks.CosmosClient).AssertNumberOfCalls(t, "BroadcastTx", 2)
	chainMock.AssertNumberOfCalls(t, "CacheBinary", 1)
}
