package network

import (
	spntypes "github.com/tendermint/spn/pkg/types"
	monitoringctypes "github.com/tendermint/spn/x/monitoringc/types"

	"github.com/ignite-hq/cli/ignite/services/network/networktypes"
)

func (n Network) CreateClient(
	launchID uint64,
	consensusState spntypes.ConsensusState,
	validatorSet spntypes.ValidatorSet,
	unbondingTime int64,
	height uint64,
) (string, error) {
	msgCreateClient := monitoringctypes.NewMsgCreateClient(
		n.account.Address(networktypes.SPN),
		launchID,
		consensusState,
		validatorSet,
		unbondingTime,
		height,
	)

	res, err := n.cosmos.BroadcastTx(n.account.Name, msgCreateClient)
	if err != nil {
		return "", err
	}

	var createClientRes monitoringctypes.MsgCreateClientResponse
	if err := res.Decode(&createClientRes); err != nil {
		return "", err
	}
	return createClientRes.ClientID, nil
}
