package cosmosfaucet

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/tendermint/starport/starport/pkg/cosmoscoin"
	"github.com/tendermint/starport/starport/pkg/xhttp"
)

type TransferRequest struct {
	// AccountAddress to request for coins.
	AccountAddress string `json:"address"`

	// Coins that are requested.
	// default ones used when this one isn't provided.
	Coins []string `json:"coins"`
}

type TransferResponse struct {
	Error string            `json:"error,omitempty"`
	Coins []cosmoscoin.Coin `json:"coins"`
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
		if err == context.Canceled {
			return
		}
		responseError(w, http.StatusInternalServerError, err)
	} else {
		responseSuccess(w, coins)
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

func (f Faucet) faucetInfoHandler(w http.ResponseWriter, r *http.Request) {
	xhttp.ResponseJSON(w, http.StatusOK, FaucetInfoResponse{
		IsAFaucet: true,
		ChainID:   f.chainID,
	})
}

// coinsFromRequest determines tokens to transfer from transfer request.
func (f Faucet) coinsFromRequest(req TransferRequest) ([]cosmoscoin.Coin, error) {
	if len(req.Coins) == 0 {
		return f.coins, nil
	}

	var coins []cosmoscoin.Coin
	for _, c := range req.Coins {
		coin, err := cosmoscoin.Parse(c)
		if err != nil {
			return nil, err
		}
		coins = append(coins, coin)
	}

	return coins, nil
}

func responseSuccess(w http.ResponseWriter, coins []cosmoscoin.Coin) {
	xhttp.ResponseJSON(w, http.StatusOK, TransferResponse{
		Coins: coins,
	})
}

func responseError(w http.ResponseWriter, code int, err error) {
	xhttp.ResponseJSON(w, code, TransferResponse{
		Error: err.Error(),
	})
}
