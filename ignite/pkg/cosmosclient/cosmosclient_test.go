package cosmosclient

import (
	"context"
	"io"
	"os"
	"path"
	"strconv"
	"testing"

	"github.com/cosmos/cosmos-sdk/client/flags"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/p2p"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/ignite/cli/ignite/pkg/cosmosaccount"
	"github.com/ignite/cli/ignite/pkg/cosmosclient/mocks"
)

//go:generate mockery --srcpkg github.com/tendermint/tendermint/rpc/client/ --name Client --structname RPCClient --filename rpclient.go --with-expecter

type suite struct {
	rpcClient *mocks.RPCClient
}

func newSuite(t *testing.T, setup func(suite)) suite {
	s := suite{
		rpcClient: mocks.NewRPCClient(t),
	}
	if setup != nil {
		setup(s)
	}
	return s
}

func TestNew(t *testing.T) {
	var (
		ctx     = context.Background()
		home, _ = os.UserHomeDir()
	)
	tests := []struct {
		name           string
		opts           []Option
		expectedClient Client
		expectedError  string
		setup          func(suite)
	}{
		{
			name: "default values",
			setup: func(s suite) {
				s.rpcClient.EXPECT().Status(ctx).Return(&ctypes.ResultStatus{
					NodeInfo: p2p.DefaultNodeInfo{Network: "mychain"},
				}, nil)
			},
			expectedClient: Client{
				chainID:            "mychain",
				nodeAddress:        defaultNodeAddress,
				homePath:           path.Join(home, ".mychain"),
				keyringServiceName: "",
				keyringDir:         path.Join(home, ".mychain"),
				keyringBackend:     cosmosaccount.KeyringTest,
				addressPrefix:      "cosmos",
				faucetAddress:      defaultFaucetAddress,
				faucetDenom:        defaultFaucetDenom,
				faucetMinAmount:    defaultFaucetMinAmount,
				out:                io.Discard,
				gas:                strconv.Itoa(defaultGasLimit),
				broadcastMode:      flags.BroadcastBlock,
			},
		},
		{
			name: "custom values",
			setup: func(s suite) {
				s.rpcClient.EXPECT().Status(ctx).Return(&ctypes.ResultStatus{
					NodeInfo: p2p.DefaultNodeInfo{Network: "mychain"},
				}, nil)
			},
			opts: []Option{
				WithHome("home"),
				WithKeyringServiceName("keyringServiceName"),
				WithKeyringBackend(cosmosaccount.KeyringOS),
				WithKeyringDir("keyringDir"),
				WithNodeAddress("addr"),
				WithAddressPrefix("prefix"),
				WithUseFaucet("faucetAddress", "denom", 42),
				WithGas("gas"),
				WithGasPrices("gasPrices"),
				WithFees("fees"),
				WithBroadcastMode("broadcastMode"),
				WithGenerateOnly(true),
			},
			expectedClient: Client{
				chainID:            "mychain",
				nodeAddress:        "addr",
				homePath:           "home",
				keyringServiceName: "keyringServiceName",
				keyringBackend:     cosmosaccount.KeyringOS,
				keyringDir:         "keyringDir",
				addressPrefix:      "prefix",
				useFaucet:          true,
				faucetAddress:      "faucetAddress",
				faucetDenom:        "denom",
				faucetMinAmount:    42,
				out:                io.Discard,
				gas:                "gas",
				gasPrices:          "gasPrices",
				fees:               "fees",
				broadcastMode:      "broadcastMode",
				generateOnly:       true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				require = require.New(t)
				assert  = assert.New(t)
				suite   = newSuite(t, tt.setup)
			)
			tt.opts = append(tt.opts,
				WithRPCClient(suite.rpcClient),
			)

			c, err := New(ctx, tt.opts...)

			if tt.expectedError != "" {
				require.EqualError(err, tt.expectedError)
				return
			}
			require.NoError(err)
			assert.Equal(tt.expectedClient.chainID, "mychain")
			assert.Equal(tt.expectedClient.nodeAddress, c.nodeAddress)
			assert.Equal(tt.expectedClient.homePath, c.homePath)
			assert.Equal(tt.expectedClient.keyringServiceName, c.keyringServiceName)
			assert.Equal(tt.expectedClient.keyringBackend, c.keyringBackend)
			assert.Equal(tt.expectedClient.keyringDir, c.keyringDir)
			assert.Equal(tt.expectedClient.addressPrefix, c.addressPrefix)
			assert.Equal(tt.expectedClient.useFaucet, c.useFaucet)
			assert.Equal(tt.expectedClient.faucetAddress, c.faucetAddress)
			assert.Equal(tt.expectedClient.faucetDenom, c.faucetDenom)
			assert.Equal(tt.expectedClient.faucetMinAmount, c.faucetMinAmount)
			assert.Equal(tt.expectedClient.out, c.out)
			assert.Equal(tt.expectedClient.gas, c.gas)
			assert.Equal(tt.expectedClient.fees, c.fees)
			assert.Equal(tt.expectedClient.gasPrices, c.gasPrices)
			assert.Equal(tt.expectedClient.broadcastMode, c.broadcastMode)
			assert.Equal(tt.expectedClient.generateOnly, c.generateOnly)
			// assert the sdk config has been updated with the addr prefix
			config := sdktypes.GetConfig()
			assert.Equal(tt.expectedClient.addressPrefix, config.GetBech32AccountAddrPrefix())
		})
	}
}
