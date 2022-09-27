package cosmosclient

import (
	"context"
	"errors"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

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

// Broadcast signs and broadcasts this tx.
// If faucet is enabled and if the from account doesn't have enough funds, is
// it automatically filled with the default amount, and the tx is broadcasted
// again. Note that this may still end with the same error if the amount is
// greater than the amount dumped by the faucet.
func (s TxService) Broadcast() (Response, error) {
	defer s.client.lockBech32Prefix()()

	accountName := s.clientContext.GetFromName()
	accountAddress := s.clientContext.GetFromAddress()

	// validate msgs.
	for _, msg := range s.txBuilder.GetTx().GetMsgs() {
		if err := msg.ValidateBasic(); err != nil {
			return Response{}, err
		}
	}

	if err := s.client.signer.Sign(s.txFactory, accountName, s.txBuilder, true); err != nil {
		return Response{}, err
	}

	txBytes, err := s.clientContext.TxConfig.TxEncoder()(s.txBuilder.GetTx())
	if err != nil {
		return Response{}, err
	}

	resp, err := s.clientContext.BroadcastTx(txBytes)
	if s.client.useFaucet && errors.Is(err, sdkerrors.ErrInsufficientFunds) {
		err = s.client.makeSureAccountHasTokens(context.Background(), accountAddress.String())
		if err != nil {
			return Response{}, err
		}
		resp, err = s.clientContext.BroadcastTx(txBytes)
	}

	return Response{
		Codec:      s.clientContext.Codec,
		TxResponse: resp,
	}, handleBroadcastResult(resp, err)
}

// EncodeJSON encodes the transaction as a json string
func (s TxService) EncodeJSON() ([]byte, error) {
	return s.client.context.TxConfig.TxJSONEncoder()(s.txBuilder.GetTx())
}
