package network

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	campaigntypes "github.com/tendermint/spn/x/campaign/types"
	launchtypes "github.com/tendermint/spn/x/launch/types"
	profiletypes "github.com/tendermint/spn/x/profile/types"
	"github.com/tendermint/starport/starport/pkg/cosmosaccount"
	"github.com/tendermint/starport/starport/pkg/cosmoserror"
	"github.com/tendermint/starport/starport/pkg/cosmosutil"
	"github.com/tendermint/starport/starport/services/network/mocks"
	"github.com/tendermint/starport/starport/services/network/networktypes"
	"github.com/tendermint/starport/starport/services/network/testutil"
)

func stubNetworkForPublish(account cosmosaccount.Account) Network {
	profileQueryMock := new(mocks.ProfileClient)
	profileQueryMock.On("CoordinatorByAddress", mock.Anything, &profiletypes.QueryGetCoordinatorByAddressRequest{
		Address: account.Address(networktypes.SPN),
	}).Return(nil, nil)
	campaignQueryMock := new(mocks.CampaignClient)
	campaignQueryMock.On("Campaign", mock.Anything, &campaigntypes.QueryGetCampaignRequest{
		CampaignID: testutil.TestCampaignID,
	}).Return(nil, nil)
	networkClientMock := new(mocks.CosmosClient)
	networkClientMock.On("BroadcastTx",
		account.Name,
		mock.AnythingOfType("*types.MsgCreateChain"),
	).Return(testutil.NewResponse(&launchtypes.MsgCreateChainResponse{
		LaunchID: testutil.TestLaunchID,
	}), nil)
	networkClientMock.On("BroadcastTx",
		account.Name,
		mock.AnythingOfType("*types.MsgCreateCampaign"),
	).Return(testutil.NewResponse(&campaigntypes.MsgCreateCampaignResponse{
		CampaignID: testutil.TestCampaignID,
	}), nil)
	networkClientMock.On("BroadcastTx",
		account.Name,
		mock.AnythingOfType("*types.MsgCreateCoordinator"),
	).Return(testutil.NewResponse(&profiletypes.MsgCreateCoordinatorResponse{}), nil)
	return Network{
		cosmos:        networkClientMock,
		account:       account,
		profileQuery:  profileQueryMock,
		campaignQuery: campaignQueryMock,
	}
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
	account, err := testutil.NewTestAccount(testutil.TestAccountName)
	assert.Nil(t, err)

	network := stubNetworkForPublish(account)
	chainMock := testutil.NewChainMock()
	launchID, campaignID, _, err := network.Publish(context.Background(), chainMock)
	require.Nil(t, err)
	require.Equal(t, testutil.TestLaunchID, launchID)
	require.Equal(t, testutil.TestCampaignID, campaignID)
	chainMock.AssertNumberOfCalls(t, "ChainID", 1)
	network.profileQuery.(*mocks.ProfileClient).AssertNumberOfCalls(t, "CoordinatorByAddress", 1)
	network.cosmos.(*mocks.CosmosClient).AssertNumberOfCalls(t, "BroadcastTx", 2)
	chainMock.AssertNumberOfCalls(t, "CacheBinary", 1)
}

func TestPublishWithCampaignID(t *testing.T) {
	account, err := testutil.NewTestAccount(testutil.TestAccountName)
	assert.Nil(t, err)

	network := stubNetworkForPublish(account)
	chainMock := testutil.NewChainMock()
	launchID, campaignID, _, err := network.Publish(context.Background(), chainMock, WithCampaign(testutil.TestCampaignID))
	require.Nil(t, err)
	require.Equal(t, testutil.TestLaunchID, launchID)
	require.Equal(t, testutil.TestCampaignID, campaignID)
	chainMock.AssertNumberOfCalls(t, "ChainID", 1)
	network.profileQuery.(*mocks.ProfileClient).AssertNumberOfCalls(t, "CoordinatorByAddress", 1)
	network.campaignQuery.(*mocks.CampaignClient).AssertNumberOfCalls(t, "Campaign", 1)
	network.cosmos.(*mocks.CosmosClient).AssertNumberOfCalls(t, "BroadcastTx", 1)
	chainMock.AssertNumberOfCalls(t, "CacheBinary", 1)
}

func TestPublishWithCustomGenesis(t *testing.T) {
	var customGenesisChainID = "test-custom-1"

	account, err := testutil.NewTestAccount(testutil.TestAccountName)
	assert.Nil(t, err)

	network := stubNetworkForPublish(account)
	chainMock := testutil.NewChainMock()
	gts := startGenesisTestServer(cosmosutil.ChainGenesis{ChainID: customGenesisChainID})
	defer gts.Close()
	launchID, campaignID, _, err := network.Publish(context.Background(), chainMock, WithCustomGenesis(gts.URL))
	require.Nil(t, err)
	require.Equal(t, testutil.TestLaunchID, launchID)
	require.Equal(t, testutil.TestCampaignID, campaignID)
	chainMock.AssertNumberOfCalls(t, "ChainID", 0)
	network.campaignQuery.(*mocks.CampaignClient).AssertNumberOfCalls(t, "Campaign", 0)
	network.profileQuery.(*mocks.ProfileClient).AssertNumberOfCalls(t, "CoordinatorByAddress", 1)
	network.cosmos.(*mocks.CosmosClient).AssertNumberOfCalls(t, "BroadcastTx", 2)
	chainMock.AssertNumberOfCalls(t, "CacheBinary", 1)
}

func TestPublishWithCustomChainID(t *testing.T) {
	var customGenesisChainID = "test-custom-1"

	account, err := testutil.NewTestAccount(testutil.TestAccountName)
	assert.Nil(t, err)

	network := stubNetworkForPublish(account)
	chainMock := testutil.NewChainMock()
	gts := startGenesisTestServer(cosmosutil.ChainGenesis{ChainID: customGenesisChainID})
	defer gts.Close()
	launchID, campaignID, _, err := network.Publish(context.Background(), chainMock, WithChainID(testutil.TestChainChainID))
	require.Nil(t, err)
	require.Equal(t, testutil.TestLaunchID, launchID)
	require.Equal(t, testutil.TestCampaignID, campaignID)
	chainMock.AssertNumberOfCalls(t, "ChainID", 0)
	network.profileQuery.(*mocks.ProfileClient).AssertNumberOfCalls(t, "CoordinatorByAddress", 1)
	network.cosmos.(*mocks.CosmosClient).AssertNumberOfCalls(t, "BroadcastTx", 2)
	chainMock.AssertNumberOfCalls(t, "CacheBinary", 1)
}

func TestPublishMainnet(t *testing.T) {
	account, err := testutil.NewTestAccount(testutil.TestAccountName)
	assert.Nil(t, err)

	network := stubNetworkForPublish(account)
	network.cosmos.(*mocks.CosmosClient).On(
		"BroadcastTx",
		account.Name,
		&campaigntypes.MsgInitializeMainnet{
			Coordinator:    account.Address(networktypes.SPN),
			CampaignID:     testutil.TestCampaignID,
			SourceURL:      testutil.TestChainSourceURL,
			SourceHash:     testutil.TestChainSourceHash,
			MainnetChainID: testutil.TestChainChainID,
		},
	).Return(testutil.NewResponse(&campaigntypes.MsgInitializeMainnetResponse{
		MainnetID: testutil.TestMainnetID,
	}), nil)
	chainMock := testutil.NewChainMock()
	launchID, campaignID, mainnetID, err := network.Publish(context.Background(), chainMock, Mainnet())
	require.Nil(t, err)
	require.Equal(t, testutil.TestLaunchID, launchID)
	require.Equal(t, testutil.TestCampaignID, campaignID)
	require.Equal(t, testutil.TestMainnetID, mainnetID)
	chainMock.AssertNumberOfCalls(t, "ChainID", 1)
	network.profileQuery.(*mocks.ProfileClient).AssertNumberOfCalls(t, "CoordinatorByAddress", 1)
	network.cosmos.(*mocks.CosmosClient).AssertNumberOfCalls(t, "BroadcastTx", 3)
	chainMock.AssertNumberOfCalls(t, "CacheBinary", 1)
}

func TestPublishWithCustomGenesisAndFailedToFetchIt(t *testing.T) {
	account, err := testutil.NewTestAccount(testutil.TestAccountName)
	assert.Nil(t, err)

	network := stubNetworkForPublish(account)
	chainMock := testutil.NewChainMock()
	gts := startNotFoundServer()
	defer gts.Close()
	_, _, _, err = network.Publish(context.Background(), chainMock, WithCustomGenesis(gts.URL))
	require.NotNil(t, err)
	chainMock.AssertNumberOfCalls(t, "ChainID", 0)
	network.campaignQuery.(*mocks.CampaignClient).AssertNumberOfCalls(t, "Campaign", 0)
	network.profileQuery.(*mocks.ProfileClient).AssertNumberOfCalls(t, "CoordinatorByAddress", 0)
	network.cosmos.(*mocks.CosmosClient).AssertNumberOfCalls(t, "BroadcastTx", 0)
	chainMock.AssertNumberOfCalls(t, "CacheBinary", 0)
}

func TestPublishFailedToInitializeMainnet(t *testing.T) {
	account, err := testutil.NewTestAccount(testutil.TestAccountName)
	assert.Nil(t, err)

	network := stubNetworkForPublish(account)
	network.cosmos.(*mocks.CosmosClient).On(
		"BroadcastTx",
		account.Name,
		&campaigntypes.MsgInitializeMainnet{
			Coordinator:    account.Address(networktypes.SPN),
			CampaignID:     testutil.TestCampaignID,
			SourceURL:      testutil.TestChainSourceURL,
			SourceHash:     testutil.TestChainSourceHash,
			MainnetChainID: testutil.TestChainChainID,
		},
	).Return(testutil.NewResponse(&campaigntypes.MsgInitializeMainnetResponse{
		MainnetID: testutil.TestMainnetID,
	}), errors.New("failed to initialize mainnet"))
	chainMock := testutil.NewChainMock()
	_, _, _, err = network.Publish(context.Background(), chainMock, Mainnet())
	require.NotNil(t, err)
	chainMock.AssertNumberOfCalls(t, "ChainID", 1)
	network.profileQuery.(*mocks.ProfileClient).AssertNumberOfCalls(t, "CoordinatorByAddress", 1)
	network.cosmos.(*mocks.CosmosClient).AssertNumberOfCalls(t, "BroadcastTx", 3)
	chainMock.AssertNumberOfCalls(t, "CacheBinary", 1)
}

func TestPublishWithCustomGenesisAndFailedToParseIt(t *testing.T) {
	var customGenesisChainID = "test-custom-1"

	account, err := testutil.NewTestAccount(testutil.TestAccountName)
	assert.Nil(t, err)

	network := stubNetworkForPublish(account)
	chainMock := testutil.NewChainMock()
	gts := startGenesisTestServer(cosmosutil.ChainGenesis{ChainID: customGenesisChainID})
	defer gts.Close()
	launchID, campaignID, _, err := network.Publish(context.Background(), chainMock, WithCustomGenesis(gts.URL))
	require.Nil(t, err)
	require.Equal(t, testutil.TestLaunchID, launchID)
	require.Equal(t, testutil.TestCampaignID, campaignID)
	chainMock.AssertNumberOfCalls(t, "ChainID", 0)
	network.campaignQuery.(*mocks.CampaignClient).AssertNumberOfCalls(t, "Campaign", 0)
	network.profileQuery.(*mocks.ProfileClient).AssertNumberOfCalls(t, "CoordinatorByAddress", 1)
	network.cosmos.(*mocks.CosmosClient).AssertNumberOfCalls(t, "BroadcastTx", 2)
	chainMock.AssertNumberOfCalls(t, "CacheBinary", 1)
}

func TestPublishWithCoordinatorCreation(t *testing.T) {
	account, err := testutil.NewTestAccount(testutil.TestAccountName)
	assert.Nil(t, err)

	network := stubNetworkForPublish(account)
	profileQueryMock := new(mocks.ProfileClient)
	profileQueryMock.On("CoordinatorByAddress", mock.Anything, &profiletypes.QueryGetCoordinatorByAddressRequest{
		Address: account.Address(networktypes.SPN),
	}).Return(nil, cosmoserror.ErrNotFound)
	network.profileQuery = profileQueryMock
	chainMock := testutil.NewChainMock()
	launchID, campaignID, _, err := network.Publish(context.Background(), chainMock)
	require.Nil(t, err)
	require.Equal(t, testutil.TestLaunchID, launchID)
	require.Equal(t, testutil.TestCampaignID, campaignID)
	chainMock.AssertNumberOfCalls(t, "ChainID", 1)
	network.profileQuery.(*mocks.ProfileClient).AssertNumberOfCalls(t, "CoordinatorByAddress", 1)
	network.cosmos.(*mocks.CosmosClient).AssertNumberOfCalls(t, "BroadcastTx", 3)
	chainMock.AssertNumberOfCalls(t, "CacheBinary", 1)
}

func TestPublishFailedToFetchCoordinator(t *testing.T) {
	account, err := testutil.NewTestAccount(testutil.TestAccountName)
	assert.Nil(t, err)

	network := stubNetworkForPublish(account)
	profileQueryMock := new(mocks.ProfileClient)
	profileQueryMock.On("CoordinatorByAddress", mock.Anything, &profiletypes.QueryGetCoordinatorByAddressRequest{
		Address: account.Address(networktypes.SPN),
	}).Return(nil, cosmoserror.ErrInternal)
	network.profileQuery = profileQueryMock
	chainMock := testutil.NewChainMock()
	_, _, _, err = network.Publish(context.Background(), chainMock)
	require.NotNil(t, err)
	chainMock.AssertNumberOfCalls(t, "ChainID", 1)
	network.profileQuery.(*mocks.ProfileClient).AssertNumberOfCalls(t, "CoordinatorByAddress", 1)
	network.cosmos.(*mocks.CosmosClient).AssertNumberOfCalls(t, "BroadcastTx", 0)
	chainMock.AssertNumberOfCalls(t, "CacheBinary", 0)
}

func TestPublishFailedToReadChainID(t *testing.T) {
	account, err := testutil.NewTestAccount(testutil.TestAccountName)
	assert.Nil(t, err)

	network := stubNetworkForPublish(account)
	chainMock := testutil.NewChainMock(testutil.WithChainIDFail())
	_, _, _, err = network.Publish(context.Background(), chainMock)
	require.NotNil(t, err)
	chainMock.AssertNumberOfCalls(t, "ChainID", 1)
	network.profileQuery.(*mocks.ProfileClient).AssertNumberOfCalls(t, "CoordinatorByAddress", 0)
	network.cosmos.(*mocks.CosmosClient).AssertNumberOfCalls(t, "BroadcastTx", 0)
	chainMock.AssertNumberOfCalls(t, "CacheBinary", 0)
}

func TestPublishFailedToQueryCampaign(t *testing.T) {
	account, err := testutil.NewTestAccount(testutil.TestAccountName)
	assert.Nil(t, err)

	network := stubNetworkForPublish(account)
	campaignQueryMock := new(mocks.CampaignClient)
	campaignQueryMock.On("Campaign", mock.Anything, &campaigntypes.QueryGetCampaignRequest{
		CampaignID: testutil.TestCampaignID,
	}).Return(nil, cosmoserror.ErrNotFound)
	network.campaignQuery = campaignQueryMock
	chainMock := testutil.NewChainMock()
	_, _, _, err = network.Publish(context.Background(), chainMock, WithCampaign(testutil.TestCampaignID))
	require.NotNil(t, err)
	chainMock.AssertNumberOfCalls(t, "ChainID", 1)
	network.profileQuery.(*mocks.ProfileClient).AssertNumberOfCalls(t, "CoordinatorByAddress", 1)
	network.cosmos.(*mocks.CosmosClient).AssertNumberOfCalls(t, "BroadcastTx", 0)
	chainMock.AssertNumberOfCalls(t, "CacheBinary", 0)
}

func TestPublishFailedToCreateCampaign(t *testing.T) {
	account, err := testutil.NewTestAccount(testutil.TestAccountName)
	assert.Nil(t, err)

	network := stubNetworkForPublish(account)
	networkClientMock := new(mocks.CosmosClient)
	networkClientMock.On("BroadcastTx",
		account.Name,
		mock.AnythingOfType("*types.MsgCreateCampaign"),
	).Return(testutil.NewResponse(&campaigntypes.MsgCreateCampaignResponse{
		CampaignID: testutil.TestCampaignID,
	}), errors.New("failed to create"))
	network.cosmos = networkClientMock
	chainMock := testutil.NewChainMock()
	_, _, _, err = network.Publish(context.Background(), chainMock)
	require.NotNil(t, err)
	chainMock.AssertNumberOfCalls(t, "ChainID", 1)
	network.profileQuery.(*mocks.ProfileClient).AssertNumberOfCalls(t, "CoordinatorByAddress", 1)
	network.cosmos.(*mocks.CosmosClient).AssertNumberOfCalls(t, "BroadcastTx", 1)
	chainMock.AssertNumberOfCalls(t, "CacheBinary", 0)
}

func TestPublishFailedToCreateChain(t *testing.T) {
	account, err := testutil.NewTestAccount(testutil.TestAccountName)
	assert.Nil(t, err)

	network := stubNetworkForPublish(account)
	networkClientMock := new(mocks.CosmosClient)
	networkClientMock.On("BroadcastTx",
		account.Name,
		mock.AnythingOfType("*types.MsgCreateChain"),
	).Return(testutil.NewResponse(&launchtypes.MsgCreateChainResponse{
		LaunchID: testutil.TestLaunchID,
	}), errors.New("failed to create chain"))
	networkClientMock.On("BroadcastTx",
		account.Name,
		mock.AnythingOfType("*types.MsgCreateCampaign"),
	).Return(testutil.NewResponse(&campaigntypes.MsgCreateCampaignResponse{
		CampaignID: testutil.TestCampaignID,
	}), nil)
	network.cosmos = networkClientMock
	chainMock := testutil.NewChainMock()
	_, _, _, err = network.Publish(context.Background(), chainMock)
	require.NotNil(t, err)
	chainMock.AssertNumberOfCalls(t, "ChainID", 1)
	network.profileQuery.(*mocks.ProfileClient).AssertNumberOfCalls(t, "CoordinatorByAddress", 1)
	network.cosmos.(*mocks.CosmosClient).AssertNumberOfCalls(t, "BroadcastTx", 2)
	chainMock.AssertNumberOfCalls(t, "CacheBinary", 0)
}

func TestPublishFailedToCacheBinary(t *testing.T) {
	account, err := testutil.NewTestAccount(testutil.TestAccountName)
	assert.Nil(t, err)

	network := stubNetworkForPublish(account)
	chainMock := testutil.NewChainMock(testutil.WithCacheBinaryFail())
	_, _, _, err = network.Publish(context.Background(), chainMock)
	require.NotNil(t, err)
	chainMock.AssertNumberOfCalls(t, "ChainID", 1)
	network.profileQuery.(*mocks.ProfileClient).AssertNumberOfCalls(t, "CoordinatorByAddress", 1)
	network.cosmos.(*mocks.CosmosClient).AssertNumberOfCalls(t, "BroadcastTx", 2)
	chainMock.AssertNumberOfCalls(t, "CacheBinary", 1)
}
