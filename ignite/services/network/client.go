package network

import (
	"context"
	"errors"

	monitoringctypes "github.com/tendermint/spn/x/monitoringc/types"

	"github.com/ignite/cli/ignite/pkg/cosmoserror"
	"github.com/ignite/cli/ignite/services/network/networktypes"
)

// CreateClient send create client message to SPN
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

	res, err := n.cosmos.BroadcastTx(n.account, msgCreateClient)
	if err != nil {
		return "", err
	}

	var createClientRes monitoringctypes.MsgCreateClientResponse
	if err := res.Decode(&createClientRes); err != nil {
		return "", err
	}
	return createClientRes.ClientID, nil
}

// verifiedClientIDs fetches the verified client ids from SPN by launch id
func (n Network) verifiedClientIDs(ctx context.Context, launchID uint64) ([]string, error) {
	res, err := n.monitoringConsumerQuery.
		VerifiedClientIds(ctx,
			&monitoringctypes.QueryGetVerifiedClientIdsRequest{
				LaunchID: launchID,
			},
		)

	if cosmoserror.Unwrap(err) == cosmoserror.ErrNotFound {
		return nil, ErrObjectNotFound
	} else if err != nil {
		return nil, err
	}
	return res.ClientIds, nil
}

// RewardIBCInfo returns IBC info to relay packets for a chain to claim rewards.
func (n Network) RewardIBCInfo(ctx context.Context, launchID uint64) (networktypes.RewardIBCInfo, error) {
	clientStates, err := n.verifiedClientIDs(ctx, launchID)
	if err != nil {
		return networktypes.RewardIBCInfo{}, err
	}
	if len(clientStates) == 0 {
		return networktypes.RewardIBCInfo{}, ErrObjectNotFound
	}

	clientID := clientStates[0]

	connections, err := n.node.clientConnections(ctx, clientID)
	if err != nil && !errors.Is(err, ErrObjectNotFound) {
		return networktypes.RewardIBCInfo{}, err
	}
	if errors.Is(err, ErrObjectNotFound) || len(connections) == 0 {
		return networktypes.RewardIBCInfo{}, nil
	}

	connectionID := connections[0]

	channels, err := n.node.connectionChannels(ctx, connectionID)
	if err != nil && !errors.Is(err, ErrObjectNotFound) {
		return networktypes.RewardIBCInfo{}, err
	}
	if errors.Is(err, ErrObjectNotFound) || len(connections) == 0 {
		return networktypes.RewardIBCInfo{}, nil
	}

	info := networktypes.RewardIBCInfo{
		ClientID:     clientID,
		ConnectionID: connectionID,
		ChannelID:    channels[0],
	}

	return info, nil
}
