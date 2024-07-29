package cosmosclient

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
)

var _ Signer = signer{}

// signer implements the Signer interface.
type signer struct{}

func (signer) Sign(ctx context.Context, txf tx.Factory, name string, txBuilder client.TxBuilder, overwriteSig bool) error {
	return tx.Sign(ctx, txf, name, txBuilder, overwriteSig)
}
