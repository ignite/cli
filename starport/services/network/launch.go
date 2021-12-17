package network

import (
	"context"
	"fmt"
	"time"

	launchtypes "github.com/tendermint/spn/x/launch/types"
	"github.com/tendermint/starport/starport/pkg/cosmoserror"
	"github.com/tendermint/starport/starport/pkg/events"
	"github.com/tendermint/starport/starport/pkg/xtime"
	"github.com/tendermint/starport/starport/services/network/networkchain"
)

// LaunchParams fetches the chain launch module params from SPN
func (n Network) LaunchParams(ctx context.Context) (launchtypes.Params, error) {
	res, err := launchtypes.NewQueryClient(n.cosmos.Context).Params(ctx, &launchtypes.QueryParamsRequest{})
	if err != nil {
		return launchtypes.Params{}, cosmoserror.Unwrap(err)
	}
	return res.GetParams(), nil
}

// TriggerLaunch launches a chain as a coordinator
func (n Network) TriggerLaunch(ctx context.Context, launchID uint64, remainingTime time.Duration) error {
	n.ev.Send(events.New(events.StatusOngoing, fmt.Sprintf("Launching chain %d", launchID)))
	params, err := n.LaunchParams(ctx)
	if err != nil {
		return cosmoserror.Unwrap(err)
	}

	var (
		minLaunch = xtime.Seconds(params.MinLaunchTime)
		maxLaunch = xtime.Seconds(params.MaxLaunchTime)
		address   = n.account.Address(networkchain.SPN)
	)
	switch {
	case remainingTime == 0:
		// if the user does not specify the remaining time, use the minimal one
		remainingTime = minLaunch
	case remainingTime < minLaunch:
		return fmt.Errorf("remaining time %s lower than minimum %s",
			xtime.NowAfter(remainingTime),
			xtime.NowAfter(minLaunch))
	case remainingTime > maxLaunch:
		return fmt.Errorf("remaining time %s greater than maximum %s",
			xtime.NowAfter(remainingTime),
			xtime.NowAfter(maxLaunch))
	}

	msg := launchtypes.NewMsgTriggerLaunch(address, launchID, uint64(remainingTime.Seconds()))
	n.ev.Send(events.New(events.StatusOngoing, "Setting launch time"))
	res, err := n.cosmos.BroadcastTx(n.account.Name, msg)
	if err != nil {
		return cosmoserror.Unwrap(err)
	}

	var launchRes launchtypes.MsgTriggerLaunchResponse
	if err := res.Decode(&launchRes); err != nil {
		return cosmoserror.Unwrap(err)
	}

	n.ev.Send(events.New(events.StatusDone,
		fmt.Sprintf("Chain %d will be launched on %s", launchID, xtime.NowAfter(remainingTime)),
	))
	return nil
}
