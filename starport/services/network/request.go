package network

import (
	"context"

	launchtypes "github.com/tendermint/spn/x/launch/types"
	"github.com/tendermint/starport/starport/pkg/yaml"
)

// FetchRequest fetches the chain request from SPN by launch and request id
func (b *Builder) FetchRequest(ctx context.Context, launchID, requestID uint64) (string, error) {
	res, err := launchtypes.NewQueryClient(b.cosmos.Context).Request(ctx, &launchtypes.QueryGetRequestRequest{
		LaunchID:  launchID,
		RequestID: requestID,
	})
	if err != nil {
		return "", err
	}

	// convert the request object to yaml
	requestYaml, err := yaml.ParseString(ctx, res.Request,
		"$.content.content.genesisValidator.genTx",
		"$.content.content.genesisValidator.consPubKey",
	)
	if err != nil {
		return "", err
	}

	return requestYaml, err
}
