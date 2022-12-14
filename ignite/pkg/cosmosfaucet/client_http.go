package cosmosfaucet

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

// ErrTransferRequest is an error that occurs when a transfer request fails.
type ErrTransferRequest struct {
	StatusCode int
}

// Error implements error.
func (err ErrTransferRequest) Error() string {
	return http.StatusText(err.StatusCode)
}

// HTTPClient is a faucet client.
type HTTPClient struct {
	addr string
}

// NewClient returns a new faucet client.
func NewClient(addr string) HTTPClient {
	return HTTPClient{addr}
}

// Transfer requests tokens from the faucet with req.
func (c HTTPClient) Transfer(ctx context.Context, req TransferRequest) (TransferResponse, error) {
	data, err := json.Marshal(req)
	if err != nil {
		return TransferResponse{}, err
	}

	hreq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.addr, bytes.NewReader(data))
	if err != nil {
		return TransferResponse{}, err
	}

	hres, err := http.DefaultClient.Do(hreq)
	if err != nil {
		return TransferResponse{}, err
	}
	defer hres.Body.Close()

	if hres.StatusCode != http.StatusOK {
		return TransferResponse{}, ErrTransferRequest{hres.StatusCode}
	}

	var res TransferResponse
	err = json.NewDecoder(hres.Body).Decode(&res)
	return res, err
}

// FaucetInfo fetch the faucet info for clients to determine if this is a real faucet and
// what is the chain id of the chain that faucet is operating for.
func (c HTTPClient) FaucetInfo(ctx context.Context) (FaucetInfoResponse, error) {
	hreq, err := http.NewRequestWithContext(ctx, http.MethodGet, c.addr+"/info", nil)
	if err != nil {
		return FaucetInfoResponse{}, err
	}

	hres, err := http.DefaultClient.Do(hreq)
	if err != nil {
		return FaucetInfoResponse{}, err
	}
	defer hres.Body.Close()

	if hres.StatusCode != http.StatusOK {
		return FaucetInfoResponse{}, errors.New(http.StatusText(hres.StatusCode))
	}

	var res FaucetInfoResponse
	err = json.NewDecoder(hres.Body).Decode(&res)
	return res, err
}
