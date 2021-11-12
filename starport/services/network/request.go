package network

import (
	"context"

	launchtypes "github.com/tendermint/spn/x/launch/types"
)

// FetchRequest fetches the chain request from SPN by launch and request id
func (b *Builder) FetchRequest(ctx context.Context, launchID, requestID uint64) (launchtypes.Request, error) {
	res, err := launchtypes.NewQueryClient(b.cosmos.Context).Request(ctx, &launchtypes.QueryGetRequestRequest{
		LaunchID:  launchID,
		RequestID: requestID,
	})
	if err != nil {
		return launchtypes.Request{}, err
	}

	return res.Request, err
}
