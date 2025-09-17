package buf

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var BufTokenURL = "https://api.ignite.com/v1/buf/token"

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
		return "", fmt.Errorf("HTTP request failed with status code: %d", resp.StatusCode)
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
