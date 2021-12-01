package network

import (
	"context"

	launchtypes "github.com/tendermint/spn/x/launch/types"
	"github.com/tendermint/starport/starport/pkg/events"
	"github.com/tendermint/starport/starport/services/network/networkchain"
)

type LaunchInfo = networkchain.Launch

// LaunchInfo fetches the chain launch from Starport Network by launch id.
func (n Network) LaunchInfo(ctx context.Context, id uint64) (LaunchInfo, error) {
	n.ev.Send(events.New(events.StatusOngoing, "Fetching chain information"))

	res, err := launchtypes.NewQueryClient(n.cosmos.Context).Chain(ctx, &launchtypes.QueryGetChainRequest{
		LaunchID: id,
	})
	if err != nil {
		return LaunchInfo{}, err
	}

	n.ev.Send(events.New(events.StatusOngoing, "Chain information fetched"))

	return networkchain.ParseLaunch(res.Chain), nil
}

// LaunchesInfo fetches the chain launches from Starport Network
func (n Network) LaunchesInfo(ctx context.Context) ([]LaunchInfo, error) {
	var launchesInfo []LaunchInfo

	n.ev.Send(events.New(events.StatusOngoing, "Fetching chains information"))
	res, err := launchtypes.NewQueryClient(n.cosmos.Context).ChainAll(ctx, &launchtypes.QueryAllChainRequest{})
	if err != nil {
		return launchesInfo, err
	}

	// Parse fetched chains
	for _, chain := range res.Chain {
		launchesInfo = append(launchesInfo, networkchain.ParseLaunch(chain))
	}

	return launchesInfo, err
}
