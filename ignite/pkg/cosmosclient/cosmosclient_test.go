package cosmosclient

import (
	"context"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	rpcmocks "github.com/tendermint/tendermint/rpc/client/mocks"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

func TestGetBlockTXs(t *testing.T) {
	m := rpcmocks.Client{}
	m.Test(t)

	// Mock the Block RPC endpoint
	ctx := context.Background()
	block := createTestBlock(1)

	m.On("Block", ctx, &block.Height).Return(&ctypes.ResultBlock{Block: &block}, nil)

	// Mock the TxSearch RPC endpoint
	searchQry := createTxSearchByHeightQuery(block.Height)
	page := 1
	perPage := defaultTXsPerPage
	rtx := ctypes.ResultTx{}
	resSearch := ctypes.ResultTxSearch{
		Txs:        []*ctypes.ResultTx{&rtx},
		TotalCount: 1,
	}

	m.On("TxSearch", ctx, searchQry, false, &page, &perPage, orderAsc).Return(&resSearch, nil)

	// Create a cosmos client with an RPC mock
	client := Client{RPC: &m}

	txs, err := client.GetBlockTXs(ctx, block.Height)
	require.NoError(t, err)
	require.Equal(t, txs, []TX{
		{
			BlockTime: block.Time,
			Raw:       &rtx,
		},
	})

	m.AssertNumberOfCalls(t, "Block", 1)
	m.AssertNumberOfCalls(t, "TxSearch", 1)
}

func TestGetBlockTXsWithBlockError(t *testing.T) {
	m := rpcmocks.Client{}
	m.Test(t)

	wantErr := errors.New("expected error")

	// Mock the Block RPC endpoint
	ctx := context.Background()
	height := int64(1)

	m.On("Block", ctx, &height).Return(nil, wantErr)

	// Create a cosmos client with an RPC mock
	client := Client{RPC: &m}

	txs, err := client.GetBlockTXs(ctx, height)
	require.ErrorIs(t, err, wantErr)
	require.Nil(t, txs)

	m.AssertNumberOfCalls(t, "Block", 1)
	m.AssertNumberOfCalls(t, "TxSearch", 0)
}

func TestGetBlockTXsPagination(t *testing.T) {
	m := rpcmocks.Client{}
	m.Test(t)

	// Mock the Block RPC endpoint
	ctx := context.Background()
	block := createTestBlock(1)

	m.On("Block", ctx, &block.Height).Return(&ctypes.ResultBlock{Block: &block}, nil)

	// Mock the TxSearch RPC endpoint and fake the number of
	// transactions so it is called twice to fetch two pages
	searchQry := createTxSearchByHeightQuery(block.Height)
	perPage := defaultTXsPerPage
	fakeCount := perPage + 1
	first := 1
	second := 2
	firstPage := ctypes.ResultTxSearch{
		Txs:        []*ctypes.ResultTx{{}},
		TotalCount: fakeCount,
	}
	secondPage := ctypes.ResultTxSearch{
		Txs:        []*ctypes.ResultTx{{}},
		TotalCount: fakeCount,
	}

	m.On("TxSearch", ctx, searchQry, false, &first, &perPage, orderAsc).Return(&firstPage, nil)
	m.On("TxSearch", ctx, searchQry, false, &second, &perPage, orderAsc).Return(&secondPage, nil)

	// Create a cosmos client with an RPC mock
	client := Client{RPC: &m}

	txs, err := client.GetBlockTXs(ctx, block.Height)
	require.NoError(t, err)
	require.Equal(t, txs, []TX{
		{
			BlockTime: block.Time,
			Raw:       firstPage.Txs[0],
		},
		{
			BlockTime: block.Time,
			Raw:       secondPage.Txs[0],
		},
	})

	m.AssertNumberOfCalls(t, "Block", 1)
	m.AssertNumberOfCalls(t, "TxSearch", 2)
}

func TestGetBlockTXsWithSearchError(t *testing.T) {
	m := rpcmocks.Client{}
	m.Test(t)

	wantErr := errors.New("expected error")

	// Mock the Block RPC endpoint
	ctx := context.Background()
	block := createTestBlock(1)

	m.On("Block", ctx, &block.Height).Return(&ctypes.ResultBlock{Block: &block}, nil)

	// Mock the TxSearch RPC endpoint
	searchQry := createTxSearchByHeightQuery(block.Height)
	perPage := defaultTXsPerPage
	page := 1

	m.On("TxSearch", ctx, searchQry, false, &page, &perPage, orderAsc).Return(nil, wantErr)

	// Create a cosmos client with an RPC mock
	client := Client{RPC: &m}

	txs, err := client.GetBlockTXs(ctx, block.Height)
	require.ErrorIs(t, err, wantErr)
	require.Nil(t, txs)

	m.AssertNumberOfCalls(t, "Block", 1)
	m.AssertNumberOfCalls(t, "TxSearch", 1)
}

func createTestBlock(height int64) tmtypes.Block {
	return tmtypes.Block{
		Header: tmtypes.Header{
			Height: height,
			Time:   time.Now(),
		},
	}
}
