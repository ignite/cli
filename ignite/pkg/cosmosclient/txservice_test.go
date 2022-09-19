package cosmosclient_test

import (
	"encoding/hex"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/ignite/cli/ignite/pkg/cosmosaccount"
	"github.com/ignite/cli/ignite/pkg/cosmosclient"
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
	sdkaddress, err := a.Record.GetAddress()
	require.NoError(t, err)
	msg := &banktypes.MsgSend{
		FromAddress: sdkaddress.String(),
		ToAddress:   "cosmos1k8e50d2d8xkdfw9c4et3m45llh69e7xzw6uzga",
		Amount: sdktypes.NewCoins(
			sdktypes.NewCoin("token", sdktypes.NewIntFromUint64((1))),
		),
	}
	tests := []struct {
		name             string
		msg              sdk.Msg
		opts             []cosmosclient.Option
		expectedResponse *sdktypes.TxResponse
		expectedError    string
		setup            func(suite)
	}{
		{
			name:          "fail: invalid msg",
			msg:           &banktypes.MsgSend{},
			expectedError: "invalid from address: empty address string is not allowed: invalid address",
			setup: func(s suite) {
				s.expectPrepareFactory(sdkaddress)
			},
		},
		{
			name:          "fail: error not found",
			msg:           msg,
			expectedError: "make sure that your account has enough balance",

			setup: func(s suite) {
				s.expectPrepareFactory(sdkaddress)
				s.signer.EXPECT().
					Sign(mock.Anything, "bob", mock.Anything, true).
					Return(nil)
				s.rpcClient.EXPECT().
					BroadcastTxCommit(mock.Anything, mock.Anything).
					Return(nil, sdkerrors.ErrNotFound)
			},
		},
		{
			name:          "fail: response code > 0",
			msg:           msg,
			expectedError: "error code: '42' msg: 'oups'",

			setup: func(s suite) {
				s.expectPrepareFactory(sdkaddress)
				s.signer.EXPECT().
					Sign(mock.Anything, "bob", mock.Anything, true).
					Return(nil)
				s.rpcClient.EXPECT().
					BroadcastTxCommit(mock.Anything, mock.Anything).
					Return(&ctypes.ResultBroadcastTxCommit{
						CheckTx: abci.ResponseCheckTx{
							Code: 42,
							Log:  "oups",
						},
					}, nil)
			},
		},
		{
			name:          "fail: ErrInsufficientFunds and disabled faucet",
			msg:           msg,
			expectedError: "error while requesting node 'http://localhost:26657': " + sdkerrors.ErrInsufficientFunds.Error(),

			setup: func(s suite) {
				s.expectPrepareFactory(sdkaddress)
				s.signer.EXPECT().
					Sign(mock.Anything, "bob", mock.Anything, true).
					Return(nil)
				s.rpcClient.EXPECT().
					BroadcastTxCommit(mock.Anything, mock.Anything).
					Return(nil, sdkerrors.ErrInsufficientFunds)
			},
		},
		{
			name:          "fail: ErrInsufficientFunds, enabled faucet but still not enough funds",
			msg:           msg,
			expectedError: "error while requesting node 'http://localhost:26657': " + sdkerrors.ErrInsufficientFunds.Error(),
			opts: []cosmosclient.Option{
				cosmosclient.WithUseFaucet("localhost:1234", "", 0),
			},

			setup: func(s suite) {
				s.expectMakeSureAccountHasToken(sdkaddress.String(), defaultFaucetMinAmount)
				s.expectPrepareFactory(sdkaddress)
				s.signer.EXPECT().
					Sign(mock.Anything, "bob", mock.Anything, true).
					Return(nil)
				s.rpcClient.EXPECT().
					BroadcastTxCommit(mock.Anything, mock.Anything).
					Return(nil, sdkerrors.ErrInsufficientFunds).
					Once()

				s.expectMakeSureAccountHasToken(sdkaddress.String(), defaultFaucetMinAmount)

				// Once balance is fine, broadcast the tx again, but still no funds
				s.rpcClient.EXPECT().
					BroadcastTxCommit(mock.Anything, mock.Anything).
					Return(nil, sdkerrors.ErrInsufficientFunds).
					Once()
			},
		},
		{
			name: "ok: basic usecase",
			msg:  msg,
			expectedResponse: &sdktypes.TxResponse{
				TxHash: txHashStr,
			},

			setup: func(s suite) {
				s.expectPrepareFactory(sdkaddress)
				s.signer.EXPECT().
					Sign(mock.Anything, "bob", mock.Anything, true).
					Return(nil)
				s.rpcClient.EXPECT().
					BroadcastTxCommit(mock.Anything, mock.Anything).
					Return(&ctypes.ResultBroadcastTxCommit{
						Hash: txHash,
					}, nil)
			},
		},
		{
			name: "ok: ErrInsufficientFunds and enabled faucet",
			msg:  msg,
			opts: []cosmosclient.Option{
				cosmosclient.WithUseFaucet("localhost:1234", "", 0),
			},
			expectedResponse: &sdktypes.TxResponse{
				TxHash: txHashStr,
			},

			setup: func(s suite) {
				s.expectMakeSureAccountHasToken(sdkaddress.String(), defaultFaucetMinAmount)
				s.expectPrepareFactory(sdkaddress)
				s.signer.EXPECT().
					Sign(mock.Anything, "bob", mock.Anything, true).
					Return(nil)
				s.rpcClient.EXPECT().
					BroadcastTxCommit(mock.Anything, mock.Anything).
					Return(nil, sdkerrors.ErrInsufficientFunds).
					Once()

				s.expectMakeSureAccountHasToken(sdkaddress.String(), defaultFaucetMinAmount)

				// Once balance is fine, broadcast the tx again
				s.rpcClient.EXPECT().
					BroadcastTxCommit(mock.Anything, mock.Anything).
					Return(&ctypes.ResultBroadcastTxCommit{
						Hash: txHash,
					}, nil).
					Once()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require := require.New(t)
			assert := assert.New(t)
			c := newClient(t, tt.setup, tt.opts...)
			account, err := c.AccountRegistry.Import(accountName, key, passphrase)
			require.NoError(err)
			ctx := c.Context().
				WithFromName(accountName).
				WithFromAddress(sdkaddress)
			txService, err := c.CreateTx(account, tt.msg)
			require.NoError(err)

			res, err := txService.Broadcast()

			if tt.expectedError != "" {
				require.EqualError(err, tt.expectedError)
				return
			}
			require.NoError(err)
			assert.Equal(ctx.Codec, res.Codec)
			assert.Equal(tt.expectedResponse, res.TxResponse)
		})
	}
}
