package cosmosclient_test

import (
	"encoding/hex"
	"testing"

	"cosmossdk.io/math"
	abci "github.com/cometbft/cometbft/abci/types"
	ctypes "github.com/cometbft/cometbft/rpc/core/types"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/v29/ignite/pkg/cosmosaccount"
	"github.com/ignite/cli/v29/ignite/pkg/cosmosclient"
	"github.com/ignite/cli/v29/ignite/pkg/errors"
)

func TestTxServiceBroadcast(t *testing.T) {
	var (
		accountName = "bob"
		passphrase  = "passphrase"
		txHash      = []byte{1, 2, 3}
		txHashStr   = hex.EncodeToString(txHash)
	)
	r, err := cosmosaccount.NewInMemory()
	require.NoError(t, err)
	a, _, err := r.Create(accountName)
	require.NoError(t, err)
	// Export created account to we can import it in the Client below.
	key, err := r.Export(accountName, passphrase)
	require.NoError(t, err)
	sdkaddr, err := a.Record.GetAddress()
	require.NoError(t, err)
	msg := &banktypes.MsgSend{
		FromAddress: sdkaddr.String(),
		ToAddress:   "cosmos1k8e50d2d8xkdfw9c4et3m45llh69e7xzw6uzga",
		Amount: sdktypes.NewCoins(
			sdktypes.NewCoin("token", math.NewIntFromUint64(1)),
		),
	}
	tests := []struct {
		name             string
		msg              sdktypes.Msg
		opts             []cosmosclient.Option
		expectedResponse *sdktypes.TxResponse
		expectedError    string
		setup            func(suite)
	}{
		{
			name:          "fail: error not found",
			msg:           msg,
			expectedError: "make sure that your account has enough balance",

			setup: func(s suite) {
				s.expectPrepareFactory(sdkaddr)
				s.signer.EXPECT().
					Sign(mock.Anything, mock.Anything, "bob", mock.Anything, true).
					Return(nil)
				s.rpcClient.EXPECT().
					BroadcastTxSync(mock.Anything, mock.Anything).
					Return(nil, sdkerrors.ErrNotFound)
			},
		},
		{
			name:          "fail: response code > 0",
			msg:           msg,
			expectedError: "error code: '42' msg: 'oups'",

			setup: func(s suite) {
				s.expectPrepareFactory(sdkaddr)
				s.signer.EXPECT().
					Sign(mock.Anything, mock.Anything, "bob", mock.Anything, true).
					Return(nil)
				s.rpcClient.EXPECT().
					BroadcastTxSync(mock.Anything, mock.Anything).
					Return(&ctypes.ResultBroadcastTx{
						Code: 42,
						Log:  "oups",
					}, nil)
			},
		},
		{
			name: "ok: tx confirmed immediately",
			msg:  msg,
			expectedResponse: &sdktypes.TxResponse{
				TxHash: txHashStr,
				RawLog: "log",
			},

			setup: func(s suite) {
				s.expectPrepareFactory(sdkaddr)
				s.signer.EXPECT().
					Sign(mock.Anything, mock.Anything, "bob", mock.Anything, true).
					Return(nil)
				s.rpcClient.EXPECT().
					BroadcastTxSync(mock.Anything, mock.Anything).
					Return(&ctypes.ResultBroadcastTx{
						Hash: txHash,
					}, nil)

				// Tx is broadcasted, now check for confirmation
				s.rpcClient.EXPECT().Tx(mock.Anything, txHash, false).
					Return(&ctypes.ResultTx{
						Hash: txHash,
						TxResult: abci.ExecTxResult{
							Log: "log",
						},
					}, nil)
			},
		},
		{
			name:          "fail: tx confirmed with error code",
			msg:           msg,
			expectedError: "error code: '42' msg: 'oups'",

			setup: func(s suite) {
				s.expectPrepareFactory(sdkaddr)
				s.signer.EXPECT().
					Sign(mock.Anything, mock.Anything, "bob", mock.Anything, true).
					Return(nil)
				s.rpcClient.EXPECT().
					BroadcastTxSync(mock.Anything, mock.Anything).
					Return(&ctypes.ResultBroadcastTx{
						Hash: txHash,
					}, nil)

				// Tx is broadcasted, now check for confirmation
				s.rpcClient.EXPECT().Tx(mock.Anything, txHash, false).
					Return(&ctypes.ResultTx{
						Hash: txHash,
						TxResult: abci.ExecTxResult{
							Code: 42,
							Log:  "oups",
						},
					}, nil)
			},
		},
		{
			name: "ok: tx confirmed after a while",
			msg:  msg,
			expectedResponse: &sdktypes.TxResponse{
				TxHash: txHashStr,
			},

			setup: func(s suite) {
				s.expectPrepareFactory(sdkaddr)
				s.signer.EXPECT().
					Sign(mock.Anything, mock.Anything, "bob", mock.Anything, true).
					Return(nil)
				s.rpcClient.EXPECT().
					BroadcastTxSync(mock.Anything, mock.Anything).
					Return(&ctypes.ResultBroadcastTx{
						Hash: txHash,
					}, nil)

				// Tx is broadcasted, now check for confirmation
				// First time the tx is not found (not confirmed yet)
				s.rpcClient.EXPECT().Tx(mock.Anything, txHash, false).
					Return(nil, errors.New("not found")).Once()
				// Wait for 1 block
				s.rpcClient.EXPECT().Status(mock.Anything).
					Return(&ctypes.ResultStatus{
						SyncInfo: ctypes.SyncInfo{LatestBlockHeight: 1},
					}, nil).Once()
				s.rpcClient.EXPECT().Status(mock.Anything).
					Return(&ctypes.ResultStatus{
						SyncInfo: ctypes.SyncInfo{LatestBlockHeight: 2},
					}, nil).Once()
				// Then try gain to fetch the tx, this time it is confirmed
				s.rpcClient.EXPECT().Tx(mock.Anything, txHash, false).
					Return(&ctypes.ResultTx{
						Hash: txHash,
					}, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newClient(t, tt.setup, tt.opts...)
			account, err := c.AccountRegistry.Import(accountName, key, passphrase)
			require.NoError(t, err)
			ctx := c.Context().
				WithFromName(accountName).
				WithFromAddress(sdkaddr)
			txService, err := c.CreateTx(ctx.CmdContext, account, tt.msg)
			require.NoError(t, err)

			res, err := txService.Broadcast(ctx.CmdContext)
			if tt.expectedError != "" {
				require.EqualError(t, err, tt.expectedError)
				return
			}
			require.NoError(t, err)
			require.Equal(t, ctx.Codec, res.Codec)
			require.Equal(t, tt.expectedResponse, res.TxResponse)
		})
	}
}

func TestTxServiceBroadcastAsync(t *testing.T) {
	var (
		accountName = "bob"
		passphrase  = "passphrase"
		txHash      = []byte{1, 2, 3}
		txHashStr   = hex.EncodeToString(txHash)
	)
	r, err := cosmosaccount.NewInMemory()
	require.NoError(t, err)
	a, _, err := r.Create(accountName)
	require.NoError(t, err)
	// Export created account to we can import it in the Client below.
	key, err := r.Export(accountName, passphrase)
	require.NoError(t, err)
	sdkaddr, err := a.Record.GetAddress()
	require.NoError(t, err)
	msg := &banktypes.MsgSend{
		FromAddress: sdkaddr.String(),
		ToAddress:   "cosmos1k8e50d2d8xkdfw9c4et3m45llh69e7xzw6uzga",
		Amount: sdktypes.NewCoins(
			sdktypes.NewCoin("token", math.NewIntFromUint64(1)),
		),
	}
	tests := []struct {
		name             string
		msg              sdktypes.Msg
		opts             []cosmosclient.Option
		expectedResponse *sdktypes.TxResponse
		expectedError    string
		setup            func(suite)
	}{
		{
			name:          "fail: error not found",
			msg:           msg,
			expectedError: "make sure that your account has enough balance",

			setup: func(s suite) {
				s.expectPrepareFactory(sdkaddr)
				s.signer.EXPECT().
					Sign(mock.Anything, mock.Anything, "bob", mock.Anything, true).
					Return(nil)
				s.rpcClient.EXPECT().
					BroadcastTxSync(mock.Anything, mock.Anything).
					Return(nil, sdkerrors.ErrNotFound)
			},
		},
		{
			name:          "fail: response code > 0",
			msg:           msg,
			expectedError: "error code: '42' msg: 'oups'",

			setup: func(s suite) {
				s.expectPrepareFactory(sdkaddr)
				s.signer.EXPECT().
					Sign(mock.Anything, mock.Anything, "bob", mock.Anything, true).
					Return(nil)
				s.rpcClient.EXPECT().
					BroadcastTxSync(mock.Anything, mock.Anything).
					Return(&ctypes.ResultBroadcastTx{
						Code: 42,
						Log:  "oups",
					}, nil)
			},
		},
		{
			name: "ok: tx broadcasted successfully",
			msg:  msg,
			expectedResponse: &sdktypes.TxResponse{
				TxHash: txHashStr,
				RawLog: "",
			},

			setup: func(s suite) {
				s.expectPrepareFactory(sdkaddr)
				s.signer.EXPECT().
					Sign(mock.Anything, mock.Anything, "bob", mock.Anything, true).
					Return(nil)
				s.rpcClient.EXPECT().
					BroadcastTxSync(mock.Anything, mock.Anything).
					Return(&ctypes.ResultBroadcastTx{
						Hash: txHash,
					}, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newClient(t, tt.setup, tt.opts...)
			account, err := c.AccountRegistry.Import(accountName, key, passphrase)
			require.NoError(t, err)
			ctx := c.Context().
				WithFromName(accountName).
				WithFromAddress(sdkaddr)
			txService, err := c.CreateTx(ctx.CmdContext, account, tt.msg)
			require.NoError(t, err)

			res, err := txService.BroadcastAsync(ctx.CmdContext)
			if tt.expectedError != "" {
				require.EqualError(t, err, tt.expectedError)
				return
			}
			require.NoError(t, err)
			require.Equal(t, ctx.Codec, res.Codec)
			require.Equal(t, tt.expectedResponse, res.TxResponse)
		})
	}
}
