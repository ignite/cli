package cosmosfaucet

import (
	"encoding/json"
	"net/http"

	"github.com/rs/cors"
	"github.com/tendermint/starport/starport/pkg/xhttp"
)

const (
	statusOK    = "ok"
	statusError = "error"
)

type transferRequest struct {
	AccountAddress string `json:"address"`
	Denom          string `json:"denom"`
}

type transferResponse struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

// ServeHTTP implements http.Handler to expose the functionality of Faucet.Transfer() via HTTP.
// request/response payloads are compatible with the previous implementation at allinbits/cosmos-faucet.
func (f Faucet) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// check method.
	if r.Method != http.MethodPost {
		xhttp.ResponseJSON(w, http.StatusMethodNotAllowed, transferResponse{
			Status: statusError,
			Error:  "method not allowed",
		})
		return
	}

	// add CORS.
	cors.Default().Handler(http.HandlerFunc(f.handler)).ServeHTTP(w, r)
}

func (f Faucet) handler(w http.ResponseWriter, r *http.Request) {
	var req transferRequest

	// decode request into req.
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		xhttp.ResponseJSON(w, http.StatusBadRequest, transferResponse{
			Status: statusError,
			Error:  err.Error(),
		})
		return
	}

	// determine the coin amount and denom to transfer.
	coin := f.coinByDenom(req.Denom)

	if err := f.Transfer(r.Context(), req.AccountAddress, coin.amount, coin.denom); err != nil {
		xhttp.ResponseJSON(w, http.StatusInternalServerError, transferResponse{
			Status: statusError,
			Error:  err.Error(),
		})
		return
	}

	xhttp.ResponseJSON(w, http.StatusOK, transferResponse{
		Status: statusOK,
	})
}

func (f Faucet) coinByDenom(denom string) coin {
	for _, coin := range f.coins {
		if coin.denom == denom {
			return coin
		}
	}

	// otherwise use the default one.
	return f.coins[0]
}
