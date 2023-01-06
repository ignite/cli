package testutil

import (
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/ignite/cli/ignite/pkg/cosmosclient/mocks"
)

//go:generate mockery --srcpkg github.com/tendermint/tendermint/rpc/client --name Client --structname RPCClient --filename rpc_client.go --output ../mocks --with-expecter
//go:generate mockery --srcpkg github.com/cosmos/cosmos-sdk/client --name AccountRetriever --filename account_retriever.go --output ../mocks --with-expecter
//go:generate mockery --srcpkg github.com/cosmos/cosmos-sdk/x/bank/types --name QueryClient --structname BankQueryClient --filename bank_query_client.go --output ../mocks --with-expecter

// NewTendermintClientMock creates a new Tendermint RPC client mock.
func NewTendermintClientMock(t *testing.T) *TendermintClientMock {
	t.Helper()
	m := TendermintClientMock{}
	m.Test(t)

	return &m
}

// TendermintClientMock mocks Tendermint's RPC client.
type TendermintClientMock struct {
	mocks.RPCClient
}

// OnStatus starts a generic call mock on the Status RPC method.
func (m *TendermintClientMock) OnStatus() *mock.Call {
	return m.On("Status", mock.Anything)
}

// OnBlock starts a generic call mock on the Block RPC method.
func (m *TendermintClientMock) OnBlock() *mock.Call {
	return m.On("Block", RepeatMockArgs(2)...)
}

// OnTxSearch starts a generic call mock on the TxSearch RPC method.
func (m *TendermintClientMock) OnTxSearch() *mock.Call {
	return m.On("TxSearch", RepeatMockArgs(6)...)
}

// RepeatMockArgs returns a slice with an N number of mock.Anything arguments.
// This function can be useful to define a number of generic consecutive arguments
// for mocked method calls.
func RepeatMockArgs(n int) (args []interface{}) {
	for i := 0; i < n; i++ {
		args = append(args, mock.Anything)
	}

	return args
}
