package network

import (
	"context"
	"fmt"
	"time"

	launchtypes "github.com/tendermint/spn/x/launch/types"
	"github.com/tendermint/starport/starport/pkg/events"
	"github.com/tendermint/starport/starport/pkg/xtime"
	"github.com/tendermint/starport/starport/services/network/networkchain"
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
	remainingTimestamp := uint64(remainingTime.Seconds())
	n.ev.Send(events.New(events.StatusOngoing, fmt.Sprintf("Launching chain %d", launchID)))

	address := n.account.Address(networkchain.SPN)
	params, err := n.LaunchParams(ctx)
	if err != nil {
		return err
	}

	switch {
	case remainingTimestamp == 0:
		// if the user does not specify the remaining time, use the minimal one
		remainingTimestamp = params.MinLaunchTime
	case remainingTimestamp < params.MinLaunchTime:
		return fmt.Errorf("remaining time %s lower than minimum %s",
			xtime.NowAfter(remainingTimestamp),
			xtime.NowAfter(params.MinLaunchTime))
	case remainingTimestamp > params.MaxLaunchTime:
		return fmt.Errorf("remaining time %s greater than maximum %s",
			xtime.NowAfter(remainingTimestamp),
			xtime.NowAfter(params.MaxLaunchTime))
	}

	msg := launchtypes.NewMsgTriggerLaunch(address, launchID, remainingTimestamp)
	n.ev.Send(events.New(events.StatusOngoing, "Broadcasting launch transaction"))
	res, err := n.cosmos.BroadcastTx(n.account.Name, msg)
	if err != nil {
		return err
	}

	var launchRes launchtypes.MsgTriggerLaunchResponse
	if err := res.Decode(&launchRes); err != nil {
		return err
	}

	n.ev.Send(events.New(events.StatusDone,
		fmt.Sprintf("Chain %d will be launched on %s", launchID, xtime.NowAfter(remainingTimestamp)),
	))
	return nil
}
