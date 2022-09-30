package network

import (
	"context"
	"errors"
	"testing"

	sdkmath "cosmossdk.io/math"
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
		account := testutil.NewTestAccount(t, testutil.TestAccountName)
		tmp := t.TempDir()
		addr, err := account.Address(networktypes.SPN)
		require.NoError(t, err)
		gentx := testutil.NewGentx(
			addr,
			TestDenom,
			TestAmountString,
			"",
			testutil.PeerAddress,
		)
		gentxPath := gentx.SaveTo(t, tmp)
		suite, network := newSuite(account)

		suite.ChainMock.On("NodeID", context.Background()).Return(testutil.NodeID, nil).Once()
		suite.CosmosClientMock.
			On(
				"BroadcastTx",
				context.Background(),
				account,
				launchtypes.NewMsgSendRequest(
					addr,
					testutil.LaunchID,
					launchtypes.NewGenesisValidator(
						testutil.LaunchID,
						addr,
						gentx.JSON(t),
						[]byte{},
						sdk.NewCoin(TestDenom, sdkmath.NewInt(TestAmountInt)),
						launchtypes.Peer{
							Id: testutil.NodeID,
							Connection: &launchtypes.Peer_TcpAddress{
								TcpAddress: testutil.TCPAddress,
							},
						}),
				),
			).
			Return(testutil.NewResponse(&launchtypes.MsgSendRequestResponse{
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
		account := testutil.NewTestAccount(t, testutil.TestAccountName)
		tmp := t.TempDir()
		addr, err := account.Address(networktypes.SPN)
		require.NoError(t, err)
		gentx := testutil.NewGentx(
			addr,
			TestDenom,
			TestAmountString,
			"",
			testutil.PeerAddress,
		)
		gentxPath := gentx.SaveTo(t, tmp)
		suite, network := newSuite(account)

		suite.CosmosClientMock.
			On(
				"BroadcastTx",
				context.Background(),
				account,
				launchtypes.NewMsgSendRequest(
					addr,
					testutil.LaunchID,
					launchtypes.NewGenesisValidator(
						testutil.LaunchID,
						addr,
						gentx.JSON(t),
						[]byte{},
						sdk.NewCoin(TestDenom, sdkmath.NewInt(TestAmountInt)),
						launchtypes.Peer{
							Id: testutil.NodeID,
							Connection: &launchtypes.Peer_TcpAddress{
								TcpAddress: testutil.TCPAddress,
							},
						},
					),
				),
			).
			Return(testutil.NewResponse(&launchtypes.MsgSendRequestResponse{
				RequestID:    TestGenesisValidatorRequestID,
				AutoApproved: false,
			}), nil).
			Once()

		joinErr := network.Join(context.Background(), suite.ChainMock, testutil.LaunchID, gentxPath)
		require.NoError(t, joinErr)
		suite.AssertAllMocks(t)
	})

	t.Run("failed to send join request, failed to broadcast join tx", func(t *testing.T) {
		account := testutil.NewTestAccount(t, testutil.TestAccountName)
		tmp := t.TempDir()
		addr, err := account.Address(networktypes.SPN)
		require.NoError(t, err)
		gentx := testutil.NewGentx(
			addr,
			TestDenom,
			TestAmountString,
			"",
			testutil.PeerAddress,
		)
		gentxPath := gentx.SaveTo(t, tmp)
		suite, network := newSuite(account)
		expectedError := errors.New("failed to add validator")

		suite.ChainMock.On("NodeID", context.Background()).Return(testutil.NodeID, nil).Once()
		suite.CosmosClientMock.
			On(
				"BroadcastTx",
				context.Background(),
				account,
				launchtypes.NewMsgSendRequest(
					addr,
					testutil.LaunchID,
					launchtypes.NewGenesisValidator(
						testutil.LaunchID,
						addr,
						gentx.JSON(t),
						[]byte{},
						sdk.NewCoin(TestDenom, sdkmath.NewInt(TestAmountInt)),
						launchtypes.Peer{
							Id: testutil.NodeID,
							Connection: &launchtypes.Peer_TcpAddress{
								TcpAddress: testutil.TCPAddress,
							},
						},
					),
				),
			).
			Return(
				testutil.NewResponse(&launchtypes.MsgSendRequestResponse{}),
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
		account := testutil.NewTestAccount(t, testutil.TestAccountName)
		tmp := t.TempDir()
		addr, err := account.Address(networktypes.SPN)
		require.NoError(t, err)
		gentx := testutil.NewGentx(
			addr,
			TestDenom,
			TestAmountString,
			"",
			testutil.PeerAddress,
		)
		gentxPath := gentx.SaveTo(t, tmp)
		suite, network := newSuite(account)

		suite.ChainMock.On("NodeID", context.Background()).Return(testutil.NodeID, nil).Once()
		suite.CosmosClientMock.
			On(
				"BroadcastTx",
				context.Background(),
				account,
				launchtypes.NewMsgSendRequest(
					addr,
					testutil.LaunchID,
					launchtypes.NewGenesisValidator(
						testutil.LaunchID,
						addr,
						gentx.JSON(t),
						[]byte{},
						sdk.NewCoin(TestDenom, sdkmath.NewInt(TestAmountInt)),
						launchtypes.Peer{
							Id: testutil.NodeID,
							Connection: &launchtypes.Peer_TcpAddress{
								TcpAddress: testutil.TCPAddress,
							},
						},
					),
				),
			).
			Return(testutil.NewResponse(&launchtypes.MsgSendRequestResponse{
				RequestID:    TestGenesisValidatorRequestID,
				AutoApproved: false,
			}), nil).
			Once()
		suite.CosmosClientMock.
			On(
				"BroadcastTx",
				context.Background(),
				account,
				launchtypes.NewMsgSendRequest(
					addr,
					testutil.LaunchID,
					launchtypes.NewGenesisAccount(
						testutil.LaunchID,
						addr,
						sdk.NewCoins(sdk.NewCoin(TestDenom, sdkmath.NewInt(TestAmountInt))),
					),
				),
			).
			Return(testutil.NewResponse(&launchtypes.MsgSendRequestResponse{
				RequestID:    TestAccountRequestID,
				AutoApproved: false,
			}), nil).
			Once()

		joinErr := network.Join(
			context.Background(),
			suite.ChainMock,
			testutil.LaunchID,
			gentxPath,
			WithAccountRequest(sdk.NewCoins(sdk.NewCoin(TestDenom, sdkmath.NewInt(TestAmountInt)))),
			WithPublicAddress(testutil.TCPAddress),
		)
		require.NoError(t, joinErr)
		suite.AssertAllMocks(t)
	})

	t.Run("failed to send join request, failed to read node id", func(t *testing.T) {
		account := testutil.NewTestAccount(t, testutil.TestAccountName)
		tmp := t.TempDir()
		addr, err := account.Address(networktypes.SPN)
		require.NoError(t, err)
		gentx := testutil.NewGentx(
			addr,
			TestDenom,
			TestAmountString,
			"",
			testutil.PeerAddress,
		)
		gentxPath := gentx.SaveTo(t, tmp)
		suite, network := newSuite(account)
		expectedError := errors.New("failed to get node id")

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
