package network

import (
	"context"

	launchtypes "github.com/tendermint/spn/x/launch/types"
)

// Requests fetches the chain requests from SPN by launch id
func (n Network) Requests(ctx context.Context, launchID uint64) ([]launchtypes.Request, error) {
	res, err := launchtypes.NewQueryClient(n.cosmos.Context).RequestAll(ctx, &launchtypes.QueryAllRequestRequest{
		LaunchID: launchID,
	})
	if err != nil {
		return nil, err
	}

	return res.Request, err
}
