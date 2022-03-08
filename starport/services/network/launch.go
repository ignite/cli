package network

import (
	"context"
	"fmt"
	"time"

	launchtypes "github.com/tendermint/spn/x/launch/types"

	"github.com/tendermint/starport/starport/pkg/events"
	"github.com/tendermint/starport/starport/pkg/xtime"
	"github.com/tendermint/starport/starport/services/network/networktypes"
)

// LaunchParams fetches the chain launch module params from SPN
func (n Network) LaunchParams(ctx context.Context) (launchtypes.Params, error) {
	res, err := launchtypes.NewQueryClient(n.cosmos.Context).Params(ctx, &launchtypes.QueryParamsRequest{})
	if err != nil {
		return launchtypes.Params{}, err
	}
	return res.GetParams(), nil
}

// TriggerLaunch launches a chain as a coordinator
func (n Network) TriggerLaunch(ctx context.Context, launchID uint64, remainingTime time.Duration) error {
	n.ev.Send(events.New(events.StatusOngoing, fmt.Sprintf("Launching chain %d", launchID)))
	params, err := n.LaunchParams(ctx)
	if err != nil {
		return err
	}

	var (
		minLaunch = xtime.Seconds(params.LaunchTimeRange.MinLaunchTime)
		maxLaunch = xtime.Seconds(params.LaunchTimeRange.MaxLaunchTime)
		address   = n.account.Address(networktypes.SPN)
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

	msg := launchtypes.NewMsgTriggerLaunch(address, launchID, int64(remainingTime.Seconds()))
	n.ev.Send(events.New(events.StatusOngoing, "Setting launch time"))
	res, err := n.cosmos.BroadcastTx(n.account.Name, msg)
	if err != nil {
		return err
	}

	var launchRes launchtypes.MsgTriggerLaunchResponse
	if err := res.Decode(&launchRes); err != nil {
		return err
	}

	n.ev.Send(events.New(events.StatusDone,
		fmt.Sprintf("Chain %d will be launched on %s", launchID, xtime.NowAfter(remainingTime)),
	))
	return nil
}
