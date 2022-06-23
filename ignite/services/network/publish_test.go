package network

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	campaigntypes "github.com/tendermint/spn/x/campaign/types"
	launchtypes "github.com/tendermint/spn/x/launch/types"
	profiletypes "github.com/tendermint/spn/x/profile/types"

	"github.com/ignite/cli/ignite/pkg/cosmoserror"
	"github.com/ignite/cli/ignite/pkg/cosmosutil"
	"github.com/ignite/cli/ignite/services/network/networktypes"
	"github.com/ignite/cli/ignite/services/network/testutil"
)

func startGenesisTestServer(genesis cosmosutil.ChainGenesis) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		encodedGenesis, _ := json.Marshal(genesis)
		w.Write(encodedGenesis)
	}))
}

func startInvalidJSONServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("invalid json"))
	}))
}

func TestPublish(t *testing.T) {
	t.Run("publish chain", func(t *testing.T) {
		var (
			account        = testutil.NewTestAccount(t, testutil.TestAccountName)
			suite, network = newSuite(account)
		)

		suite.ProfileQueryMock.
			On(
				"CoordinatorByAddress",
				context.Background(),
				&profiletypes.QueryGetCoordinatorByAddressRequest{
					Address: account.Address(networktypes.SPN),
				},
			).
			Return(nil, nil).
			Once()
		suite.CosmosClientMock.
			On(
				"BroadcastTx",
				account.Name,
				&launchtypes.MsgCreateChain{
					Coordinator:    account.Address(networktypes.SPN),
					GenesisChainID: testutil.ChainID,
					SourceURL:      testutil.ChainSourceURL,
					SourceHash:     testutil.ChainSourceHash,
					GenesisURL:     "",
					GenesisHash:    "",
					HasCampaign:    true,
					CampaignID:     testutil.CampaignID,
				},
			).
			Return(testutil.NewResponse(&launchtypes.MsgCreateChainResponse{
				LaunchID: testutil.LaunchID,
			}), nil).
			Once()
		suite.CosmosClientMock.
			On(
				"BroadcastTx",
				account.Name,
				&campaigntypes.MsgCreateCampaign{
					Coordinator:  account.Address(networktypes.SPN),
					CampaignName: testutil.ChainName,
					Metadata:     []byte{},
				},
			).
			Return(testutil.NewResponse(&campaigntypes.MsgCreateCampaignResponse{
				CampaignID: testutil.CampaignID,
			}), nil).
			Once()
		suite.ChainMock.On("SourceHash").Return(testutil.ChainSourceHash).Once()
		suite.ChainMock.On("SourceURL").Return(testutil.ChainSourceURL).Once()
		suite.ChainMock.On("Name").Return(testutil.ChainName).Once()
		suite.ChainMock.On("ChainID").Return(testutil.ChainID, nil).Once()
		suite.ChainMock.On("CacheBinary", testutil.LaunchID).Return(nil).Once()

		launchID, campaignID, publishError := network.Publish(context.Background(), suite.ChainMock)
		require.NoError(t, publishError)
		require.Equal(t, testutil.LaunchID, launchID)
		require.Equal(t, testutil.CampaignID, campaignID)
		suite.AssertAllMocks(t)
	})

	t.Run("publish chain with shares", func(t *testing.T) {
		var (
			account        = testutil.NewTestAccount(t, testutil.TestAccountName)
			suite, network = newSuite(account)
		)

		suite.ProfileQueryMock.
			On(
				"CoordinatorByAddress",
				context.Background(),
				&profiletypes.QueryGetCoordinatorByAddressRequest{
					Address: account.Address(networktypes.SPN),
				},
			).
			Return(nil, nil).
			Once()
		suite.CampaignQueryMock.
			On(
				"TotalShares",
				context.Background(),
				&campaigntypes.QueryTotalSharesRequest{},
			).
			Return(&campaigntypes.QueryTotalSharesResponse{
				TotalShares: 100000,
			}, nil).
			Once()
		suite.CosmosClientMock.
			On(
				"BroadcastTx",
				account.Name,
				campaigntypes.NewMsgMintVouchers(
					account.Address(networktypes.SPN),
					testutil.CampaignID,
					campaigntypes.NewSharesFromCoins(sdk.NewCoins(sdk.NewInt64Coin("foo", 2000), sdk.NewInt64Coin("staking", 50000))),
				),
			).
			Return(testutil.NewResponse(&campaigntypes.MsgMintVouchersResponse{}), nil).
			Once()
		suite.CosmosClientMock.
			On(
				"BroadcastTx",
				account.Name,
				&launchtypes.MsgCreateChain{
					Coordinator:    account.Address(networktypes.SPN),
					GenesisChainID: testutil.ChainID,
					SourceURL:      testutil.ChainSourceURL,
					SourceHash:     testutil.ChainSourceHash,
					GenesisURL:     "",
					GenesisHash:    "",
					HasCampaign:    true,
					CampaignID:     testutil.CampaignID,
				},
			).
			Return(testutil.NewResponse(&launchtypes.MsgCreateChainResponse{
				LaunchID: testutil.LaunchID,
			}), nil).
			Once()
		suite.CosmosClientMock.
			On(
				"BroadcastTx",
				account.Name,
				&campaigntypes.MsgCreateCampaign{
					Coordinator:  account.Address(networktypes.SPN),
					CampaignName: testutil.ChainName,
					Metadata:     []byte{},
				},
			).
			Return(testutil.NewResponse(&campaigntypes.MsgCreateCampaignResponse{
				CampaignID: testutil.CampaignID,
			}), nil).
			Once()
		suite.ChainMock.On("SourceHash").Return(testutil.ChainSourceHash).Once()
		suite.ChainMock.On("SourceURL").Return(testutil.ChainSourceURL).Once()
		suite.ChainMock.On("Name").Return(testutil.ChainName).Once()
		suite.ChainMock.On("ChainID").Return(testutil.ChainID, nil).Once()
		suite.ChainMock.On("CacheBinary", testutil.LaunchID).Return(nil).Once()

		launchID, campaignID, publishError := network.Publish(context.Background(), suite.ChainMock,
			WithPercentageShares([]SharePercent{
				SampleSharePercent(t, "foo", 2, 100),
				SampleSharePercent(t, "staking", 50, 100),
			}),
		)
		require.NoError(t, publishError)
		require.Equal(t, testutil.LaunchID, launchID)
		require.Equal(t, testutil.CampaignID, campaignID)
		suite.AssertAllMocks(t)
	})

	t.Run("publish chain with pre created campaign", func(t *testing.T) {
		var (
			account        = testutil.NewTestAccount(t, testutil.TestAccountName)
			suite, network = newSuite(account)
		)

		suite.ProfileQueryMock.
			On(
				"CoordinatorByAddress",
				context.Background(),
				&profiletypes.QueryGetCoordinatorByAddressRequest{
					Address: account.Address(networktypes.SPN),
				},
			).
			Return(nil, nil).
			Once()
		suite.CampaignQueryMock.
			On(
				"Campaign",
				context.Background(),
				&campaigntypes.QueryGetCampaignRequest{
					CampaignID: testutil.CampaignID,
				},
			).
			Return(nil, nil).
			Once()
		suite.CosmosClientMock.
			On(
				"BroadcastTx",
				account.Name,
				&launchtypes.MsgCreateChain{
					Coordinator:    account.Address(networktypes.SPN),
					GenesisChainID: testutil.ChainID,
					SourceURL:      testutil.ChainSourceURL,
					SourceHash:     testutil.ChainSourceHash,
					GenesisURL:     "",
					GenesisHash:    "",
					HasCampaign:    true,
					CampaignID:     testutil.CampaignID,
				},
			).
			Return(testutil.NewResponse(&launchtypes.MsgCreateChainResponse{
				LaunchID: testutil.LaunchID,
			}), nil).
			Once()
		suite.ChainMock.On("SourceHash").Return(testutil.ChainSourceHash).Once()
		suite.ChainMock.On("SourceURL").Return(testutil.ChainSourceURL).Once()
		suite.ChainMock.On("ChainID").Return(testutil.ChainID, nil).Once()
		suite.ChainMock.On("CacheBinary", testutil.LaunchID).Return(nil).Once()

		launchID, campaignID, publishError := network.Publish(context.Background(), suite.ChainMock, WithCampaign(testutil.CampaignID))
		require.NoError(t, publishError)
		require.Equal(t, testutil.LaunchID, launchID)
		require.Equal(t, testutil.CampaignID, campaignID)
		suite.AssertAllMocks(t)
	})

	t.Run("publish chain with custom genesis", func(t *testing.T) {
		var (
			account              = testutil.NewTestAccount(t, testutil.TestAccountName)
			customGenesisChainID = "test-custom-1"
			customGenesisHash    = "72a80a32e33513cd74423354502cef035e96b0bff59c754646b453b201d12d07"
			gts                  = startGenesisTestServer(cosmosutil.ChainGenesis{ChainID: customGenesisChainID})
			suite, network       = newSuite(account)
		)
		defer gts.Close()

		suite.ProfileQueryMock.
			On(
				"CoordinatorByAddress",
				context.Background(),
				&profiletypes.QueryGetCoordinatorByAddressRequest{
					Address: account.Address(networktypes.SPN),
				},
			).
			Return(nil, nil).
			Once()
		suite.CosmosClientMock.
			On(
				"BroadcastTx",
				account.Name,
				&launchtypes.MsgCreateChain{
					Coordinator:    account.Address(networktypes.SPN),
					GenesisChainID: customGenesisChainID,
					SourceURL:      testutil.ChainSourceURL,
					SourceHash:     testutil.ChainSourceHash,
					GenesisURL:     gts.URL,
					GenesisHash:    customGenesisHash,
					HasCampaign:    true,
					CampaignID:     testutil.CampaignID,
				},
			).
			Return(testutil.NewResponse(&launchtypes.MsgCreateChainResponse{
				LaunchID: testutil.LaunchID,
			}), nil).
			Once()
		suite.CosmosClientMock.
			On(
				"BroadcastTx",
				account.Name,
				&campaigntypes.MsgCreateCampaign{
					Coordinator:  account.Address(networktypes.SPN),
					CampaignName: testutil.ChainName,
					Metadata:     []byte{},
				},
			).
			Return(testutil.NewResponse(&campaigntypes.MsgCreateCampaignResponse{
				CampaignID: testutil.CampaignID,
			}), nil).
			Once()
		suite.ChainMock.On("SourceHash").Return(testutil.ChainSourceHash).Once()
		suite.ChainMock.On("SourceURL").Return(testutil.ChainSourceURL).Once()
		suite.ChainMock.On("Name").Return(testutil.ChainName).Once()
		suite.ChainMock.On("CacheBinary", testutil.LaunchID).Return(nil).Once()

		launchID, campaignID, publishError := network.Publish(context.Background(), suite.ChainMock, WithCustomGenesis(gts.URL))
		require.NoError(t, publishError)
		require.Equal(t, testutil.LaunchID, launchID)
		require.Equal(t, testutil.CampaignID, campaignID)
		suite.AssertAllMocks(t)
	})

	t.Run("publish chain with custom chain id", func(t *testing.T) {
		var (
			account        = testutil.NewTestAccount(t, testutil.TestAccountName)
			suite, network = newSuite(account)
		)

		suite.ProfileQueryMock.
			On(
				"CoordinatorByAddress",
				context.Background(),
				&profiletypes.QueryGetCoordinatorByAddressRequest{
					Address: account.Address(networktypes.SPN),
				},
			).
			Return(nil, nil).
			Once()
		suite.CosmosClientMock.
			On(
				"BroadcastTx",
				account.Name,
				&launchtypes.MsgCreateChain{
					Coordinator:    account.Address(networktypes.SPN),
					GenesisChainID: testutil.ChainID,
					SourceURL:      testutil.ChainSourceURL,
					SourceHash:     testutil.ChainSourceHash,
					GenesisURL:     "",
					GenesisHash:    "",
					HasCampaign:    true,
					CampaignID:     testutil.CampaignID,
				},
			).
			Return(testutil.NewResponse(&launchtypes.MsgCreateChainResponse{
				LaunchID: testutil.LaunchID,
			}), nil).
			Once()
		suite.CosmosClientMock.
			On(
				"BroadcastTx",
				account.Name,
				&campaigntypes.MsgCreateCampaign{
					Coordinator:  account.Address(networktypes.SPN),
					CampaignName: testutil.ChainName,
					Metadata:     []byte{},
				},
			).
			Return(testutil.NewResponse(&campaigntypes.MsgCreateCampaignResponse{
				CampaignID: testutil.CampaignID,
			}), nil).
			Once()
		suite.ChainMock.On("SourceHash").Return(testutil.ChainSourceHash).Once()
		suite.ChainMock.On("SourceURL").Return(testutil.ChainSourceURL).Once()
		suite.ChainMock.On("Name").Return(testutil.ChainName).Once()
		suite.ChainMock.On("CacheBinary", testutil.LaunchID).Return(nil).Once()

		launchID, campaignID, publishError := network.Publish(context.Background(), suite.ChainMock, WithChainID(testutil.ChainID))
		require.NoError(t, publishError)
		require.Equal(t, testutil.LaunchID, launchID)
		require.Equal(t, testutil.CampaignID, campaignID)
		suite.AssertAllMocks(t)
	})

	t.Run("publish chain with mainnet", func(t *testing.T) {
		var (
			account              = testutil.NewTestAccount(t, testutil.TestAccountName)
			customGenesisChainID = "test-custom-1"
			gts                  = startGenesisTestServer(cosmosutil.ChainGenesis{ChainID: customGenesisChainID})
			suite, network       = newSuite(account)
		)
		defer gts.Close()

		suite.ProfileQueryMock.
			On(
				"CoordinatorByAddress",
				context.Background(),
				&profiletypes.QueryGetCoordinatorByAddressRequest{
					Address: account.Address(networktypes.SPN),
				},
			).
			Return(nil, nil).
			Once()
		suite.CosmosClientMock.
			On(
				"BroadcastTx",
				account.Name,
				&campaigntypes.MsgCreateCampaign{
					Coordinator:  account.Address(networktypes.SPN),
					CampaignName: testutil.ChainName,
					Metadata:     []byte{},
				},
			).
			Return(testutil.NewResponse(&campaigntypes.MsgCreateCampaignResponse{
				CampaignID: testutil.CampaignID,
			}), nil).
			Once()
		suite.CosmosClientMock.
			On(
				"BroadcastTx",
				account.Name,
				&campaigntypes.MsgInitializeMainnet{
					Coordinator:    account.Address(networktypes.SPN),
					CampaignID:     testutil.CampaignID,
					SourceURL:      testutil.ChainSourceURL,
					SourceHash:     testutil.ChainSourceHash,
					MainnetChainID: testutil.ChainID,
				},
			).
			Return(testutil.NewResponse(&campaigntypes.MsgInitializeMainnetResponse{
				MainnetID: testutil.MainnetID,
			}), nil).
			Once()
		suite.ChainMock.On("SourceHash").Return(testutil.ChainSourceHash).Once()
		suite.ChainMock.On("SourceURL").Return(testutil.ChainSourceURL).Once()
		suite.ChainMock.On("Name").Return(testutil.ChainName).Once()
		suite.ChainMock.On("ChainID").Return(testutil.ChainID, nil).Once()
		suite.ChainMock.On("CacheBinary", testutil.LaunchID).Return(nil).Once()

		launchID, campaignID, publishError := network.Publish(context.Background(), suite.ChainMock, Mainnet())
		require.NoError(t, publishError)
		require.Equal(t, testutil.LaunchID, launchID)
		require.Equal(t, testutil.CampaignID, campaignID)
		suite.AssertAllMocks(t)
	})

	t.Run("failed to publish chain with mainnet, failed to initialize mainnet", func(t *testing.T) {
		var (
			account        = testutil.NewTestAccount(t, testutil.TestAccountName)
			suite, network = newSuite(account)
			expectedError  = errors.New("failed to initialize mainnet")
		)

		suite.ProfileQueryMock.
			On(
				"CoordinatorByAddress",
				context.Background(),
				&profiletypes.QueryGetCoordinatorByAddressRequest{
					Address: account.Address(networktypes.SPN),
				},
			).
			Return(nil, nil).
			Once()
		suite.CosmosClientMock.
			On(
				"BroadcastTx",
				account.Name,
				&campaigntypes.MsgCreateCampaign{
					Coordinator:  account.Address(networktypes.SPN),
					CampaignName: testutil.ChainName,
					Metadata:     []byte{},
				},
			).
			Return(testutil.NewResponse(&campaigntypes.MsgCreateCampaignResponse{
				CampaignID: testutil.CampaignID,
			}), nil).
			Once()
		suite.CosmosClientMock.
			On(
				"BroadcastTx",
				account.Name,
				&campaigntypes.MsgInitializeMainnet{
					Coordinator:    account.Address(networktypes.SPN),
					CampaignID:     testutil.CampaignID,
					SourceURL:      testutil.ChainSourceURL,
					SourceHash:     testutil.ChainSourceHash,
					MainnetChainID: testutil.ChainID,
				},
			).
			Return(testutil.NewResponse(&campaigntypes.MsgInitializeMainnetResponse{
				MainnetID: testutil.MainnetID,
			}), expectedError).
			Once()
		suite.ChainMock.On("SourceHash").Return(testutil.ChainSourceHash).Once()
		suite.ChainMock.On("SourceURL").Return(testutil.ChainSourceURL).Once()
		suite.ChainMock.On("Name").Return(testutil.ChainName).Once()
		suite.ChainMock.On("ChainID").Return(testutil.ChainID, nil).Once()

		_, _, publishError := network.Publish(context.Background(), suite.ChainMock, Mainnet())
		require.Error(t, publishError)
		require.Equal(t, expectedError, publishError)
		suite.AssertAllMocks(t)
	})

	t.Run("failed to publish chain with custom genesis, failed to parse custom genesis", func(t *testing.T) {
		var (
			account        = testutil.NewTestAccount(t, testutil.TestAccountName)
			suite, network = newSuite(account)
			gts            = startInvalidJSONServer()
			expectedError  = errors.New("cannot unmarshal the chain genesis file: invalid character 'i' looking for beginning of value")
		)
		defer gts.Close()

		_, _, publishError := network.Publish(context.Background(), suite.ChainMock, WithCustomGenesis(gts.URL))
		require.Error(t, publishError)
		require.Equal(t, expectedError.Error(), publishError.Error())
		suite.AssertAllMocks(t)
	})

	t.Run("publish chain with coordinator creation", func(t *testing.T) {
		var (
			account        = testutil.NewTestAccount(t, testutil.TestAccountName)
			suite, network = newSuite(account)
		)

		suite.ProfileQueryMock.
			On("CoordinatorByAddress", mock.Anything, &profiletypes.QueryGetCoordinatorByAddressRequest{
				Address: account.Address(networktypes.SPN),
			}).
			Return(nil, cosmoserror.ErrNotFound).
			Once()
		suite.CosmosClientMock.
			On(
				"BroadcastTx",
				account.Name,
				&launchtypes.MsgCreateChain{
					Coordinator:    account.Address(networktypes.SPN),
					GenesisChainID: testutil.ChainID,
					SourceURL:      testutil.ChainSourceURL,
					SourceHash:     testutil.ChainSourceHash,
					GenesisURL:     "",
					GenesisHash:    "",
					HasCampaign:    true,
					CampaignID:     testutil.CampaignID,
				},
			).
			Return(testutil.NewResponse(&launchtypes.MsgCreateChainResponse{
				LaunchID: testutil.LaunchID,
			}), nil).
			Once()
		suite.CosmosClientMock.
			On(
				"BroadcastTx",
				account.Name,
				&campaigntypes.MsgCreateCampaign{
					Coordinator:  account.Address(networktypes.SPN),
					CampaignName: testutil.ChainName,
					Metadata:     []byte{},
				},
			).
			Return(testutil.NewResponse(&campaigntypes.MsgCreateCampaignResponse{
				CampaignID: testutil.CampaignID,
			}), nil).
			Once()
		suite.CosmosClientMock.
			On(
				"BroadcastTx",
				account.Name,
				&profiletypes.MsgCreateCoordinator{
					Address: account.Address(networktypes.SPN),
				},
			).
			Return(testutil.NewResponse(&profiletypes.MsgCreateCoordinatorResponse{}), nil).
			Once()
		suite.ChainMock.On("SourceHash").Return(testutil.ChainSourceHash).Once()
		suite.ChainMock.On("SourceURL").Return(testutil.ChainSourceURL).Once()
		suite.ChainMock.On("Name").Return(testutil.ChainName).Once()
		suite.ChainMock.On("ChainID").Return(testutil.ChainID, nil).Once()
		suite.ChainMock.On("CacheBinary", testutil.LaunchID).Return(nil).Once()

		launchID, campaignID, publishError := network.Publish(context.Background(), suite.ChainMock)
		require.NoError(t, publishError)
		require.Equal(t, testutil.LaunchID, launchID)
		require.Equal(t, testutil.CampaignID, campaignID)
		suite.AssertAllMocks(t)
	})

	t.Run("failed to publish chain, failed to fetch coordinator profile", func(t *testing.T) {
		var (
			account        = testutil.NewTestAccount(t, testutil.TestAccountName)
			suite, network = newSuite(account)
			expectedError  = errors.New("failed to fetch coordinator")
		)

		suite.ProfileQueryMock.
			On("CoordinatorByAddress", mock.Anything, &profiletypes.QueryGetCoordinatorByAddressRequest{
				Address: account.Address(networktypes.SPN),
			}).
			Return(nil, expectedError).
			Once()
		suite.ChainMock.On("ChainID").Return(testutil.ChainID, nil).Once()

		_, _, publishError := network.Publish(context.Background(), suite.ChainMock)
		require.Error(t, publishError)
		require.Equal(t, expectedError, publishError)
		suite.AssertAllMocks(t)
	})

	t.Run("failed to publish chain, failed to read chain id", func(t *testing.T) {
		var (
			account        = testutil.NewTestAccount(t, testutil.TestAccountName)
			suite, network = newSuite(account)
			expectedError  = errors.New("failed to get chainID")
		)

		suite.ChainMock.
			On("ChainID").
			Return("", expectedError).
			Once()

		_, _, publishError := network.Publish(context.Background(), suite.ChainMock)
		require.Error(t, publishError)
		require.Equal(t, expectedError, publishError)
		suite.AssertAllMocks(t)
	})

	t.Run("failed to publish chain, failed to fetch existed campaign", func(t *testing.T) {
		var (
			account        = testutil.NewTestAccount(t, testutil.TestAccountName)
			suite, network = newSuite(account)
		)

		suite.ProfileQueryMock.
			On(
				"CoordinatorByAddress",
				context.Background(),
				&profiletypes.QueryGetCoordinatorByAddressRequest{
					Address: account.Address(networktypes.SPN),
				},
			).
			Return(nil, nil).
			Once()
		suite.CampaignQueryMock.
			On("Campaign", mock.Anything, &campaigntypes.QueryGetCampaignRequest{
				CampaignID: testutil.CampaignID,
			}).
			Return(nil, cosmoserror.ErrNotFound).
			Once()
		suite.ChainMock.On("ChainID").Return(testutil.ChainID, nil).Once()

		_, _, publishError := network.Publish(context.Background(), suite.ChainMock, WithCampaign(testutil.CampaignID))
		require.Error(t, publishError)
		require.Equal(t, cosmoserror.ErrNotFound, publishError)
		suite.AssertAllMocks(t)
	})

	t.Run("failed to publish chain, failed to create campaign", func(t *testing.T) {
		var (
			account        = testutil.NewTestAccount(t, testutil.TestAccountName)
			suite, network = newSuite(account)
			expectedError  = errors.New("failed to create")
		)

		suite.ProfileQueryMock.
			On(
				"CoordinatorByAddress",
				context.Background(),
				&profiletypes.QueryGetCoordinatorByAddressRequest{
					Address: account.Address(networktypes.SPN),
				},
			).
			Return(nil, nil).
			Once()
		suite.CosmosClientMock.
			On(
				"BroadcastTx",
				account.Name,
				&campaigntypes.MsgCreateCampaign{
					Coordinator:  account.Address(networktypes.SPN),
					CampaignName: testutil.ChainName,
					Metadata:     []byte{},
				},
			).
			Return(testutil.NewResponse(&campaigntypes.MsgCreateCampaignResponse{
				CampaignID: testutil.CampaignID,
			}), expectedError).
			Once()
		suite.ChainMock.On("Name").Return(testutil.ChainName).Once()
		suite.ChainMock.On("ChainID").Return(testutil.ChainID, nil).Once()

		_, _, publishError := network.Publish(context.Background(), suite.ChainMock)
		require.Error(t, publishError)
		require.Equal(t, expectedError, publishError)
		suite.AssertAllMocks(t)
	})

	t.Run("failed to publish chain, failed to create chain", func(t *testing.T) {
		var (
			account        = testutil.NewTestAccount(t, testutil.TestAccountName)
			suite, network = newSuite(account)
			expectedError  = errors.New("failed to create chain")
		)

		suite.ProfileQueryMock.
			On(
				"CoordinatorByAddress",
				context.Background(),
				&profiletypes.QueryGetCoordinatorByAddressRequest{
					Address: account.Address(networktypes.SPN),
				},
			).
			Return(nil, nil).
			Once()
		suite.CosmosClientMock.
			On(
				"BroadcastTx",
				account.Name,
				&launchtypes.MsgCreateChain{
					Coordinator:    account.Address(networktypes.SPN),
					GenesisChainID: testutil.ChainID,
					SourceURL:      testutil.ChainSourceURL,
					SourceHash:     testutil.ChainSourceHash,
					GenesisURL:     "",
					GenesisHash:    "",
					HasCampaign:    true,
					CampaignID:     testutil.CampaignID,
				},
			).
			Return(testutil.NewResponse(&launchtypes.MsgCreateChainResponse{
				LaunchID: testutil.LaunchID,
			}), expectedError).
			Once()
		suite.CosmosClientMock.
			On(
				"BroadcastTx",
				account.Name,
				&campaigntypes.MsgCreateCampaign{
					Coordinator:  account.Address(networktypes.SPN),
					CampaignName: testutil.ChainName,
					Metadata:     []byte{},
				},
			).
			Return(testutil.NewResponse(&campaigntypes.MsgCreateCampaignResponse{
				CampaignID: testutil.CampaignID,
			}), nil).
			Once()
		suite.ChainMock.On("SourceHash").Return(testutil.ChainSourceHash).Once()
		suite.ChainMock.On("SourceURL").Return(testutil.ChainSourceURL).Once()
		suite.ChainMock.On("Name").Return(testutil.ChainName).Once()
		suite.ChainMock.On("ChainID").Return(testutil.ChainID, nil).Once()

		_, _, publishError := network.Publish(context.Background(), suite.ChainMock)
		require.Error(t, publishError)
		require.Equal(t, expectedError, publishError)
		suite.AssertAllMocks(t)
	})

	t.Run("failed to publish chain, failed to cache binary", func(t *testing.T) {
		var (
			account        = testutil.NewTestAccount(t, testutil.TestAccountName)
			suite, network = newSuite(account)
			expectedError  = errors.New("failed to cache binary")
		)

		suite.ProfileQueryMock.
			On(
				"CoordinatorByAddress",
				context.Background(),
				&profiletypes.QueryGetCoordinatorByAddressRequest{
					Address: account.Address(networktypes.SPN),
				},
			).
			Return(nil, nil).
			Once()
		suite.CosmosClientMock.
			On(
				"BroadcastTx",
				account.Name,
				&launchtypes.MsgCreateChain{
					Coordinator:    account.Address(networktypes.SPN),
					GenesisChainID: testutil.ChainID,
					SourceURL:      testutil.ChainSourceURL,
					SourceHash:     testutil.ChainSourceHash,
					GenesisURL:     "",
					GenesisHash:    "",
					HasCampaign:    true,
					CampaignID:     testutil.CampaignID,
				},
			).
			Return(testutil.NewResponse(&launchtypes.MsgCreateChainResponse{
				LaunchID: testutil.LaunchID,
			}), nil).
			Once()
		suite.CosmosClientMock.
			On(
				"BroadcastTx",
				account.Name,
				&campaigntypes.MsgCreateCampaign{
					Coordinator:  account.Address(networktypes.SPN),
					CampaignName: testutil.ChainName,
					Metadata:     []byte{},
				},
			).
			Return(testutil.NewResponse(&campaigntypes.MsgCreateCampaignResponse{
				CampaignID: testutil.CampaignID,
			}), nil).
			Once()
		suite.ChainMock.On("SourceHash").Return(testutil.ChainSourceHash).Once()
		suite.ChainMock.On("SourceURL").Return(testutil.ChainSourceURL).Once()
		suite.ChainMock.On("Name").Return(testutil.ChainName).Once()
		suite.ChainMock.On("ChainID").Return(testutil.ChainID, nil).Once()
		suite.ChainMock.
			On("CacheBinary", testutil.LaunchID).
			Return(expectedError).
			Once()

		_, _, publishError := network.Publish(context.Background(), suite.ChainMock)
		require.Error(t, publishError)
		require.Equal(t, expectedError, publishError)
		suite.AssertAllMocks(t)
	})
}
