package cosmosclient

import (
	"context"
	"fmt"

	"github.com/tendermint/tendermint/libs/bytes"
	"github.com/tendermint/tendermint/rpc/client"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	"github.com/tendermint/tendermint/types"
)

// rpcWrapper is a rpclient.Client but with more contextualized errors.
// Useful because the original implementation may return JSON errors when the
// requested node is busy, which is confusing for the user. With rpcWrapper,
// the error is prefixed with 'error while requesting node xxx: JSON error'.
type rpcWrapper struct {
	rpcclient.Client
	nodeAddress string
}

func rpcError(node string, err error) error {
	return fmt.Errorf("error while requesting node '%s': %w", node, err)
}

// Reading from abci app
func (rpc rpcWrapper) ABCIInfo(ctx context.Context) (*ctypes.ResultABCIInfo, error) {
	res, err := rpc.Client.ABCIInfo(ctx)
	if err != nil {
		return nil, rpcError(rpc.nodeAddress, err)
	}
	return res, nil
}

func (rpc rpcWrapper) ABCIQuery(ctx context.Context, path string, data bytes.HexBytes) (*ctypes.ResultABCIQuery, error) {
	res, err := rpc.Client.ABCIQuery(ctx, path, data)
	if err != nil {
		return nil, rpcError(rpc.nodeAddress, err)
	}
	return res, nil
	// TODO: Implement
}

func (rpc rpcWrapper) ABCIQueryWithOptions(ctx context.Context, path string, data bytes.HexBytes, opts client.ABCIQueryOptions) (*ctypes.ResultABCIQuery, error) {
	res, err := rpc.Client.ABCIQueryWithOptions(ctx, path, data, opts)
	if err != nil {
		return nil, rpcError(rpc.nodeAddress, err)
	}
	return res, nil
}

// Writing to abci app
func (rpc rpcWrapper) BroadcastTxCommit(ctx context.Context, tx types.Tx) (*ctypes.ResultBroadcastTxCommit, error) {
	res, err := rpc.Client.BroadcastTxCommit(ctx, tx)
	if err != nil {
		return nil, rpcError(rpc.nodeAddress, err)
	}
	return res, nil
}

func (rpc rpcWrapper) BroadcastTxAsync(ctx context.Context, tx types.Tx) (*ctypes.ResultBroadcastTx, error) {
	res, err := rpc.Client.BroadcastTxAsync(ctx, tx)
	if err != nil {
		return nil, rpcError(rpc.nodeAddress, err)
	}
	return res, nil
}

func (rpc rpcWrapper) BroadcastTxSync(ctx context.Context, tx types.Tx) (*ctypes.ResultBroadcastTx, error) {
	res, err := rpc.Client.BroadcastTxSync(ctx, tx)
	if err != nil {
		return nil, rpcError(rpc.nodeAddress, err)
	}
	return res, nil
}

func (rpc rpcWrapper) GenesisChunked(ctx context.Context, n uint) (*ctypes.ResultGenesisChunk, error) {
	res, err := rpc.Client.GenesisChunked(ctx, n)
	if err != nil {
		return nil, rpcError(rpc.nodeAddress, err)
	}
	return res, nil
}

func (rpc rpcWrapper) BlockchainInfo(ctx context.Context, minHeight int64, maxHeight int64) (*ctypes.ResultBlockchainInfo, error) {
	res, err := rpc.Client.BlockchainInfo(ctx, minHeight, maxHeight)
	if err != nil {
		return nil, rpcError(rpc.nodeAddress, err)
	}
	return res, nil
}

func (rpc rpcWrapper) NetInfo(ctx context.Context) (*ctypes.ResultNetInfo, error) {
	res, err := rpc.Client.NetInfo(ctx)
	if err != nil {
		return nil, rpcError(rpc.nodeAddress, err)
	}
	return res, nil
}

func (rpc rpcWrapper) DumpConsensusState(ctx context.Context) (*ctypes.ResultDumpConsensusState, error) {
	res, err := rpc.Client.DumpConsensusState(ctx)
	if err != nil {
		return nil, rpcError(rpc.nodeAddress, err)
	}
	return res, nil
}

func (rpc rpcWrapper) ConsensusState(ctx context.Context) (*ctypes.ResultConsensusState, error) {
	res, err := rpc.Client.ConsensusState(ctx)
	if err != nil {
		return nil, rpcError(rpc.nodeAddress, err)
	}
	return res, nil
}

func (rpc rpcWrapper) ConsensusParams(ctx context.Context, height *int64) (*ctypes.ResultConsensusParams, error) {
	res, err := rpc.Client.ConsensusParams(ctx, height)
	if err != nil {
		return nil, rpcError(rpc.nodeAddress, err)
	}
	return res, nil
}

func (rpc rpcWrapper) Health(ctx context.Context) (*ctypes.ResultHealth, error) {
	res, err := rpc.Client.Health(ctx)
	if err != nil {
		return nil, rpcError(rpc.nodeAddress, err)
	}
	return res, nil
}

func (rpc rpcWrapper) Block(ctx context.Context, height *int64) (*ctypes.ResultBlock, error) {
	res, err := rpc.Client.Block(ctx, height)
	if err != nil {
		return nil, rpcError(rpc.nodeAddress, err)
	}
	return res, nil
}

func (rpc rpcWrapper) BlockByHash(ctx context.Context, hash []byte) (*ctypes.ResultBlock, error) {
	res, err := rpc.Client.BlockByHash(ctx, hash)
	if err != nil {
		return nil, rpcError(rpc.nodeAddress, err)
	}
	return res, nil
}

func (rpc rpcWrapper) BlockResults(ctx context.Context, height *int64) (*ctypes.ResultBlockResults, error) {
	res, err := rpc.Client.BlockResults(ctx, height)
	if err != nil {
		return nil, rpcError(rpc.nodeAddress, err)
	}
	return res, nil
}

func (rpc rpcWrapper) Commit(ctx context.Context, height *int64) (*ctypes.ResultCommit, error) {
	res, err := rpc.Client.Commit(ctx, height)
	if err != nil {
		return nil, rpcError(rpc.nodeAddress, err)
	}
	return res, nil
}

func (rpc rpcWrapper) Validators(ctx context.Context, height *int64, page *int, perPage *int) (*ctypes.ResultValidators, error) {
	res, err := rpc.Client.Validators(ctx, height, page, perPage)
	if err != nil {
		return nil, rpcError(rpc.nodeAddress, err)
	}
	return res, nil
}

func (rpc rpcWrapper) Tx(ctx context.Context, hash []byte, prove bool) (*ctypes.ResultTx, error) {
	res, err := rpc.Client.Tx(ctx, hash, prove)
	if err != nil {
		return nil, rpcError(rpc.nodeAddress, err)
	}
	return res, nil
}

func (rpc rpcWrapper) TxSearch(ctx context.Context, query string, prove bool, page *int, perPage *int, orderBy string) (*ctypes.ResultTxSearch, error) {
	res, err := rpc.Client.TxSearch(ctx, query, prove, page, perPage, orderBy)
	if err != nil {
		return nil, rpcError(rpc.nodeAddress, err)
	}
	return res, nil
}

func (rpc rpcWrapper) BlockSearch(ctx context.Context, query string, page *int, perPage *int, orderBy string) (*ctypes.ResultBlockSearch, error) {
	res, err := rpc.Client.BlockSearch(ctx, query, page, perPage, orderBy)
	if err != nil {
		return nil, rpcError(rpc.nodeAddress, err)
	}
	return res, nil
}

func (rpc rpcWrapper) Status(ctx context.Context) (*ctypes.ResultStatus, error) {
	res, err := rpc.Client.Status(ctx)
	if err != nil {
		return nil, rpcError(rpc.nodeAddress, err)
	}
	return res, nil
}

func (rpc rpcWrapper) BroadcastEvidence(ctx context.Context, e types.Evidence) (*ctypes.ResultBroadcastEvidence, error) {
	res, err := rpc.Client.BroadcastEvidence(ctx, e)
	if err != nil {
		return nil, rpcError(rpc.nodeAddress, err)
	}
	return res, nil
}

func (rpc rpcWrapper) UnconfirmedTxs(ctx context.Context, limit *int) (*ctypes.ResultUnconfirmedTxs, error) {
	res, err := rpc.Client.UnconfirmedTxs(ctx, limit)
	if err != nil {
		return nil, rpcError(rpc.nodeAddress, err)
	}
	return res, nil
}

func (rpc rpcWrapper) NumUnconfirmedTxs(ctx context.Context) (*ctypes.ResultUnconfirmedTxs, error) {
	res, err := rpc.Client.NumUnconfirmedTxs(ctx)
	if err != nil {
		return nil, rpcError(rpc.nodeAddress, err)
	}
	return res, nil
}

func (rpc rpcWrapper) CheckTx(ctx context.Context, tx types.Tx) (*ctypes.ResultCheckTx, error) {
	res, err := rpc.Client.CheckTx(ctx, tx)
	if err != nil {
		return nil, rpcError(rpc.nodeAddress, err)
	}
	return res, nil
}
