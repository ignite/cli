package cosmosclient

import (
	gogogrpc "github.com/cosmos/gogoproto/grpc"

	"cosmossdk.io/core/transaction"

	"github.com/cosmos/cosmos-sdk/client/tx"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
)

// gasometer implements the Gasometer interface.
type gasometer struct{}

func (gasometer) CalculateGas(clientCtx gogogrpc.ClientConn, txf tx.Factory, msgs ...transaction.Msg) (*txtypes.SimulateResponse, uint64, error) {
	return tx.CalculateGas(clientCtx, txf, msgs...)
}
