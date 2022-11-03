package network

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	campaigntypes "github.com/tendermint/spn/x/campaign/types"
	launchtypes "github.com/tendermint/spn/x/launch/types"
	profiletypes "github.com/tendermint/spn/x/profile/types"

	"github.com/ignite/cli/ignite/pkg/cosmoserror"
	"github.com/ignite/cli/ignite/services/network/networktypes"
	"github.com/ignite/cli/ignite/services/network/testutil"
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

func startInvalidJSONServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("invalid json"))
	}))
}

func TestPublish(t *testing.T) {
	t.Run("publish chain without campaign", func(t *testing.T) {
		var (
			account        = testutil.NewTestAccount(t, testutil.TestAccountName)
			suite, network = newSuite(account)
		)

		addr, err := account.Address(networktypes.SPN)
		require.NoError(t, err)

		suite.ProfileQueryMock.
			On(
				"CoordinatorByAddress",
				context.Background(),
				&profiletypes.QueryGetCoordinatorByAddressRequest{
					Address: addr,
				},
			).
			Return(&profiletypes.QueryGetCoordinatorByAddressResponse{
				CoordinatorByAddress: profiletypes.CoordinatorByAddress{
					Address:       addr,
					CoordinatorID: 1,
				},
			}, nil).
			Once()
		suite.CosmosClientMock.
			On(
				"BroadcastTx",
				context.Background(),
				account,
				&launchtypes.MsgCreateChain{
					Coordinator:    addr,
					GenesisChainID: testutil.ChainID,
					SourceURL:      testutil.ChainSourceURL,
					SourceHash:     testutil.ChainSourceHash,
					InitialGenesis: launchtypes.NewDefaultInitialGenesis(),
					HasCampaign:    false,
					CampaignID:     0,
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

		launchID, campaignID, publishError := network.Publish(context.Background(), suite.ChainMock)
		require.NoError(t, publishError)
		require.Equal(t, testutil.LaunchID, launchID)
		require.Equal(t, uint64(0), campaignID)
		suite.AssertAllMocks(t)
	})

	t.Run("publish chain with custom account balance", func(t *testing.T) {
		var (
			account        = testutil.NewTestAccount(t, testutil.TestAccountName)
			suite, network = newSuite(account)
		)

		accountBalance, err := sdk.ParseCoinsNormalized("1000foo,500bar")
		require.NoError(t, err)

		addr, err := account.Address(networktypes.SPN)
		require.NoError(t, err)

		suite.ProfileQueryMock.
			On(
				"CoordinatorByAddress",
				context.Background(),
				&profiletypes.QueryGetCoordinatorByAddressRequest{
					Address: addr,
				},
			).
			Return(&profiletypes.QueryGetCoordinatorByAddressResponse{
				CoordinatorByAddress: profiletypes.CoordinatorByAddress{
					Address:       addr,
					CoordinatorID: 1,
				},
			}, nil).
			Once()
		suite.CosmosClientMock.
			On(
				"BroadcastTx",
				context.Background(),
				account,
				&launchtypes.MsgCreateChain{
					Coordinator:    addr,
					GenesisChainID: testutil.ChainID,
					SourceURL:      testutil.ChainSourceURL,
					SourceHash:     testutil.ChainSourceHash,
					InitialGenesis: launchtypes.NewDefaultInitialGenesis(),
					HasCampaign:    false,
					CampaignID:     0,
					AccountBalance: accountBalance,
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

		launchID, campaignID, publishError := network.Publish(
			context.Background(),
			suite.ChainMock,
			WithAccountBalance(accountBalance),
		)
		require.NoError(t, publishError)
		require.Equal(t, testutil.LaunchID, launchID)
		require.Equal(t, uint64(0), campaignID)
		suite.AssertAllMocks(t)
	})

	t.Run("publish chain with pre created campaign", func(t *testing.T) {
		var (
			account        = testutil.NewTestAccount(t, testutil.TestAccountName)
			suite, network = newSuite(account)
		)

		addr, err := account.Address(networktypes.SPN)
		require.NoError(t, err)

		suite.ProfileQueryMock.
			On(
				"CoordinatorByAddress",
				context.Background(),
				&profiletypes.QueryGetCoordinatorByAddressRequest{
					Address: addr,
				},
			).
			Return(&profiletypes.QueryGetCoordinatorByAddressResponse{
				CoordinatorByAddress: profiletypes.CoordinatorByAddress{
					Address:       addr,
					CoordinatorID: 1,
				},
			}, nil).
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
				context.Background(),
				account,
				&launchtypes.MsgCreateChain{
					Coordinator:    addr,
					GenesisChainID: testutil.ChainID,
					SourceURL:      testutil.ChainSourceURL,
					SourceHash:     testutil.ChainSourceHash,
					InitialGenesis: launchtypes.NewDefaultInitialGenesis(),
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

	t.Run("publish chain with a pre created campaign with shares", func(t *testing.T) {
		var (
			account        = testutil.NewTestAccount(t, testutil.TestAccountName)
			suite, network = newSuite(account)
		)

		addr, err := account.Address(networktypes.SPN)
		require.NoError(t, err)

		suite.ProfileQueryMock.
			On(
				"CoordinatorByAddress",
				context.Background(),
				&profiletypes.QueryGetCoordinatorByAddressRequest{
					Address: addr,
				},
			).
			Return(&profiletypes.QueryGetCoordinatorByAddressResponse{
				CoordinatorByAddress: profiletypes.CoordinatorByAddress{
					Address:       addr,
					CoordinatorID: 1,
				},
			}, nil).
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
				context.Background(),
				account,
				campaigntypes.NewMsgMintVouchers(
					addr,
					testutil.CampaignID,
					campaigntypes.NewSharesFromCoins(sdk.NewCoins(sdk.NewInt64Coin("foo", 2000), sdk.NewInt64Coin("staking", 50000))),
				),
			).
			Return(testutil.NewResponse(&campaigntypes.MsgMintVouchersResponse{}), nil).
			Once()
		suite.CosmosClientMock.
			On(
				"BroadcastTx",
				context.Background(),
				account,
				&launchtypes.MsgCreateChain{
					Coordinator:    addr,
					GenesisChainID: testutil.ChainID,
					SourceURL:      testutil.ChainSourceURL,
					SourceHash:     testutil.ChainSourceHash,
					InitialGenesis: launchtypes.NewDefaultInitialGenesis(),
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

		launchID, campaignID, publishError := network.Publish(context.Background(), suite.ChainMock, WithCampaign(testutil.CampaignID),
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

	t.Run("publish chain with custom genesis url", func(t *testing.T) {
		var (
			account              = testutil.NewTestAccount(t, testutil.TestAccountName)
			customGenesisChainID = "test-custom-1"
			customGenesisHash    = "86167654c1af18c801837d443563fd98b3fe5e8d337e70faad181cdf2100da52"
			gts                  = startGenesisTestServer("mocks/data/genesis.json")
			suite, network       = newSuite(account)
		)
		defer gts.Close()

		addr, err := account.Address(networktypes.SPN)
		require.NoError(t, err)

		suite.ProfileQueryMock.
			On(
				"CoordinatorByAddress",
				context.Background(),
				&profiletypes.QueryGetCoordinatorByAddressRequest{
					Address: addr,
				},
			).
			Return(&profiletypes.QueryGetCoordinatorByAddressResponse{
				CoordinatorByAddress: profiletypes.CoordinatorByAddress{
					Address:       addr,
					CoordinatorID: 1,
				},
			}, nil).
			Once()
		suite.CosmosClientMock.
			On(
				"BroadcastTx",
				context.Background(),
				account,
				&launchtypes.MsgCreateChain{
					Coordinator:    addr,
					GenesisChainID: customGenesisChainID,
					SourceURL:      testutil.ChainSourceURL,
					SourceHash:     testutil.ChainSourceHash,
					InitialGenesis: launchtypes.NewGenesisURL(
						gts.URL,
						customGenesisHash,
					),
					HasCampaign: false,
					CampaignID:  0,
				},
			).
			Return(testutil.NewResponse(&launchtypes.MsgCreateChainResponse{
				LaunchID: testutil.LaunchID,
			}), nil).
			Once()
		suite.ChainMock.On("SourceHash").Return(testutil.ChainSourceHash).Once()
		suite.ChainMock.On("SourceURL").Return(testutil.ChainSourceURL).Once()
		suite.ChainMock.On("CacheBinary", testutil.LaunchID).Return(nil).Once()

		launchID, campaignID, publishError := network.Publish(
			context.Background(),
			suite.ChainMock,
			WithCustomGenesisURL(gts.URL),
		)
		require.NoError(t, publishError)
		require.Equal(t, testutil.LaunchID, launchID)
		require.Equal(t, uint64(0), campaignID)
		suite.AssertAllMocks(t)
	})

	t.Run("publish chain with custom chain id", func(t *testing.T) {
		var (
			account        = testutil.NewTestAccount(t, testutil.TestAccountName)
			suite, network = newSuite(account)
		)

		addr, err := account.Address(networktypes.SPN)
		require.NoError(t, err)

		suite.ProfileQueryMock.
			On(
				"CoordinatorByAddress",
				context.Background(),
				&profiletypes.QueryGetCoordinatorByAddressRequest{
					Address: addr,
				},
			).
			Return(&profiletypes.QueryGetCoordinatorByAddressResponse{
				CoordinatorByAddress: profiletypes.CoordinatorByAddress{
					Address:       addr,
					CoordinatorID: 1,
				},
			}, nil).
			Once()
		suite.CosmosClientMock.
			On(
				"BroadcastTx",
				context.Background(),
				account,
				&launchtypes.MsgCreateChain{
					Coordinator:    addr,
					GenesisChainID: testutil.ChainID,
					SourceURL:      testutil.ChainSourceURL,
					SourceHash:     testutil.ChainSourceHash,
					InitialGenesis: launchtypes.NewDefaultInitialGenesis(),
					HasCampaign:    false,
					CampaignID:     0,
				},
			).
			Return(testutil.NewResponse(&launchtypes.MsgCreateChainResponse{
				LaunchID: testutil.LaunchID,
			}), nil).
			Once()
		suite.ChainMock.On("SourceHash").Return(testutil.ChainSourceHash).Once()
		suite.ChainMock.On("SourceURL").Return(testutil.ChainSourceURL).Once()
		suite.ChainMock.On("CacheBinary", testutil.LaunchID).Return(nil).Once()

		launchID, campaignID, publishError := network.Publish(context.Background(), suite.ChainMock, WithChainID(testutil.ChainID))
		require.NoError(t, publishError)
		require.Equal(t, testutil.LaunchID, launchID)
		require.Equal(t, uint64(0), campaignID)
		suite.AssertAllMocks(t)
	})

	t.Run("publish chain with custom genesis config", func(t *testing.T) {
		var (
			account        = testutil.NewTestAccount(t, testutil.TestAccountName)
			suite, network = newSuite(account)
		)

		addr, err := account.Address(networktypes.SPN)
		require.NoError(t, err)

		suite.ProfileQueryMock.
			On(
				"CoordinatorByAddress",
				context.Background(),
				&profiletypes.QueryGetCoordinatorByAddressRequest{
					Address: addr,
				},
			).
			Return(&profiletypes.QueryGetCoordinatorByAddressResponse{
				CoordinatorByAddress: profiletypes.CoordinatorByAddress{
					Address:       addr,
					CoordinatorID: 1,
				},
			}, nil).
			Once()
		suite.CosmosClientMock.
			On(
				"BroadcastTx",
				context.Background(),
				account,
				&launchtypes.MsgCreateChain{
					Coordinator:    addr,
					GenesisChainID: testutil.ChainID,
					SourceURL:      testutil.ChainSourceURL,
					SourceHash:     testutil.ChainSourceHash,
					InitialGenesis: launchtypes.NewConfigGenesis(
						testutil.ChainConfigYML,
					),
					HasCampaign: false,
					CampaignID:  0,
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

		launchID, campaignID, publishError := network.Publish(
			context.Background(),
			suite.ChainMock,
			WithCustomGenesisConfig(testutil.ChainConfigYML),
		)
		require.NoError(t, publishError)
		require.Equal(t, testutil.LaunchID, launchID)
		require.Equal(t, uint64(0), campaignID)
		suite.AssertAllMocks(t)
	})

	t.Run("publish chain with custom chain id", func(t *testing.T) {
		var (
			account        = testutil.NewTestAccount(t, testutil.TestAccountName)
			suite, network = newSuite(account)
		)

		addr, err := account.Address(networktypes.SPN)
		require.NoError(t, err)

		suite.ProfileQueryMock.
			On(
				"CoordinatorByAddress",
				context.Background(),
				&profiletypes.QueryGetCoordinatorByAddressRequest{
					Address: addr,
				},
			).
			Return(&profiletypes.QueryGetCoordinatorByAddressResponse{
				CoordinatorByAddress: profiletypes.CoordinatorByAddress{
					Address:       addr,
					CoordinatorID: 1,
				},
			}, nil).
			Once()
		suite.CosmosClientMock.
			On(
				"BroadcastTx",
				context.Background(),
				account,
				&launchtypes.MsgCreateChain{
					Coordinator:    addr,
					GenesisChainID: testutil.ChainID,
					SourceURL:      testutil.ChainSourceURL,
					SourceHash:     testutil.ChainSourceHash,
					InitialGenesis: launchtypes.NewDefaultInitialGenesis(),
					HasCampaign:    false,
					CampaignID:     0,
				},
			).
			Return(testutil.NewResponse(&launchtypes.MsgCreateChainResponse{
				LaunchID: testutil.LaunchID,
			}), nil).
			Once()
		suite.ChainMock.On("SourceHash").Return(testutil.ChainSourceHash).Once()
		suite.ChainMock.On("SourceURL").Return(testutil.ChainSourceURL).Once()
		suite.ChainMock.On("CacheBinary", testutil.LaunchID).Return(nil).Once()

		launchID, campaignID, publishError := network.Publish(context.Background(), suite.ChainMock, WithChainID(testutil.ChainID))
		require.NoError(t, publishError)
		require.Equal(t, testutil.LaunchID, launchID)
		require.Equal(t, uint64(0), campaignID)
		suite.AssertAllMocks(t)
	})

	t.Run("publish chain with mainnet", func(t *testing.T) {
		var (
			account        = testutil.NewTestAccount(t, testutil.TestAccountName)
			gts            = startGenesisTestServer("mocks/data/genesis.json")
			suite, network = newSuite(account)
		)
		defer gts.Close()

		addr, err := account.Address(networktypes.SPN)
		require.NoError(t, err)

		suite.ProfileQueryMock.
			On(
				"CoordinatorByAddress",
				context.Background(),
				&profiletypes.QueryGetCoordinatorByAddressRequest{
					Address: addr,
				},
			).
			Return(&profiletypes.QueryGetCoordinatorByAddressResponse{
				CoordinatorByAddress: profiletypes.CoordinatorByAddress{
					Address:       addr,
					CoordinatorID: 1,
				},
			}, nil).
			Once()
		suite.CosmosClientMock.
			On(
				"BroadcastTx",
				context.Background(),
				account,
				&campaigntypes.MsgCreateCampaign{
					Coordinator:  addr,
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
				context.Background(),
				account,
				&campaigntypes.MsgInitializeMainnet{
					Coordinator:    addr,
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

		addr, err := account.Address(networktypes.SPN)
		require.NoError(t, err)

		suite.ProfileQueryMock.
			On(
				"CoordinatorByAddress",
				context.Background(),
				&profiletypes.QueryGetCoordinatorByAddressRequest{
					Address: addr,
				},
			).
			Return(&profiletypes.QueryGetCoordinatorByAddressResponse{
				CoordinatorByAddress: profiletypes.CoordinatorByAddress{
					Address:       addr,
					CoordinatorID: 1,
				},
			}, nil).
			Once()
		suite.CosmosClientMock.
			On(
				"BroadcastTx",
				context.Background(),
				account,
				&campaigntypes.MsgCreateCampaign{
					Coordinator:  addr,
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
				context.Background(),
				account,
				&campaigntypes.MsgInitializeMainnet{
					Coordinator:    addr,
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
		require.ErrorIs(t, publishError, expectedError)
		suite.AssertAllMocks(t)
	})

	t.Run("failed to publish chain with custom genesis, failed to parse custom genesis", func(t *testing.T) {
		var (
			account        = testutil.NewTestAccount(t, testutil.TestAccountName)
			suite, network = newSuite(account)
			gts            = startInvalidJSONServer()
			expectedError  = errors.New("JSON field not found")
		)
		defer gts.Close()

		_, _, publishError := network.Publish(
			context.Background(),
			suite.ChainMock,
			WithCustomGenesisURL(gts.URL),
		)
		require.Error(t, publishError)
		require.Equal(t, expectedError.Error(), publishError.Error())
		suite.AssertAllMocks(t)
	})

	t.Run("publish chain with coordinator creation", func(t *testing.T) {
		var (
			account        = testutil.NewTestAccount(t, testutil.TestAccountName)
			suite, network = newSuite(account)
		)

		addr, err := account.Address(networktypes.SPN)
		require.NoError(t, err)

		suite.ProfileQueryMock.
			On("CoordinatorByAddress", mock.Anything, &profiletypes.QueryGetCoordinatorByAddressRequest{
				Address: addr,
			}).
			Return(nil, cosmoserror.ErrNotFound).
			Once()
		suite.CosmosClientMock.
			On(
				"BroadcastTx",
				context.Background(),
				account,
				&launchtypes.MsgCreateChain{
					Coordinator:    addr,
					GenesisChainID: testutil.ChainID,
					SourceURL:      testutil.ChainSourceURL,
					SourceHash:     testutil.ChainSourceHash,
					InitialGenesis: launchtypes.NewDefaultInitialGenesis(),
					HasCampaign:    false,
					CampaignID:     0,
				},
			).
			Return(testutil.NewResponse(&launchtypes.MsgCreateChainResponse{
				LaunchID: testutil.LaunchID,
			}), nil).
			Once()
		suite.CosmosClientMock.
			On(
				"BroadcastTx",
				context.Background(),
				account,
				&profiletypes.MsgCreateCoordinator{
					Address: addr,
				},
			).
			Return(testutil.NewResponse(&profiletypes.MsgCreateCoordinatorResponse{}), nil).
			Once()
		suite.ChainMock.On("SourceHash").Return(testutil.ChainSourceHash).Once()
		suite.ChainMock.On("SourceURL").Return(testutil.ChainSourceURL).Once()
		suite.ChainMock.On("ChainID").Return(testutil.ChainID, nil).Once()
		suite.ChainMock.On("CacheBinary", testutil.LaunchID).Return(nil).Once()

		launchID, campaignID, publishError := network.Publish(context.Background(), suite.ChainMock)
		require.NoError(t, publishError)
		require.Equal(t, testutil.LaunchID, launchID)
		require.Equal(t, uint64(0), campaignID)
		suite.AssertAllMocks(t)
	})

	t.Run("failed to publish chain, failed to fetch coordinator profile", func(t *testing.T) {
		var (
			account        = testutil.NewTestAccount(t, testutil.TestAccountName)
			suite, network = newSuite(account)
			expectedError  = errors.New("failed to fetch coordinator")
		)

		addr, err := account.Address(networktypes.SPN)
		require.NoError(t, err)

		suite.ProfileQueryMock.
			On("CoordinatorByAddress", mock.Anything, &profiletypes.QueryGetCoordinatorByAddressRequest{
				Address: addr,
			}).
			Return(nil, expectedError).
			Once()
		suite.ChainMock.On("ChainID").Return(testutil.ChainID, nil).Once()

		_, _, publishError := network.Publish(context.Background(), suite.ChainMock)
		require.Error(t, publishError)
		require.ErrorIs(t, publishError, expectedError)
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
		require.ErrorIs(t, publishError, expectedError)
		suite.AssertAllMocks(t)
	})

	t.Run("failed to publish chain, failed to fetch existed campaign", func(t *testing.T) {
		var (
			account        = testutil.NewTestAccount(t, testutil.TestAccountName)
			suite, network = newSuite(account)
		)

		addr, err := account.Address(networktypes.SPN)
		require.NoError(t, err)

		suite.ProfileQueryMock.
			On(
				"CoordinatorByAddress",
				context.Background(),
				&profiletypes.QueryGetCoordinatorByAddressRequest{
					Address: addr,
				},
			).
			Return(&profiletypes.QueryGetCoordinatorByAddressResponse{
				CoordinatorByAddress: profiletypes.CoordinatorByAddress{
					Address:       addr,
					CoordinatorID: 1,
				},
			}, nil).
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
		require.ErrorIs(t, publishError, cosmoserror.ErrNotFound)
		suite.AssertAllMocks(t)
	})

	t.Run("failed to publish chain, failed to create chain", func(t *testing.T) {
		var (
			account        = testutil.NewTestAccount(t, testutil.TestAccountName)
			suite, network = newSuite(account)
			expectedError  = errors.New("failed to create chain")
		)

		addr, err := account.Address(networktypes.SPN)
		require.NoError(t, err)

		suite.ProfileQueryMock.
			On(
				"CoordinatorByAddress",
				context.Background(),
				&profiletypes.QueryGetCoordinatorByAddressRequest{
					Address: addr,
				},
			).
			Return(&profiletypes.QueryGetCoordinatorByAddressResponse{
				CoordinatorByAddress: profiletypes.CoordinatorByAddress{
					Address:       addr,
					CoordinatorID: 1,
				},
			}, nil).
			Once()
		suite.CosmosClientMock.
			On(
				"BroadcastTx",
				context.Background(),
				account,
				&launchtypes.MsgCreateChain{
					Coordinator:    addr,
					GenesisChainID: testutil.ChainID,
					SourceURL:      testutil.ChainSourceURL,
					SourceHash:     testutil.ChainSourceHash,
					InitialGenesis: launchtypes.NewDefaultInitialGenesis(),
					HasCampaign:    false,
					CampaignID:     0,
				},
			).
			Return(testutil.NewResponse(&launchtypes.MsgCreateChainResponse{
				LaunchID: testutil.LaunchID,
			}), expectedError).
			Once()
		suite.ChainMock.On("SourceHash").Return(testutil.ChainSourceHash).Once()
		suite.ChainMock.On("SourceURL").Return(testutil.ChainSourceURL).Once()
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

		addr, err := account.Address(networktypes.SPN)
		require.NoError(t, err)

		suite.ProfileQueryMock.
			On(
				"CoordinatorByAddress",
				context.Background(),
				&profiletypes.QueryGetCoordinatorByAddressRequest{
					Address: addr,
				},
			).
			Return(&profiletypes.QueryGetCoordinatorByAddressResponse{
				CoordinatorByAddress: profiletypes.CoordinatorByAddress{
					Address:       addr,
					CoordinatorID: 1,
				},
			}, nil).
			Once()
		suite.CosmosClientMock.
			On(
				"BroadcastTx",
				context.Background(),
				account,
				&launchtypes.MsgCreateChain{
					Coordinator:    addr,
					GenesisChainID: testutil.ChainID,
					SourceURL:      testutil.ChainSourceURL,
					SourceHash:     testutil.ChainSourceHash,
					InitialGenesis: launchtypes.NewDefaultInitialGenesis(),
					HasCampaign:    false,
					CampaignID:     0,
				},
			).
			Return(testutil.NewResponse(&launchtypes.MsgCreateChainResponse{
				LaunchID: testutil.LaunchID,
			}), nil).
			Once()
		suite.ChainMock.On("SourceHash").Return(testutil.ChainSourceHash).Once()
		suite.ChainMock.On("SourceURL").Return(testutil.ChainSourceURL).Once()
		suite.ChainMock.On("ChainID").Return(testutil.ChainID, nil).Once()
		suite.ChainMock.
			On("CacheBinary", testutil.LaunchID).
			Return(expectedError).
			Once()

		_, _, publishError := network.Publish(context.Background(), suite.ChainMock)
		require.Error(t, publishError)
		require.ErrorIs(t, publishError, expectedError)
		suite.AssertAllMocks(t)
	})
}
