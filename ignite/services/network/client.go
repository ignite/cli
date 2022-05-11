package network

import (
	monitoringctypes "github.com/tendermint/spn/x/monitoringc/types"

	"github.com/ignite-hq/cli/ignite/services/network/networktypes"
)

func (n Network) CreateClient(
	launchID uint64,
	unbondingTime int64,
	ibcInfo networktypes.IBCInfo,
) (string, error) {
	msgCreateClient := monitoringctypes.NewMsgCreateClient(
		n.account.Address(networktypes.SPN),
		launchID,
		ibcInfo.ConsensusState,
		ibcInfo.ValidatorSet,
		unbondingTime,
		ibcInfo.RevisionHeight,
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
