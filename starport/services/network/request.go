package network

import (
	"context"

	launchtypes "github.com/tendermint/spn/x/launch/types"
	"github.com/tendermint/starport/starport/pkg/events"
)

// ChainLaunches returns the list of chain launches in the network
func (b *Builder) ChainLaunches(ctx context.Context) ([]launchtypes.Chain, error) {
	b.ev.Send(events.New(events.StatusOngoing, "Fetching chain launches"))
	res, err := launchtypes.NewQueryClient(b.cosmos.Context).ChainAll(ctx, &launchtypes.QueryAllChainRequest{})
	if err != nil {
		return []launchtypes.Chain{}, err
	}
	return res.Chain, err
}
