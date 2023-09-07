package plugin

import (
	"context"
	"encoding/json"

	chainconfig "github.com/ignite/cli/ignite/config/chain"
	"github.com/ignite/cli/ignite/pkg/cosmosanalysis/chain"
)

type Chainer interface {
	// AppPath returns the configured App's path.
	AppPath() string

	// Config returns the configured App's configuration.
	Config() (*chainconfig.Config, error)

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

func (api clientAPI) GetChainInfo(context.Context) (*ChainInfo, error) {
	chainID, err := api.chain.ID()
	if err != nil {
		return nil, err
	}

	rpc, err := api.chain.RPCPublicAddress()
	if err != nil {
		return nil, err
	}

	return &ChainInfo{
		ChainId:    chainID,
		AppPath:    api.chain.AppPath(),
		ConfigPath: api.chain.ConfigPath(),
		RpcAddress: rpc,
	}, nil
}

func (api clientAPI) GetModuleList(ctx context.Context) (*ModuleList, error) {
	conf, err := api.chain.Config()
	if err != nil {
		return nil, err
	}
	mods, err := chain.GetModuleList(ctx, api.chain.AppPath(), conf.Build.Proto.Path, conf.Build.Proto.ThirdPartyPaths)
	if err != nil {
		return nil, err
	}
	bz, err := json.Marshal(mods)
	if err != nil {
		return nil, err
	}
	var moduleList ModuleList
	if err := json.Unmarshal(bz, &moduleList); err != nil {
		return nil, err
	}
	return &moduleList, nil
}
