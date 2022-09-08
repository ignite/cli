package cosmosclient_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/p2p"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/ignite/cli/ignite/pkg/cosmosaccount"
	"github.com/ignite/cli/ignite/pkg/cosmosclient"
	"github.com/ignite/cli/ignite/pkg/cosmosclient/mocks"
)

//go:generate mockery --srcpkg github.com/tendermint/tendermint/rpc/client/ --name Client --structname RPCClient --filename rpclient.go --with-expecter

type suite struct {
	rpcClient *mocks.RPCClient
}

func newClient(t *testing.T, setup func(suite), opts ...cosmosclient.Option) cosmosclient.Client {
	s := suite{
		rpcClient: mocks.NewRPCClient(t),
	}
	s.rpcClient.EXPECT().Status(mock.Anything).
		Return(&ctypes.ResultStatus{}, nil).Once()
	if setup != nil {
		setup(s)
	}
	opts = append(opts, []cosmosclient.Option{
		cosmosclient.WithRPCClient(s.rpcClient),
		cosmosclient.WithKeyringBackend(cosmosaccount.KeyringMemory),
	}...)
	c, err := cosmosclient.New(context.Background(), opts...)
	require.NoError(t, err)
	return c
}

func TestClientWaitForBlockHeight(t *testing.T) {
	var (
		ctx                 = context.Background()
		canceledCtx, cancel = context.WithTimeout(ctx, 0)
		targetBlockHeight   = int64(42)
	)
	cancel()
	tests := []struct {
		name              string
		ctx               context.Context
		waitBlockDuration time.Duration
		expectedError     string
		setup             func(suite)
	}{
		{
			name: "ok: no wait",
			ctx:  ctx,
			setup: func(s suite) {
				s.rpcClient.EXPECT().Status(ctx).Return(&ctypes.ResultStatus{
					SyncInfo: ctypes.SyncInfo{LatestBlockHeight: targetBlockHeight},
				}, nil)
			},
		},
		{
			name:              "ok: wait 1 time",
			ctx:               ctx,
			waitBlockDuration: time.Second * 2, // must exceed the wait loop duration
			setup: func(s suite) {
				s.rpcClient.EXPECT().Status(ctx).Return(&ctypes.ResultStatus{
					SyncInfo: ctypes.SyncInfo{LatestBlockHeight: targetBlockHeight - 1},
				}, nil).Once()
				s.rpcClient.EXPECT().Status(ctx).Return(&ctypes.ResultStatus{
					SyncInfo: ctypes.SyncInfo{LatestBlockHeight: targetBlockHeight},
				}, nil).Once()
			},
		},
		{
			name:              "fail: wait expired",
			ctx:               ctx,
			waitBlockDuration: time.Millisecond,
			expectedError:     "timeout exceeded waiting for block",
			setup: func(s suite) {
				s.rpcClient.EXPECT().Status(ctx).Return(&ctypes.ResultStatus{
					SyncInfo: ctypes.SyncInfo{LatestBlockHeight: targetBlockHeight - 1},
				}, nil)
			},
		},
		{
			name:              "fail: canceled context",
			ctx:               canceledCtx,
			waitBlockDuration: time.Millisecond,
			expectedError:     canceledCtx.Err().Error(),
			setup: func(s suite) {
				s.rpcClient.EXPECT().Status(canceledCtx).Return(&ctypes.ResultStatus{
					SyncInfo: ctypes.SyncInfo{LatestBlockHeight: targetBlockHeight - 1},
				}, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require := require.New(t)
			c := newClient(t, tt.setup)

			err := c.WaitForBlockHeight(tt.ctx, targetBlockHeight, tt.waitBlockDuration)

			if tt.expectedError != "" {
				require.EqualError(err, tt.expectedError)
				return
			}
			require.NoError(err)
		})
	}
}

func TestClientAccount(t *testing.T) {
	var (
		accountName = "bob"
		passphrase  = "passphrase"
	)
	r, err := cosmosaccount.NewInMemory()
	require.NoError(t, err)
	expectedAccount, _, err := r.Create(accountName)
	require.NoError(t, err)
	expectedAddr, err := expectedAccount.Address("cosmos")
	require.NoError(t, err)
	// Export created account to we can import it in the Client below.
	key, err := r.Export(accountName, passphrase)
	require.NoError(t, err)

	tests := []struct {
		name          string
		addressOrName string
		expectedError string
	}{
		{
			name:          "ok: find by name",
			addressOrName: expectedAccount.Name,
		},
		{
			name:          "ok: find by address",
			addressOrName: expectedAddr,
		},
		{
			name:          "fail: name not found",
			addressOrName: "unknown",
			expectedError: "decoding bech32 failed: invalid bech32 string length 7",
		},
		{
			name:          "fail: address not found",
			addressOrName: "cosmos1cs4hpwrpna6ucsgsa78jfp403l7gdynukrxkrv",
			expectedError: `account "cosmos1cs4hpwrpna6ucsgsa78jfp403l7gdynukrxkrv" does not exist`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				require = require.New(t)
				assert  = assert.New(t)
				c       = newClient(t, nil)
			)
			_, err := c.AccountRegistry.Import(accountName, key, passphrase)
			require.NoError(err)

			account, err := c.Account(tt.addressOrName)

			if tt.expectedError != "" {
				require.EqualError(err, tt.expectedError)
				return
			}
			require.NoError(err)
			assert.Equal(expectedAccount, account)
		})
	}
}

func TestClientAddress(t *testing.T) {
	var (
		accountName = "bob"
		passphrase  = "passphrase"
	)
	r, err := cosmosaccount.NewInMemory()
	require.NoError(t, err)
	expectedAccount, _, err := r.Create(accountName)
	require.NoError(t, err)
	// Export created account to we can import it in the Client below.
	key, err := r.Export(accountName, passphrase)
	require.NoError(t, err)

	tests := []struct {
		name           string
		accountName    string
		opts           []cosmosclient.Option
		expectedError  string
		expectedPrefix string
	}{
		{
			name:           "ok: name exists",
			accountName:    expectedAccount.Name,
			expectedPrefix: "cosmos",
		},
		{
			name: "ok: name exists with different prefix",
			opts: []cosmosclient.Option{
				cosmosclient.WithAddressPrefix("test"),
			},
			accountName:    expectedAccount.Name,
			expectedPrefix: "test",
		},
		{
			name:          "fail: name not found",
			accountName:   "unknown",
			expectedError: `account "unknown" does not exist`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				require = require.New(t)
				assert  = assert.New(t)
				c       = newClient(t, nil, tt.opts...)
			)
			_, err := c.AccountRegistry.Import(accountName, key, passphrase)
			require.NoError(err)

			address, err := c.Address(tt.accountName)

			if tt.expectedError != "" {
				require.EqualError(err, tt.expectedError)
				return
			}
			require.NoError(err)
			expectedAddr, err := expectedAccount.Address(tt.expectedPrefix)
			require.NoError(err)
			assert.Equal(expectedAddr, address)
		})
	}
}

func TestClientStatus(t *testing.T) {
	var (
		ctx            = context.Background()
		expectedStatus = &ctypes.ResultStatus{
			NodeInfo: p2p.DefaultNodeInfo{Network: "mychain"},
		}
	)
	c := newClient(t, func(s suite) {
		s.rpcClient.EXPECT().Status(ctx).Return(expectedStatus, nil).Once()
	})

	status, err := c.Status(ctx)

	require.NoError(t, err)
	assert.Equal(t, expectedStatus, status)
}
