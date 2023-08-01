package cosmosclient

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
)

// signer implements the Signer interface.
type signer struct{}

func (signer) Sign(txf tx.Factory, name string, txBuilder client.TxBuilder, overwriteSig bool) error {
	return tx.Sign(txf, name, txBuilder, overwriteSig)
}
