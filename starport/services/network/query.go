package network

import (
	"context"

	launchtypes "github.com/tendermint/spn/x/launch/types"
	"github.com/tendermint/starport/starport/pkg/events"
)

// ChainLaunch returns the chain launch data from SPN
func (b Builder) ChainLaunch(ctx context.Context, launchID uint64) (launchtypes.Chain, error) {
	b.ev.Send(events.New(events.StatusOngoing, "Fetching chain launch data"))
	res, err := launchtypes.NewQueryClient(b.cosmos.Context).Chain(ctx, &launchtypes.QueryGetChainRequest{
		LaunchID: launchID,
	})
	if err != nil {
		return launchtypes.Chain{}, err
	}
	return res.Chain, err
}
