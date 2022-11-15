package tendermintrpc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

const (
	endpointNetInfo = "/net_info"
	endpointGenesis = "/genesis"
	endpointStatus  = "/status"
)

// Client is a Tendermint RPC client.
type Client struct {
	addr string
}

// New creates a new Tendermint RPC client.
func New(addr string) Client {
	return Client{addr: addr}
}

// NetInfo represents Network Info.
type NetInfo struct {
	ConnectedPeers int
}

func (c Client) url(endpoint string) string {
	return fmt.Sprintf("%s%s", c.addr, endpoint)
}

// GetNetInfo retrieves network info.
func (c Client) GetNetInfo(ctx context.Context) (NetInfo, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.url(endpointNetInfo), nil)
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

// Genesis represents Genesis.
type Genesis struct {
	ChainID string `json:"chain_id"`
}

// GetGenesis retrieves Genesis.
func (c Client) GetGenesis(ctx context.Context) (Genesis, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.url(endpointGenesis), nil)
	if err != nil {
		return Genesis{}, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Genesis{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Genesis{}, fmt.Errorf("%d", resp.StatusCode)
	}

	var out struct {
		Result struct {
			Genesis Genesis `json:"genesis"`
		} `json:"Result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return Genesis{}, err
	}

	return out.Result.Genesis, nil
}

// NodeInfo holds node info.
type NodeInfo struct {
	Network string
}

// Status retrieves node Status.
func (c Client) Status(ctx context.Context) (NodeInfo, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.url(endpointStatus), nil)
	if err != nil {
		return NodeInfo{}, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return NodeInfo{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return NodeInfo{}, fmt.Errorf("%d", resp.StatusCode)
	}

	var (
		info NodeInfo
		b    = &bytes.Buffer{}
		r    = io.TeeReader(resp.Body, b)
	)

	var out struct {
		Result struct {
			NodeInfo NodeInfo `json:"node_info"`
		} `json:"result"`
	}

	if err := json.NewDecoder(r).Decode(&out); err != nil {
		return NodeInfo{}, err
	}

	info = out.Result.NodeInfo

	// some versions have a different response payload.
	if info.Network == "" {
		var out struct {
			Result struct {
				NodeInfo NodeInfo `json:"NodeInfo"`
			} `json:"result"`
		}

		if err := json.NewDecoder(b).Decode(&out); err != nil {
			return NodeInfo{}, err
		}

		info = out.Result.NodeInfo
	}

	return info, nil
}
