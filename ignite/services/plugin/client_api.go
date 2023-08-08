package plugin

import (
	"context"
	"encoding/json"

	"github.com/ignite/cli/ignite/pkg/cosmosanalysis/chain"
	chainservice "github.com/ignite/cli/ignite/services/chain"
)

// NewClientAPI creates a new app ClientAPI.
func NewClientAPI(c *chainservice.Chain) ClientAPI {
	return clientAPI{chain: c}
}

type clientAPI struct {
	chain *chainservice.Chain
}

// TODO: Implement dependency ClientAPI.

// Deoendencies returns chain app dependencies.
func (c clientAPI) Dependencies(ctx context.Context) (*Dependencies, error) {
	conf, err := c.chain.Config()
	if err != nil {
		return nil, err
	}
	mods, err := chain.GetModuleList(ctx, c.chain.AppPath(), conf.Build.Proto.Path, conf.Build.Proto.ThirdPartyPaths)
	if err != nil {
		return nil, err
	}
	bz, err := json.Marshal(mods)
	if err != nil {
		return nil, err
	}
	var d Dependencies
	if err := json.Unmarshal(bz, &d); err != nil {
		return nil, err
	}
	return &d, nil
}
