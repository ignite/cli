package network

import (
	"context"
	"errors"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ignite-hq/cli/ignite/pkg/cosmoserror"
	"github.com/ignite-hq/cli/ignite/services/network/networktypes"
	"github.com/ignite-hq/cli/ignite/services/network/testutil"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	launchtypes "github.com/tendermint/spn/x/launch/types"
)

const (
	TestDenom                     = "stake"
	TestAmountString              = "95000000"
	TestAmountInt                 = int64(95000000)
	TestAccountRequestID          = uint64(1)
	TestGenesisValidatorRequestID = uint64(2)
)

func TestJoin(t *testing.T) {
	t.Run("successfully send join request", func(t *testing.T) {
		var (
			account = testutil.NewTestAccount(t, testutil.TestAccountName)
			tmp     = t.TempDir()
			gentx   = testutil.NewGentx(
				account.Address(networktypes.SPN),
				TestDenom,
				TestAmountString,
				"",
				testutil.PeerAddress,
			)
			gentxPath      = gentx.SaveTo(t, tmp)
			genesisPath    = testutil.NewGenesis(testutil.ChainID).SaveTo(t, tmp)
			suite, network = newSuite(account)
		)

		suite.ChainMock.On("NodeID", context.Background()).Return(testutil.NodeID, nil).Once()
		suite.ChainMock.On("DefaultGentxPath").Return(gentxPath, nil).Once()
		suite.ChainMock.On("GenesisPath").Return(genesisPath, nil).Once()
		suite.LaunchQueryMock.
			On(
				"GenesisValidator",
				context.Background(),
				&launchtypes.QueryGetGenesisValidatorRequest{
					Address:  account.Address(networktypes.SPN),
					LaunchID: testutil.LaunchID,
				},
			).
			Return(nil, cosmoserror.ErrNotFound).
			Once()
		suite.CosmosClientMock.
			On(
				"BroadcastTx",
				account.Name,
				&launchtypes.MsgRequestAddValidator{
					Creator:        account.Address(networktypes.SPN),
					LaunchID:       testutil.LaunchID,
					ValAddress:     account.Address(networktypes.SPN),
					GenTx:          gentx.JSON(t),
					ConsPubKey:     []byte{},
					SelfDelegation: sdk.NewCoin(TestDenom, sdk.NewInt(TestAmountInt)),
					Peer: launchtypes.Peer{
						Id: testutil.NodeID,
						Connection: &launchtypes.Peer_TcpAddress{
							TcpAddress: testutil.TCPAddress,
						},
					},
				},
			).
			Return(testutil.NewResponse(&launchtypes.MsgRequestAddValidatorResponse{
				RequestID:    TestGenesisValidatorRequestID,
				AutoApproved: false,
			}), nil).
			Once()

		joinErr := network.Join(context.Background(), suite.ChainMock, testutil.LaunchID, WithPublicAddress(testutil.TCPAddress))
		require.NoError(t, joinErr)
		suite.AssertAllMocks(t)
	})

	t.Run("successfully send join request with custom gentx", func(t *testing.T) {
		var (
			account = testutil.NewTestAccount(t, testutil.TestAccountName)
			tmp     = t.TempDir()
			gentx   = testutil.NewGentx(
				account.Address(networktypes.SPN),
				TestDenom,
				TestAmountString,
				"",
				testutil.PeerAddress,
			)
			gentxPath      = gentx.SaveTo(t, tmp)
			genesisPath    = testutil.NewGenesis(testutil.ChainID).SaveTo(t, tmp)
			suite, network = newSuite(account)
		)
		suite.ChainMock.On("GenesisPath").Return(genesisPath, nil).Once()
		suite.LaunchQueryMock.
			On(
				"GenesisValidator",
				context.Background(),
				&launchtypes.QueryGetGenesisValidatorRequest{
					Address:  account.Address(networktypes.SPN),
					LaunchID: testutil.LaunchID,
				},
			).
			Return(nil, cosmoserror.ErrNotFound).
			Once()
		suite.CosmosClientMock.
			On(
				"BroadcastTx",
				account.Name,
				&launchtypes.MsgRequestAddValidator{
					Creator:        account.Address(networktypes.SPN),
					LaunchID:       testutil.LaunchID,
					ValAddress:     account.Address(networktypes.SPN),
					GenTx:          gentx.JSON(t),
					ConsPubKey:     []byte{},
					SelfDelegation: sdk.NewCoin(TestDenom, sdk.NewInt(TestAmountInt)),
					Peer: launchtypes.Peer{
						Id: testutil.NodeID,
						Connection: &launchtypes.Peer_TcpAddress{
							TcpAddress: testutil.TCPAddress,
						},
					},
				},
			).
			Return(testutil.NewResponse(&launchtypes.MsgRequestAddValidatorResponse{
				RequestID:    TestGenesisValidatorRequestID,
				AutoApproved: false,
			}), nil).
			Once()

		joinErr := network.Join(context.Background(), suite.ChainMock, testutil.LaunchID, WithCustomGentxPath(gentxPath))
		require.NoError(t, joinErr)
		suite.AssertAllMocks(t)
	})

	t.Run("failed to send join request, validator already exists", func(t *testing.T) {
		var (
			account = testutil.NewTestAccount(t, testutil.TestAccountName)
			tmp     = t.TempDir()
			gentx   = testutil.NewGentx(
				account.Address(networktypes.SPN),
				TestDenom,
				TestAmountString,
				"",
				testutil.PeerAddress,
			)
			gentxPath      = gentx.SaveTo(t, tmp)
			genesisPath    = testutil.NewGenesis(testutil.ChainID).SaveTo(t, tmp)
			suite, network = newSuite(account)
		)

		suite.ChainMock.On("NodeID", context.Background()).Return(testutil.NodeID, nil).Once()
		suite.ChainMock.On("GenesisPath").Return(genesisPath, nil).Once()
		suite.ChainMock.On("DefaultGentxPath").Return(gentxPath, nil).Once()
		suite.LaunchQueryMock.
			On(
				"GenesisValidator",
				context.Background(),
				&launchtypes.QueryGetGenesisValidatorRequest{
					Address:  account.Address(networktypes.SPN),
					LaunchID: testutil.LaunchID,
				},
			).
			Return(nil, nil).
			Once()

		joinErr := network.Join(context.Background(), suite.ChainMock, testutil.LaunchID, WithPublicAddress(testutil.TCPAddress))
		require.Error(t, joinErr)
		suite.AssertAllMocks(t)
	})

	t.Run("failed to send join request, failed to check validator existence", func(t *testing.T) {
		var (
			account = testutil.NewTestAccount(t, testutil.TestAccountName)
			tmp     = t.TempDir()
			gentx   = testutil.NewGentx(
				account.Address(networktypes.SPN),
				TestDenom,
				TestAmountString,
				"",
				testutil.PeerAddress,
			)
			gentxPath      = gentx.SaveTo(t, tmp)
			genesisPath    = testutil.NewGenesis(testutil.ChainID).SaveTo(t, tmp)
			suite, network = newSuite(account)
		)

		suite.ChainMock.On("NodeID", context.Background()).Return(testutil.NodeID, nil).Once()
		suite.ChainMock.On("GenesisPath").Return(genesisPath, nil).Once()
		suite.ChainMock.On("DefaultGentxPath").Return(gentxPath, nil).Once()
		suite.LaunchQueryMock.
			On(
				"GenesisValidator",
				context.Background(),
				&launchtypes.QueryGetGenesisValidatorRequest{
					Address:  account.Address(networktypes.SPN),
					LaunchID: testutil.LaunchID,
				},
			).
			Return(nil, errors.New("failed to perform request")).
			Once()

		joinErr := network.Join(context.Background(), suite.ChainMock, testutil.LaunchID, WithPublicAddress(testutil.TCPAddress))
		require.Error(t, joinErr)
		suite.AssertAllMocks(t)
	})

	t.Run("failed to send join request, failed to broadcast join tx", func(t *testing.T) {
		var (
			account = testutil.NewTestAccount(t, testutil.TestAccountName)
			tmp     = t.TempDir()
			gentx   = testutil.NewGentx(
				account.Address(networktypes.SPN),
				TestDenom,
				TestAmountString,
				"",
				testutil.PeerAddress,
			)
			gentxPath      = gentx.SaveTo(t, tmp)
			genesisPath    = testutil.NewGenesis(testutil.ChainID).SaveTo(t, tmp)
			suite, network = newSuite(account)
		)

		suite.ChainMock.On("NodeID", context.Background()).Return(testutil.NodeID, nil).Once()
		suite.ChainMock.On("GenesisPath").Return(genesisPath, nil).Once()
		suite.ChainMock.On("DefaultGentxPath").Return(gentxPath, nil).Once()
		suite.LaunchQueryMock.
			On(
				"GenesisValidator",
				context.Background(),
				&launchtypes.QueryGetGenesisValidatorRequest{
					Address:  account.Address(networktypes.SPN),
					LaunchID: testutil.LaunchID,
				},
			).
			Return(nil, cosmoserror.ErrNotFound).
			Once()
		suite.CosmosClientMock.
			On(
				"BroadcastTx",
				account.Name,
				&launchtypes.MsgRequestAddValidator{
					Creator:        account.Address(networktypes.SPN),
					LaunchID:       testutil.LaunchID,
					ValAddress:     account.Address(networktypes.SPN),
					GenTx:          gentx.JSON(t),
					ConsPubKey:     []byte{},
					SelfDelegation: sdk.NewCoin(TestDenom, sdk.NewInt(TestAmountInt)),
					Peer: launchtypes.Peer{
						Id: testutil.NodeID,
						Connection: &launchtypes.Peer_TcpAddress{
							TcpAddress: testutil.TCPAddress,
						},
					},
				},
			).
			Return(
				testutil.NewResponse(&launchtypes.MsgRequestAddValidatorResponse{}),
				errors.New("failed to add validator"),
			).
			Once()

		joinErr := network.Join(context.Background(), suite.ChainMock, testutil.LaunchID, WithPublicAddress(testutil.TCPAddress))
		require.Error(t, joinErr)
		suite.AssertAllMocks(t)
	})

	t.Run("successfully send join request with account request", func(t *testing.T) {
		var (
			account = testutil.NewTestAccount(t, testutil.TestAccountName)
			tmp     = t.TempDir()
			gentx   = testutil.NewGentx(
				account.Address(networktypes.SPN),
				TestDenom,
				TestAmountString,
				"",
				testutil.PeerAddress,
			)
			gentxPath      = gentx.SaveTo(t, tmp)
			genesisPath    = testutil.NewGenesis(testutil.ChainID).SaveTo(t, tmp)
			suite, network = newSuite(account)
		)

		suite.ChainMock.On("NodeID", context.Background()).Return(testutil.NodeID, nil).Once()
		suite.ChainMock.On("GenesisPath").Return(genesisPath, nil).Once()
		suite.ChainMock.On("DefaultGentxPath").Return(gentxPath, nil).Once()
		suite.LaunchQueryMock.
			On(
				"GenesisValidator",
				context.Background(),
				&launchtypes.QueryGetGenesisValidatorRequest{
					Address:  account.Address(networktypes.SPN),
					LaunchID: testutil.LaunchID,
				},
			).
			Return(nil, cosmoserror.ErrNotFound).
			Once()
		suite.LaunchQueryMock.
			On(
				"VestingAccount",
				context.Background(),
				&launchtypes.QueryGetVestingAccountRequest{
					Address:  account.Address(networktypes.SPN),
					LaunchID: testutil.LaunchID,
				}).
			Return(nil, cosmoserror.ErrNotFound).
			Once()
		suite.CosmosClientMock.
			On(
				"BroadcastTx",
				account.Name,
				&launchtypes.MsgRequestAddValidator{
					Creator:        account.Address(networktypes.SPN),
					LaunchID:       testutil.LaunchID,
					ValAddress:     account.Address(networktypes.SPN),
					GenTx:          gentx.JSON(t),
					ConsPubKey:     []byte{},
					SelfDelegation: sdk.NewCoin(TestDenom, sdk.NewInt(TestAmountInt)),
					Peer: launchtypes.Peer{
						Id: testutil.NodeID,
						Connection: &launchtypes.Peer_TcpAddress{
							TcpAddress: testutil.TCPAddress,
						},
					},
				},
			).
			Return(testutil.NewResponse(&launchtypes.MsgRequestAddValidatorResponse{
				RequestID:    TestGenesisValidatorRequestID,
				AutoApproved: false,
			}), nil).
			Once()
		suite.CosmosClientMock.
			On(
				"BroadcastTx",
				account.Name,
				&launchtypes.MsgRequestAddAccount{
					Creator:  account.Address(networktypes.SPN),
					LaunchID: testutil.LaunchID,
					Address:  account.Address(networktypes.SPN),
					Coins:    sdk.NewCoins(sdk.NewCoin(TestDenom, sdk.NewInt(TestAmountInt))),
				},
			).
			Return(testutil.NewResponse(&launchtypes.MsgRequestAddAccountResponse{
				RequestID:    TestAccountRequestID,
				AutoApproved: false,
			}), nil).
			Once()

		joinErr := network.Join(
			context.Background(),
			suite.ChainMock,
			testutil.LaunchID,
			WithAccountRequest(sdk.NewCoins(sdk.NewCoin(TestDenom, sdk.NewInt(TestAmountInt)))),
			WithPublicAddress(testutil.TCPAddress),
		)
		require.NoError(t, joinErr)
		suite.AssertAllMocks(t)
	})

	t.Run("failed to send join request with account request, account exists in genesis", func(t *testing.T) {
		var (
			account = testutil.NewTestAccount(t, testutil.TestAccountName)
			tmp     = t.TempDir()
			gentx   = testutil.NewGentx(
				account.Address(networktypes.SPN),
				TestDenom,
				TestAmountString,
				"",
				testutil.PeerAddress,
			)
			gentxPath      = gentx.SaveTo(t, tmp)
			genesis        = testutil.NewGenesis(testutil.ChainID).AddAccount(account.Address(networktypes.SPN))
			genesisPath    = genesis.SaveTo(t, tmp)
			suite, network = newSuite(account)
		)

		suite.ChainMock.On("NodeID", context.Background()).Return(testutil.NodeID, nil).Once()
		suite.ChainMock.On("GenesisPath").Return(genesisPath, nil).Once()
		suite.ChainMock.On("DefaultGentxPath").Return(gentxPath, nil).Once()

		joinErr := network.Join(
			context.Background(),
			suite.ChainMock,
			testutil.LaunchID,
			WithAccountRequest(sdk.NewCoins(sdk.NewCoin(TestDenom, sdk.NewInt(TestAmountInt)))),
			WithPublicAddress(testutil.TCPAddress),
		)
		require.Error(t, joinErr)
		suite.AssertAllMocks(t)
	})

	t.Run("failed to send join request with account request, failed to broadcast account tx", func(t *testing.T) {
		var (
			account = testutil.NewTestAccount(t, testutil.TestAccountName)
			tmp     = t.TempDir()
			gentx   = testutil.NewGentx(
				account.Address(networktypes.SPN),
				TestDenom,
				TestAmountString,
				"",
				testutil.PeerAddress,
			)
			gentxPath      = gentx.SaveTo(t, tmp)
			genesisPath    = testutil.NewGenesis(testutil.ChainID).SaveTo(t, tmp)
			suite, network = newSuite(account)
		)

		suite.ChainMock.On("NodeID", context.Background()).Return(testutil.NodeID, nil).Once()
		suite.ChainMock.On("GenesisPath").Return(genesisPath, nil).Once()
		suite.ChainMock.On("DefaultGentxPath").Return(gentxPath, nil).Once()

		suite.LaunchQueryMock.
			On(
				"VestingAccount",
				context.Background(),
				&launchtypes.QueryGetVestingAccountRequest{
					Address:  account.Address(networktypes.SPN),
					LaunchID: testutil.LaunchID,
				}).
			Return(nil, cosmoserror.ErrNotFound).
			Once()
		suite.CosmosClientMock.
			On(
				"BroadcastTx",
				account.Name,
				&launchtypes.MsgRequestAddAccount{
					Creator:  account.Address(networktypes.SPN),
					LaunchID: testutil.LaunchID,
					Address:  account.Address(networktypes.SPN),
					Coins:    sdk.NewCoins(sdk.NewCoin(TestDenom, sdk.NewInt(TestAmountInt))),
				},
			).
			Return(
				testutil.NewResponse(&launchtypes.MsgRequestAddAccountResponse{}),
				errors.New("failed to create account"),
			).
			Once()

		joinErr := network.Join(
			context.Background(),
			suite.ChainMock,
			testutil.LaunchID,
			WithAccountRequest(sdk.NewCoins(sdk.NewCoin(TestDenom, sdk.NewInt(TestAmountInt)))),
			WithPublicAddress(testutil.TCPAddress),
		)
		require.Error(t, joinErr)
		suite.AssertAllMocks(t)
	})

	t.Run("failed to send join request, failed to read node id", func(t *testing.T) {
		var (
			account        = testutil.NewTestAccount(t, testutil.TestAccountName)
			suite, network = newSuite(account)
		)
		suite.ChainMock.
			On("NodeID", mock.Anything).
			Return("", errors.New("failed to get node id")).
			Once()

		joinErr := network.Join(context.Background(), suite.ChainMock, testutil.LaunchID, WithPublicAddress(testutil.TCPAddress))
		require.Error(t, joinErr)
		suite.AssertAllMocks(t)
	})

	t.Run("failed to send join request, failed to read default gentx", func(t *testing.T) {
		var (
			account        = testutil.NewTestAccount(t, testutil.TestAccountName)
			suite, network = newSuite(account)
		)

		suite.ChainMock.On("NodeID", context.Background()).Return(testutil.NodeID, nil).Once()
		suite.ChainMock.
			On("DefaultGentxPath").
			Return("", errors.New("failed to get default gentx path")).
			Once()

		joinErr := network.Join(context.Background(), suite.ChainMock, testutil.LaunchID, WithPublicAddress(testutil.TCPAddress))
		require.Error(t, joinErr)
		suite.AssertAllMocks(t)
	})

	t.Run("failed to send join request, failed to read genesis", func(t *testing.T) {
		var (
			account = testutil.NewTestAccount(t, testutil.TestAccountName)
			tmp     = t.TempDir()
			gentx   = testutil.NewGentx(
				account.Address(networktypes.SPN),
				TestDenom,
				TestAmountString,
				"",
				testutil.PeerAddress,
			)
			gentxPath      = gentx.SaveTo(t, tmp)
			suite, network = newSuite(account)
		)

		suite.ChainMock.On("NodeID", context.Background()).Return(testutil.NodeID, nil).Once()
		suite.ChainMock.On("DefaultGentxPath").Return(gentxPath, nil).Once()
		suite.ChainMock.
			On("GenesisPath").
			Return("", errors.New("failed to get genesis path")).
			Once()

		joinErr := network.Join(context.Background(), suite.ChainMock, testutil.LaunchID, WithPublicAddress(testutil.TCPAddress))
		require.Error(t, joinErr)
		suite.AssertAllMocks(t)
	})

	t.Run("failed to send join request, failed to read custom genesis", func(t *testing.T) {
		var (
			account        = testutil.NewTestAccount(t, testutil.TestAccountName)
			gentxPath      = "invalid/path"
			suite, network = newSuite(account)
		)

		joinErr := network.Join(context.Background(), suite.ChainMock, testutil.LaunchID, WithCustomGentxPath(gentxPath))
		require.Error(t, joinErr)
		suite.ChainMock.AssertNumberOfCalls(t, "NodeID", 0)
		suite.ChainMock.AssertNumberOfCalls(t, "GenesisPath", 0)
		suite.ChainMock.AssertNumberOfCalls(t, "DefaultGentxPath", 0)
	})
}
