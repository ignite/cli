package cosmosclient

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
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

	if err := s.client.prepareBroadcast(context.Background(), accountName, []sdktypes.Msg{}); err != nil {
		return Response{}, err
	}

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
