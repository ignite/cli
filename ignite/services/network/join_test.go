package network

import (
	"context"
	"errors"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	launchtypes "github.com/tendermint/spn/x/launch/types"

	"github.com/ignite/cli/ignite/pkg/cosmoserror"
	"github.com/ignite/cli/ignite/services/network/networktypes"
	"github.com/ignite/cli/ignite/services/network/testutil"
)

const (
	TestDenom                     = "stake"
	TestAmountString              = "95000000"
	TestCustomNodeID              = "b91f4adbfb0c0b513040d914bfb717303c0eaa17"
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
				"",
			)
			gentxPath      = gentx.SaveTo(t, tmp)
			genesisPath    = testutil.NewGenesis(testutil.ChainID).SaveTo(t, tmp)
			suite, network = newSuite(account)
		)

		suite.ChainMock.On("NodeID", context.Background()).Return(testutil.NodeID, nil).Once()
		suite.ChainMock.On("GenesisPath").Return(genesisPath, nil).Once()
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

		joinErr := network.Join(
			context.Background(),
			suite.ChainMock,
			testutil.LaunchID,
			gentxPath,
			WithPublicAddress(testutil.TCPAddress),
		)
		require.NoError(t, joinErr)
		suite.AssertAllMocks(t)
	})

	t.Run("successfully send join request with custom node id", func(t *testing.T) {
		var (
			account = testutil.NewTestAccount(t, testutil.TestAccountName)
			tmp     = t.TempDir()
			gentx   = testutil.NewGentx(
				account.Address(networktypes.SPN),
				TestDenom,
				TestAmountString,
				"",
				"",
			)
			gentxPath      = gentx.SaveTo(t, tmp)
			genesisPath    = testutil.NewGenesis(testutil.ChainID).SaveTo(t, tmp)
			suite, network = newSuite(account)
		)

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
						Id: TestCustomNodeID,
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

		joinErr := network.Join(
			context.Background(),
			suite.ChainMock,
			testutil.LaunchID,
			gentxPath,
			WithPublicAddress(testutil.TCPAddress),
			WithNodeID(TestCustomNodeID),
		)
		require.NoError(t, joinErr)
		suite.AssertAllMocks(t)
	})

	t.Run("successfully send join request with hidden public address", func(t *testing.T) {
		var (
			account = testutil.NewTestAccount(t, testutil.TestAccountName)
			tmp     = t.TempDir()
			gentx   = testutil.NewGentx(
				account.Address(networktypes.SPN),
				TestDenom,
				TestAmountString,
				"",
				"random memo",
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
						Connection: &launchtypes.Peer_None{
							None: &launchtypes.Peer_EmptyConnection{},
						},
					},
				},
			).
			Return(testutil.NewResponse(&launchtypes.MsgRequestAddValidatorResponse{
				RequestID:    TestGenesisValidatorRequestID,
				AutoApproved: false,
			}), nil).
			Once()

		joinErr := network.Join(
			context.Background(),
			suite.ChainMock,
			testutil.LaunchID,
			gentxPath,
			WithPublicAddress(""),
		)
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
				"",
			)
			gentxPath      = gentx.SaveTo(t, tmp)
			genesisPath    = testutil.NewGenesis(testutil.ChainID).SaveTo(t, tmp)
			suite, network = newSuite(account)
		)

		suite.ChainMock.On("GenesisPath").Return(genesisPath, nil).Once()
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

		joinErr := network.Join(context.Background(), suite.ChainMock, testutil.LaunchID, gentxPath)
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
				"",
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

		joinErr := network.Join(
			context.Background(),
			suite.ChainMock,
			testutil.LaunchID,
			gentxPath,
			WithPublicAddress(testutil.TCPAddress),
		)
		require.Errorf(t, joinErr, "validator %s already exist", account.Address(networktypes.SPN))
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
				"",
			)
			gentxPath      = gentx.SaveTo(t, tmp)
			genesisPath    = testutil.NewGenesis(testutil.ChainID).SaveTo(t, tmp)
			suite, network = newSuite(account)
			expectedError  = errors.New("failed to perform request")
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
			Return(nil, expectedError).
			Once()

		joinErr := network.Join(
			context.Background(),
			suite.ChainMock,
			testutil.LaunchID,
			gentxPath,
			WithPublicAddress(testutil.TCPAddress),
		)
		require.Error(t, joinErr)
		require.Equal(t, expectedError, joinErr)
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
				"",
			)
			gentxPath      = gentx.SaveTo(t, tmp)
			genesisPath    = testutil.NewGenesis(testutil.ChainID).SaveTo(t, tmp)
			suite, network = newSuite(account)
			expectedError  = errors.New("failed to add validator")
		)

		suite.ChainMock.On("NodeID", context.Background()).Return(testutil.NodeID, nil).Once()
		suite.ChainMock.On("GenesisPath").Return(genesisPath, nil).Once()
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
				expectedError,
			).
			Once()

		joinErr := network.Join(
			context.Background(),
			suite.ChainMock,
			testutil.LaunchID,
			gentxPath,
			WithPublicAddress(testutil.TCPAddress),
		)
		require.Error(t, joinErr)
		require.Equal(t, expectedError, joinErr)
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
				"",
			)
			gentxPath      = gentx.SaveTo(t, tmp)
			genesisPath    = testutil.NewGenesis(testutil.ChainID).SaveTo(t, tmp)
			suite, network = newSuite(account)
		)

		suite.ChainMock.On("NodeID", context.Background()).Return(testutil.NodeID, nil).Once()
		suite.ChainMock.On("GenesisPath").Return(genesisPath, nil).Once()
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
			gentxPath,
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
				"",
			)
			gentxPath      = gentx.SaveTo(t, tmp)
			genesis        = testutil.NewGenesis(testutil.ChainID).AddAccount(account.Address(networktypes.SPN))
			genesisPath    = genesis.SaveTo(t, tmp)
			suite, network = newSuite(account)
		)

		suite.ChainMock.On("NodeID", context.Background()).Return(testutil.NodeID, nil).Once()
		suite.ChainMock.On("GenesisPath").Return(genesisPath, nil).Once()

		joinErr := network.Join(
			context.Background(),
			suite.ChainMock,
			testutil.LaunchID,
			gentxPath,
			WithAccountRequest(sdk.NewCoins(sdk.NewCoin(TestDenom, sdk.NewInt(TestAmountInt)))),
			WithPublicAddress(testutil.TCPAddress),
		)
		require.Errorf(t, joinErr, "account %s already exist", account.Address(networktypes.SPN))
		suite.AssertAllMocks(t)
	})

	t.Run("failed to send join request, failed to read node id", func(t *testing.T) {
		var (
			account = testutil.NewTestAccount(t, testutil.TestAccountName)
			tmp     = t.TempDir()
			gentx   = testutil.NewGentx(
				account.Address(networktypes.SPN),
				TestDenom,
				TestAmountString,
				"",
				"",
			)
			gentxPath      = gentx.SaveTo(t, tmp)
			suite, network = newSuite(account)
			expectedError  = errors.New("failed to get node id")
		)
		suite.ChainMock.
			On("NodeID", mock.Anything).
			Return("", expectedError).
			Once()

		joinErr := network.Join(
			context.Background(),
			suite.ChainMock,
			testutil.LaunchID,
			gentxPath,
			WithPublicAddress(testutil.TCPAddress),
		)
		require.Error(t, joinErr)
		require.Equal(t, expectedError, joinErr)
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
				"",
			)
			gentxPath      = gentx.SaveTo(t, tmp)
			suite, network = newSuite(account)
			expectedError  = errors.New("failed to get genesis path")
		)

		suite.ChainMock.On("NodeID", context.Background()).Return(testutil.NodeID, nil).Once()
		suite.ChainMock.
			On("GenesisPath").
			Return("", expectedError).
			Once()

		joinErr := network.Join(
			context.Background(),
			suite.ChainMock,
			testutil.LaunchID,
			gentxPath,
			WithPublicAddress(testutil.TCPAddress),
		)
		require.Error(t, joinErr)
		require.Equal(t, expectedError, joinErr)
		suite.AssertAllMocks(t)
	})

	t.Run("failed to send join request, failed to read custom genesis", func(t *testing.T) {
		var (
			account        = testutil.NewTestAccount(t, testutil.TestAccountName)
			gentxPath      = "invalid/path"
			suite, network = newSuite(account)
			expectedError  = errors.New("chain home folder is not initialized yet: invalid/path")
		)

		joinErr := network.Join(context.Background(), suite.ChainMock, testutil.LaunchID, gentxPath)
		require.Error(t, joinErr)
		require.Equal(t, expectedError, joinErr)
		suite.AssertAllMocks(t)
	})

	t.Run("failed to send join request, unsupported peer address format", func(t *testing.T) {
		var (
			account        = testutil.NewTestAccount(t, testutil.TestAccountName)
			suite, network = newSuite(account)
			tmp            = t.TempDir()
			gentx          = testutil.NewGentx(
				account.Address(networktypes.SPN),
				TestDenom,
				TestAmountString,
				"",
				"",
			)
			gentxPath     = gentx.SaveTo(t, tmp)
			expectedError = errors.New("unsupported public address format: invalid")
		)

		suite.ChainMock.On("NodeID", context.Background()).Return(testutil.NodeID, nil).Once()

		joinErr := network.Join(
			context.Background(),
			suite.ChainMock,
			testutil.LaunchID,
			gentxPath,
			WithPublicAddress("invalid"),
		)
		require.Error(t, joinErr)
		require.Equal(t, expectedError, joinErr)
		suite.AssertAllMocks(t)
	})

}
