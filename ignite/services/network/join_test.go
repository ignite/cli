package network

import (
	"context"
	"errors"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	launchtypes "github.com/tendermint/spn/x/launch/types"

	"github.com/ignite/cli/ignite/services/network/networktypes"
	"github.com/ignite/cli/ignite/services/network/testutil"
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
			suite, network = newSuite(account)
		)

		suite.ChainMock.On("NodeID", context.Background()).Return(testutil.NodeID, nil).Once()
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
			suite, network = newSuite(account)
		)

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
			suite, network = newSuite(account)
			expectedError  = errors.New("failed to add validator")
		)

		suite.ChainMock.On("NodeID", context.Background()).Return(testutil.NodeID, nil).Once()
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
				testutil.PeerAddress,
			)
			gentxPath      = gentx.SaveTo(t, tmp)
			suite, network = newSuite(account)
		)

		suite.ChainMock.On("NodeID", context.Background()).Return(testutil.NodeID, nil).Once()
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

	t.Run("failed to send join request, failed to read node id", func(t *testing.T) {
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

	t.Run("failed to send join request, failed to read gentx", func(t *testing.T) {
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
}
