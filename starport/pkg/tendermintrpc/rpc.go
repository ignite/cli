package tendermintrpc

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// Client is a Tendermint RPC client.
type Client struct {
	addr string
}

// New creates a new Tendermint RPC client.
func New(addr string) *Client {
	return &Client{addr: addr}
}

// NetInfo represents Network Info.
type NetInfo struct {
	ConnectedPeers int
}

// GetNetInfo retrieves network info.
func (c *Client) GetNetInfo(ctx context.Context) (NetInfo, error) {
	endpoint := fmt.Sprintf("%s/net_info", c.addr)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return NetInfo{}, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return NetInfo{}, err
	}
	defer resp.Body.Close()

	var res struct {
		Result struct {
			Peers string `json:"n_peers"`
		} `json:"result"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return NetInfo{}, err
	}

	peers, err := strconv.ParseUint(res.Result.Peers, 10, 64)
	if err != nil {
		return NetInfo{}, err
	}

	return NetInfo{
		ConnectedPeers: int(peers),
	}, nil
}
