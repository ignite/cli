package plugin

import (
	"context"
	"encoding/json"

	"github.com/ignite/cli/ignite/pkg/cosmosanalysis/chain"
	chainservice "github.com/ignite/cli/ignite/services/chain"
)

// NewClientAPI creates a new app ClientAPI.
func NewClientAPI(c *chainservice.Chain) clientAPI {
	return clientAPI{chain: c}
}

type clientAPI struct {
	chain *chainservice.Chain
}

// TODO: Implement dependency ClientAPI.

// Deoendencies returns chain app dependencies.
func (c clientAPI) Dependencies(ctx context.Context) (*Dependencies, error) {

	mods, err := chain.GetModuleList(ctx, c.chain)

	if err != nil {
		return nil, err
	}
	ret := &Dependencies{}
	bytes, err := json.Marshal(mods)

	if err != nil {
		return nil, err
	}
	json.Unmarshal(bytes, ret)
	return ret, nil
}
