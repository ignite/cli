package cosmosfaucet

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ignite/cli/ignite/pkg/xhttp"
)

type TransferRequest struct {
	// AccountAddress to request for coins.
	AccountAddress string `json:"address"`

	// Coins that are requested.
	// default ones used when this one isn't provided.
	Coins []string `json:"coins"`
}

func NewTransferRequest(accountAddress string, coins []string) TransferRequest {
	return TransferRequest{
		AccountAddress: accountAddress,
		Coins:          coins,
	}
}

type TransferResponse struct {
	Error string `json:"error,omitempty"`
}

func (f Faucet) faucetHandler(w http.ResponseWriter, r *http.Request) {
	var req TransferRequest

	// decode request into req.
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responseError(w, http.StatusBadRequest, err)
		return
	}

	// determine coins to transfer.
	coins, err := f.coinsFromRequest(req)
	if err != nil {
		responseError(w, http.StatusBadRequest, err)
		return
	}

	// try performing the transfer
	if err := f.Transfer(r.Context(), req.AccountAddress, coins); err != nil {
		if errors.Is(err, context.Canceled) {
			return
		}
		responseError(w, http.StatusInternalServerError, err)
	} else {
		responseSuccess(w)
	}
}

// FaucetInfoResponse is the faucet info payload.
type FaucetInfoResponse struct {
	// IsAFaucet indicates that this is a faucet endpoint.
	// useful for auto discoveries.
	IsAFaucet bool `json:"is_a_faucet"`

	// ChainID is chain id of the chain that faucet is running for.
	ChainID string `json:"chain_id"`
}

func (f Faucet) faucetInfoHandler(w http.ResponseWriter, _ *http.Request) {
	xhttp.ResponseJSON(w, http.StatusOK, FaucetInfoResponse{
		IsAFaucet: true,
		ChainID:   f.chainID,
	})
}

// coinsFromRequest determines tokens to transfer from transfer request.
func (f Faucet) coinsFromRequest(req TransferRequest) (sdk.Coins, error) {
	if len(req.Coins) == 0 {
		return f.coins, nil
	}

	var coins []sdk.Coin
	for _, c := range req.Coins {
		coin, err := sdk.ParseCoinNormalized(c)
		if err != nil {
			return nil, err
		}
		coins = append(coins, coin)
	}

	return coins, nil
}

func responseSuccess(w http.ResponseWriter) {
	xhttp.ResponseJSON(w, http.StatusOK, TransferResponse{})
}

func responseError(w http.ResponseWriter, code int, err error) {
	xhttp.ResponseJSON(w, code, TransferResponse{
		Error: err.Error(),
	})
}
