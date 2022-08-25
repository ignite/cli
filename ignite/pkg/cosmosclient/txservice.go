package cosmosclient

import (
	"context"

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
func (s TxService) Broadcast() (Response, error) {
	defer s.client.lockBech32Prefix()()

	accountName := s.clientContext.GetFromName()
	accountAddress := s.clientContext.GetFromAddress()

	// TODO uncomment after https://github.com/tendermint/spn/issues/363
	// validate msgs.
	//  for _, msg := range msgs {
	//  if err := msg.ValidateBasic(); err != nil {
	//  return err
	//  }
	//  }

	if err := tx.Sign(s.txFactory, accountName, s.txBuilder, true); err != nil {
		return Response{}, err
	}

	txBytes, err := s.clientContext.TxConfig.TxEncoder()(s.txBuilder.GetTx())
	if err != nil {
		return Response{}, err
	}

	resp, err := s.clientContext.BroadcastTx(txBytes)
	if err == sdkerrors.ErrInsufficientFunds {
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
