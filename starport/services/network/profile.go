package network

import (
	"fmt"

	profiletypes "github.com/tendermint/spn/x/profile/types"

	"github.com/tendermint/starport/starport/pkg/events"
	"github.com/tendermint/starport/starport/services/network/networktypes"
)

// SetValidatorOperatorAddress associates a Tendermint operator address to a specific validator address on SPN
func (n Network) SetValidatorOperatorAddress(operatorAddress string) error {
	acc := n.account.Address(networktypes.SPN)
	n.ev.Send(events.New(events.StatusOngoing,
		fmt.Sprintf("Adding the operator address %s to validator %s", operatorAddress, acc)))

	// Create and broadcast the transaction
	msg := profiletypes.NewMsgSAddValidatorOperatorAddress(
		n.account.Address(networktypes.SPN),
		operatorAddress,
	)
	res, err := n.cosmos.BroadcastTx(n.account.Name, msg)
	if err != nil {
		return err
	}

	var opAddrRes profiletypes.MsgAddValidatorOperatorAddressResponse
	if err := res.Decode(&opAddrRes); err != nil {
		return err
	}

	n.ev.Send(events.New(events.StatusDone,
		fmt.Sprintf("Validator operator address %s added ", operatorAddress)))
	return nil
}
