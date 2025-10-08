package buf

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ignite/cli/v29/ignite/pkg/errors"
)

var BufTokenURL = "https://api.ignite.com/v1/buf" //nolint:gosec // URL is hardcoded and not user-provided

// FetchToken fetches the buf token from the Ignite API.
func FetchToken() (string, error) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(BufTokenURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.Errorf("HTTP request failed with status code: %d", resp.StatusCode)
	}

	type tokenResponse struct {
		Token string `json:"token"`
	}
	var tokenResp tokenResponse

	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", err
	}

	return tokenResp.Token, nil
}
