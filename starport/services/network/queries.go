package network

import (
	"context"

	launchtypes "github.com/tendermint/spn/x/launch/types"
	"github.com/tendermint/starport/starport/pkg/events"
	"github.com/tendermint/starport/starport/services/network/networktypes"
)

// ChainLaunch fetches the chain launch from Starport Network by launch id.
func (n Network) ChainLaunch(ctx context.Context, id uint64) (networktypes.ChainLaunch, error) {
	n.ev.Send(events.New(events.StatusOngoing, "Fetching chain information"))

	res, err := launchtypes.NewQueryClient(n.cosmos.Context).Chain(ctx, &launchtypes.QueryGetChainRequest{
		LaunchID: id,
	})
	if err != nil {
		return networktypes.ChainLaunch{}, err
	}

	n.ev.Send(events.New(events.StatusOngoing, "Chain information fetched"))

	return networktypes.ParseChainLaunch(res.Chain), nil
}

// ChainLaunches fetches the chain launches from Starport Network
func (n Network) ChainLaunches(ctx context.Context) ([]networktypes.ChainLaunch, error) {
	var chainLaunches []networktypes.ChainLaunch

	n.ev.Send(events.New(events.StatusOngoing, "Fetching chains information"))
	res, err := launchtypes.NewQueryClient(n.cosmos.Context).ChainAll(ctx, &launchtypes.QueryAllChainRequest{})
	if err != nil {
		return chainLaunches, err
	}

	// Parse fetched chains
	for _, chain := range res.Chain {
		chainLaunches = append(chainLaunches, networktypes.ParseChainLaunch(chain))
	}

	return chainLaunches, err
}
