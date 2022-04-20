package network

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	campaigntypes "github.com/tendermint/spn/x/campaign/types"
	launchtypes "github.com/tendermint/spn/x/launch/types"
	profiletypes "github.com/tendermint/spn/x/profile/types"

	"github.com/ignite-hq/cli/ignite/pkg/cosmoserror"
	"github.com/ignite-hq/cli/ignite/pkg/cosmosutil/genesis"
	"github.com/ignite-hq/cli/ignite/services/network/networktypes"
	"github.com/ignite-hq/cli/ignite/services/network/testutil"
)

func startGenesisTestServer(filepath string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		file, err := os.ReadFile(filepath)
		if err != nil {
			panic(err)
		}
		if _, err = w.Write(file); err != nil {
			panic(err)
		}
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

		launchID, campaignID, _, publishError := network.Publish(context.Background(), suite.ChainMock)
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

		launchID, campaignID, _, publishError := network.Publish(context.Background(), suite.ChainMock, WithCampaign(testutil.CampaignID))
		require.NoError(t, publishError)
		require.Equal(t, testutil.LaunchID, launchID)
		require.Equal(t, testutil.CampaignID, campaignID)
		suite.AssertAllMocks(t)
	})

	t.Run("publish chain with custom genesis", func(t *testing.T) {
		var (
			account              = testutil.NewTestAccount(t, testutil.TestAccountName)
			customGenesisChainID = "test-custom-1"
			customGenesisHash    = "d4ecd98aae9006940b7f350f57873bb6e56689cbeb73769f5d169bde8d001a92"
			gts                  = startGenesisTestServer("mocks/data/genesis.json")
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

		filepath := fmt.Sprintf("%s/genesis.json", t.TempDir())
		gen, err := genesis.FromURL(context.Background(), gts.URL, filepath)
		require.NoError(t, err)
		launchID, campaignID, _, publishError := network.Publish(context.Background(), suite.ChainMock, WithCustomGenesis(gen))
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

		launchID, campaignID, _, publishError := network.Publish(context.Background(), suite.ChainMock, WithChainID(testutil.ChainID))
		require.NoError(t, publishError)
		require.Equal(t, testutil.LaunchID, launchID)
		require.Equal(t, testutil.CampaignID, campaignID)
		suite.AssertAllMocks(t)
	})

	t.Run("publish chain with mainnet", func(t *testing.T) {
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
		suite.ChainMock.On("SourceHash").Return(testutil.ChainSourceHash).Times(2)
		suite.ChainMock.On("SourceURL").Return(testutil.ChainSourceURL).Times(2)
		suite.ChainMock.On("Name").Return(testutil.ChainName).Once()
		suite.ChainMock.On("ChainID").Return(testutil.ChainID, nil).Once()
		suite.ChainMock.On("CacheBinary", testutil.LaunchID).Return(nil).Once()

		launchID, campaignID, mainnetID, publishError := network.Publish(context.Background(), suite.ChainMock, Mainnet())
		require.NoError(t, publishError)
		require.Equal(t, testutil.LaunchID, launchID)
		require.Equal(t, testutil.CampaignID, campaignID)
		require.Equal(t, testutil.MainnetID, mainnetID)
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
		suite.ChainMock.On("SourceHash").Return(testutil.ChainSourceHash).Times(2)
		suite.ChainMock.On("SourceURL").Return(testutil.ChainSourceURL).Times(2)
		suite.ChainMock.On("Name").Return(testutil.ChainName).Once()
		suite.ChainMock.On("ChainID").Return(testutil.ChainID, nil).Once()
		suite.ChainMock.On("CacheBinary", testutil.LaunchID).Return(nil).Once()

		_, _, _, publishError := network.Publish(context.Background(), suite.ChainMock, Mainnet())
		require.Error(t, publishError)
		require.Equal(t, expectedError, publishError)
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

		launchID, campaignID, _, publishError := network.Publish(context.Background(), suite.ChainMock)
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

		_, _, _, publishError := network.Publish(context.Background(), suite.ChainMock)
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

		_, _, _, publishError := network.Publish(context.Background(), suite.ChainMock)
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

		_, _, _, publishError := network.Publish(context.Background(), suite.ChainMock, WithCampaign(testutil.CampaignID))
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

		_, _, _, publishError := network.Publish(context.Background(), suite.ChainMock)
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

		_, _, _, publishError := network.Publish(context.Background(), suite.ChainMock)
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

		_, _, _, publishError := network.Publish(context.Background(), suite.ChainMock)
		require.Error(t, publishError)
		require.Equal(t, expectedError, publishError)
		suite.AssertAllMocks(t)
	})
}
