package network

import (
	"context"

	monitoringctypes "github.com/tendermint/spn/x/monitoringc/types"

	"github.com/ignite-hq/cli/ignite/services/network/networktypes"
)

func (n Network) CreateClient(
	launchID uint64,
	unbondingTime int64,
	rewardsInfo networktypes.Reward,
) (string, error) {
	msgCreateClient := monitoringctypes.NewMsgCreateClient(
		n.account.Address(networktypes.SPN),
		launchID,
		rewardsInfo.ConsensusState,
		rewardsInfo.ValidatorSet,
		unbondingTime,
		rewardsInfo.RevisionHeight,
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

// FindClientID find an IBC client id by chain
func (n Network) FindClientID(ctx context.Context, chainID string) (string, error) {
	return findClientID(ctx, n.cosmos, chainID)
}
