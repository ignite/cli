package cosmosfaucet

import (
	"encoding/json"
	"fmt"
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
}

type transferResponse struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

func (f Faucet) handler(w http.ResponseWriter, r *http.Request) {
	var req transferRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		xhttp.ResponseJSON(w, http.StatusBadRequest, transferResponse{
			Status: statusError,
			Error:  err.Error(),
		})
		return
	}

	amount := fmt.Sprintf("%d%s", f.creditAmount, f.denom)

	if err := f.Transfer(r.Context(), req.AccountAddress, amount); err != nil {
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

// ServeHTTP implements http.Handler to expose the functionality of Faucet.Transfer() via HTTP.
// request/response payloads are compatible with the previous implementation at allinbits/cosmos-faucet.
func (f Faucet) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		xhttp.ResponseJSON(w, http.StatusMethodNotAllowed, transferResponse{
			Status: statusError,
			Error:  "method not allowed",
		})
		return
	}

	cors.Default().Handler(http.HandlerFunc(f.handler)).ServeHTTP(w, r)
}
