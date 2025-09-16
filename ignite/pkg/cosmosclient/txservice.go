package cosmosclient

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdktypes "github.com/cosmos/cosmos-sdk/types"

	"github.com/ignite/cli/v29/ignite/pkg/errors"
)

// BroadcastOption configures broadcast behavior.
type BroadcastOption func(*broadcastConfig)

// broadcastConfig holds configuration for broadcast operations.
type broadcastConfig struct {
	sequence  uint64
	unordered bool
}

// WithSequence overrides the account sequence for the transaction.
func WithSequence(sequence uint64) BroadcastOption {
	return func(cfg *broadcastConfig) {
		cfg.sequence = sequence
	}
}

// WithUnordered sets the unordered flag for the transaction.
// NOTE: This is only supported for Cosmos SDK versions >= v0.53.
func WithUnordered(unordered bool) BroadcastOption {
	return func(cfg *broadcastConfig) {
		cfg.unordered = unordered
	}
}

type TxService struct {
	client        Client
	clientContext client.Context
	txBuilder     client.TxBuilder
	txFactory     tx.Factory
}

// Gas is gas decided to use for this tx.
// either calculated or configured by the caller.
func (s TxService) Gas() uint64 {
	return s.txBuilder.GetTx().GetGas()
}

// broadcast signs and broadcasts the transaction returning the initial broadcast response.
func (s TxService) broadcast(ctx context.Context, opts ...BroadcastOption) (*sdktypes.TxResponse, error) {
	defer s.client.lockBech32Prefix()()

	// validate msgs.
	for _, msg := range s.txBuilder.GetTx().GetMsgs() {
		msg, ok := msg.(sdktypes.HasValidateBasic)
		if !ok {
			continue
		}
		if err := msg.ValidateBasic(); err != nil {
			return nil, errors.WithStack(err)
		}
	}

	// Apply broadcast options
	cfg := &broadcastConfig{}
	for _, opt := range opts {
		opt(cfg)
	}

	// Override sequence if specified
	if cfg.sequence != 0 {
		s.txFactory = s.txFactory.WithSequence(cfg.sequence)
	}
	s.txFactory = s.txFactory.WithUnordered(cfg.unordered)

	accountName := s.clientContext.FromName
	if err := s.client.signer.Sign(ctx, s.txFactory, accountName, s.txBuilder, true); err != nil {
		return nil, errors.WithStack(err)
	}

	txBytes, err := s.clientContext.TxConfig.TxEncoder()(s.txBuilder.GetTx())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	resp, err := s.clientContext.BroadcastTx(txBytes)
	if err := handleBroadcastResult(resp, err); err != nil {
		return nil, err
	}

	return resp, nil
}

// Broadcast signs and broadcasts this tx.
// If faucet is enabled and if the "from" account doesn't have enough funds, is
// it automatically filled with the default amount, and the tx is broadcasted
// again. Note that this may still end with the same error if the amount is
// greater than the amount dumped by the faucet.
func (s TxService) Broadcast(ctx context.Context, opts ...BroadcastOption) (Response, error) {
	resp, err := s.broadcast(ctx, opts...)
	if err != nil {
		return Response{}, err
	}

	res, err := s.client.WaitForTx(ctx, resp.TxHash)
	if err != nil {
		return Response{}, err
	}
	// NOTE(tb) second and third parameters are omitted:
	// - second parameter represents the tx and should be of type sdktypes.Any,
	// but it is very ugly to decode, not sure if it's worth it (see sdk code
	// x/auth/query.go method makeTxResult)
	// - third parameter represents the timestamp of the tx, which must be
	// fetched from the block itself. So it requires another API call to
	// fetch the block from res.Height, not sure if it's worth it too.
	resp = sdktypes.NewResponseResultTx(res, nil, "")

	return Response{
		Codec:      s.clientContext.Codec,
		TxResponse: resp,
	}, handleBroadcastResult(resp, err)
}

// BroadcastAsync signs and broadcasts this tx.
// It is similar to Broadcast but it does not wait for the transaction to be included in a block.
func (s TxService) BroadcastAsync(ctx context.Context, opts ...BroadcastOption) (Response, error) {
	resp, err := s.broadcast(ctx, opts...)
	if err != nil {
		return Response{}, err
	}

	return Response{
		Codec:      s.clientContext.Codec,
		TxResponse: resp,
	}, handleBroadcastResult(resp, err)
}

// EncodeJSON encodes the transaction as a json string.
func (s TxService) EncodeJSON() ([]byte, error) {
	return s.client.context.TxConfig.TxJSONEncoder()(s.txBuilder.GetTx())
}
