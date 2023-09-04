package plugin

import (
	"context"
)

type Chainer interface {
	// AppPath returns the configured App's path.
	AppPath() string
	// ID returns the configured App's chain id.
	ID() (string, error)
	// ConfigPath returns the path to the App's config file.
	ConfigPath() string
	// RPCPublicAddress returns the configured App's rpc endpoint.
	RPCPublicAddress() (string, error)
}

// NewClientAPI creates a new app ClientAPI.
func NewClientAPI(c Chainer) ClientAPI {
	return clientAPI{chain: c}
}

type clientAPI struct {
	chain Chainer
}

func (c clientAPI) GetChainInfo(ctx context.Context) (*ChainInfo, error) {
	chain_id, err := c.chain.ID()
	if err != nil {
		return nil, err
	}
	app_path := c.chain.AppPath()
	config_path := c.chain.ConfigPath()
	rpc, err := c.chain.RPCPublicAddress()
	if err != nil {
		return nil, err
	}
	return &ChainInfo{
		ChainId:    chain_id,
		AppPath:    app_path,
		ConfigPath: config_path,
		RpcAddress: rpc,
	}, nil
}
