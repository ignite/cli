package cosmosclient

import (
	"context"

	"github.com/cometbft/cometbft/libs/bytes"
	rpcclient "github.com/cometbft/cometbft/rpc/client"
	ctypes "github.com/cometbft/cometbft/rpc/core/types"
	"github.com/cometbft/cometbft/types"
	"github.com/pkg/errors"
)

// rpcWrapper is a rpclient.Client but with more contextualized errors.
// Useful because the original implementation may return JSON errors when the
// requested node is busy, which is confusing for the user. With rpcWrapper,
// the error is prefixed with 'error while requesting node xxx: JSON error'.
// TODO(tb): we may remove this wrapper once https://github.com/tendermint/tendermint/issues/9312 is fixed.
type rpcWrapper struct {
	rpcclient.Client
	nodeAddress string
}

func rpcError(node string, err error) error {
	return errors.Wrapf(err, "error while requesting node '%s'", node)
}

func (rpc rpcWrapper) ABCIInfo(ctx context.Context) (*ctypes.ResultABCIInfo, error) {
	res, err := rpc.Client.ABCIInfo(ctx)
	return res, rpcError(rpc.nodeAddress, err)
}

func (rpc rpcWrapper) ABCIQuery(ctx context.Context, path string, data bytes.HexBytes) (*ctypes.ResultABCIQuery, error) {
	res, err := rpc.Client.ABCIQuery(ctx, path, data)
	return res, rpcError(rpc.nodeAddress, err)
}

func (rpc rpcWrapper) ABCIQueryWithOptions(ctx context.Context, path string, data bytes.HexBytes, opts rpcclient.ABCIQueryOptions) (*ctypes.ResultABCIQuery, error) {
	res, err := rpc.Client.ABCIQueryWithOptions(ctx, path, data, opts)
	return res, rpcError(rpc.nodeAddress, err)
}

func (rpc rpcWrapper) BroadcastTxCommit(ctx context.Context, tx types.Tx) (*ctypes.ResultBroadcastTxCommit, error) {
	res, err := rpc.Client.BroadcastTxCommit(ctx, tx)
	return res, rpcError(rpc.nodeAddress, err)
}

func (rpc rpcWrapper) BroadcastTxAsync(ctx context.Context, tx types.Tx) (*ctypes.ResultBroadcastTx, error) {
	res, err := rpc.Client.BroadcastTxAsync(ctx, tx)
	return res, rpcError(rpc.nodeAddress, err)
}

func (rpc rpcWrapper) BroadcastTxSync(ctx context.Context, tx types.Tx) (*ctypes.ResultBroadcastTx, error) {
	res, err := rpc.Client.BroadcastTxSync(ctx, tx)
	return res, rpcError(rpc.nodeAddress, err)
}

func (rpc rpcWrapper) GenesisChunked(ctx context.Context, n uint) (*ctypes.ResultGenesisChunk, error) {
	res, err := rpc.Client.GenesisChunked(ctx, n)
	return res, rpcError(rpc.nodeAddress, err)
}

func (rpc rpcWrapper) BlockchainInfo(ctx context.Context, minHeight int64, maxHeight int64) (*ctypes.ResultBlockchainInfo, error) {
	res, err := rpc.Client.BlockchainInfo(ctx, minHeight, maxHeight)
	return res, rpcError(rpc.nodeAddress, err)
}

func (rpc rpcWrapper) NetInfo(ctx context.Context) (*ctypes.ResultNetInfo, error) {
	res, err := rpc.Client.NetInfo(ctx)
	return res, rpcError(rpc.nodeAddress, err)
}

func (rpc rpcWrapper) DumpConsensusState(ctx context.Context) (*ctypes.ResultDumpConsensusState, error) {
	res, err := rpc.Client.DumpConsensusState(ctx)
	return res, rpcError(rpc.nodeAddress, err)
}

func (rpc rpcWrapper) ConsensusState(ctx context.Context) (*ctypes.ResultConsensusState, error) {
	res, err := rpc.Client.ConsensusState(ctx)
	return res, rpcError(rpc.nodeAddress, err)
}

func (rpc rpcWrapper) ConsensusParams(ctx context.Context, height *int64) (*ctypes.ResultConsensusParams, error) {
	res, err := rpc.Client.ConsensusParams(ctx, height)
	return res, rpcError(rpc.nodeAddress, err)
}

func (rpc rpcWrapper) Health(ctx context.Context) (*ctypes.ResultHealth, error) {
	res, err := rpc.Client.Health(ctx)
	return res, rpcError(rpc.nodeAddress, err)
}

func (rpc rpcWrapper) Block(ctx context.Context, height *int64) (*ctypes.ResultBlock, error) {
	res, err := rpc.Client.Block(ctx, height)
	return res, rpcError(rpc.nodeAddress, err)
}

func (rpc rpcWrapper) BlockByHash(ctx context.Context, hash []byte) (*ctypes.ResultBlock, error) {
	res, err := rpc.Client.BlockByHash(ctx, hash)
	return res, rpcError(rpc.nodeAddress, err)
}

func (rpc rpcWrapper) BlockResults(ctx context.Context, height *int64) (*ctypes.ResultBlockResults, error) {
	res, err := rpc.Client.BlockResults(ctx, height)
	return res, rpcError(rpc.nodeAddress, err)
}

func (rpc rpcWrapper) Commit(ctx context.Context, height *int64) (*ctypes.ResultCommit, error) {
	res, err := rpc.Client.Commit(ctx, height)
	return res, rpcError(rpc.nodeAddress, err)
}

func (rpc rpcWrapper) Validators(ctx context.Context, height *int64, page *int, perPage *int) (*ctypes.ResultValidators, error) {
	res, err := rpc.Client.Validators(ctx, height, page, perPage)
	return res, rpcError(rpc.nodeAddress, err)
}

func (rpc rpcWrapper) Tx(ctx context.Context, hash []byte, prove bool) (*ctypes.ResultTx, error) {
	res, err := rpc.Client.Tx(ctx, hash, prove)
	return res, rpcError(rpc.nodeAddress, err)
}

func (rpc rpcWrapper) TxSearch(ctx context.Context, query string, prove bool, page *int, perPage *int, orderBy string) (*ctypes.ResultTxSearch, error) {
	res, err := rpc.Client.TxSearch(ctx, query, prove, page, perPage, orderBy)
	return res, rpcError(rpc.nodeAddress, err)
}

func (rpc rpcWrapper) BlockSearch(ctx context.Context, query string, page *int, perPage *int, orderBy string) (*ctypes.ResultBlockSearch, error) {
	res, err := rpc.Client.BlockSearch(ctx, query, page, perPage, orderBy)
	return res, rpcError(rpc.nodeAddress, err)
}

func (rpc rpcWrapper) Status(ctx context.Context) (*ctypes.ResultStatus, error) {
	res, err := rpc.Client.Status(ctx)
	return res, rpcError(rpc.nodeAddress, err)
}

func (rpc rpcWrapper) BroadcastEvidence(ctx context.Context, e types.Evidence) (*ctypes.ResultBroadcastEvidence, error) {
	res, err := rpc.Client.BroadcastEvidence(ctx, e)
	return res, rpcError(rpc.nodeAddress, err)
}

func (rpc rpcWrapper) UnconfirmedTxs(ctx context.Context, limit *int) (*ctypes.ResultUnconfirmedTxs, error) {
	res, err := rpc.Client.UnconfirmedTxs(ctx, limit)
	return res, rpcError(rpc.nodeAddress, err)
}

func (rpc rpcWrapper) NumUnconfirmedTxs(ctx context.Context) (*ctypes.ResultUnconfirmedTxs, error) {
	res, err := rpc.Client.NumUnconfirmedTxs(ctx)
	return res, rpcError(rpc.nodeAddress, err)
}

func (rpc rpcWrapper) CheckTx(ctx context.Context, tx types.Tx) (*ctypes.ResultCheckTx, error) {
	res, err := rpc.Client.CheckTx(ctx, tx)
	return res, rpcError(rpc.nodeAddress, err)
}
