package network

import (
	"fmt"

	launchtypes "github.com/tendermint/spn/x/launch/types"
	"github.com/tendermint/starport/starport/pkg/cosmosutil"
	"github.com/tendermint/starport/starport/pkg/events"
	"github.com/tendermint/starport/starport/services/network/networkchain"
)

// Launch launch a chain as a coordinator
func (n Network) Launch(launchID uint64) error {
	address := n.account.Address(networkchain.SPN)
	spnAddress, err := cosmosutil.ChangeAddressPrefix(address, networkchain.SPN)
	if err != nil {
		return err
	}

	n.ev.Send(events.New(events.StatusOngoing, fmt.Sprintf("Launching chain %d", launchID)))

	msg := launchtypes.NewMsgTriggerLaunch(spnAddress, launchID, 0)
	n.ev.Send(events.New(events.StatusOngoing, "Broadcasting launch transaction"))
	res, err := n.cosmos.BroadcastTx(n.account.Name, msg)
	if err != nil {
		return err
	}

	var requestRes launchtypes.MsgRequestAddAccountResponse
	if err := res.Decode(&requestRes); err != nil {
		return err
	}

	n.ev.Send(events.New(events.StatusDone, "The chain was launched"))
	return nil
}
