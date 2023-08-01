package cosmosclient

import (
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
	gogogrpc "github.com/cosmos/gogoproto/grpc"
)

// gasometer implements the Gasometer interface.
type gasometer struct{}

func (gasometer) CalculateGas(clientCtx gogogrpc.ClientConn, txf tx.Factory, msgs ...sdktypes.Msg) (*txtypes.SimulateResponse, uint64, error) {
	return tx.CalculateGas(clientCtx, txf, msgs...)
}
