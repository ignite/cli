package network

import (
	"context"
	"fmt"

	launchtypes "github.com/tendermint/spn/x/launch/types"
	"github.com/tendermint/starport/starport/pkg/cosmosutil"
	"github.com/tendermint/starport/starport/pkg/events"
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

// Launch launch a chain as a coordinator
func (n Network) Launch(ctx context.Context, launchID, remainingTime uint64) error {
	address := n.account.Address(networkchain.SPN)
	spnAddress, err := cosmosutil.ChangeAddressPrefix(address, networkchain.SPN)
	if err != nil {
		return err
	}

	n.ev.Send(events.New(events.StatusOngoing, fmt.Sprintf("Launching chain %d", launchID)))

	params, err := n.LaunchParams(ctx)
	if err != nil {
		return err
	}

	if remainingTime < params.MinLaunchTime {
		remainingTime = params.MinLaunchTime
	} else if remainingTime > params.MaxLaunchTime {
		remainingTime = params.MaxLaunchTime
	}

	msg := launchtypes.NewMsgTriggerLaunch(spnAddress, launchID, remainingTime)
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
		fmt.Sprintf("Chain %d will be launched in %d seconds", launchID, remainingTime),
	))
	return nil
}
