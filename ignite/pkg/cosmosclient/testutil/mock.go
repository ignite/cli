package testutil

import (
	"testing"

	"github.com/ignite/cli/ignite/pkg/cosmosclient/mocks"
	"github.com/stretchr/testify/mock"
)

// NewTendermintClientMock creates a new Tendermint RPC client mock.
func NewTendermintClientMock(t *testing.T) *TendermintClientMock {
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

// OnStatus starts a generic call mock on the Block RPC method.
func (m *TendermintClientMock) OnBlock() *mock.Call {
	return m.On("Block", RepeatMockArgs(2)...)
}

// OnStatus starts a generic call mock on the TxSearch RPC method.
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
