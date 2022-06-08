package network

import (
	"context"

	monitoringctypes "github.com/tendermint/spn/x/monitoringc/types"

	"github.com/ignite-hq/cli/ignite/pkg/cosmoserror"
	"github.com/ignite-hq/cli/ignite/services/network/networktypes"
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

// FindClientID find client, connection and channel id by the chain launch id
func (n Network) FindClientID(ctx context.Context, launchID uint64) (relayer networktypes.Relayer, err error) {
	clientStates, err := n.verifiedClientIDs(ctx, launchID)
	if err != nil {
		return relayer, err
	}
	if len(clientStates) == 0 {
		return relayer, ErrObjectNotFound
	}
	relayer.ClientID = clientStates[0]

	connections, err := n.node.clientConnections(ctx, relayer.ClientID)
	if err != nil && err != ErrObjectNotFound {
		return relayer, err
	}
	if err == ErrObjectNotFound || len(connections) == 0 {
		return relayer, nil
	}
	relayer.ConnectionID = connections[0]

	channels, err := n.node.connectionChannels(ctx, relayer.ConnectionID)
	if err != nil && err != ErrObjectNotFound {
		return relayer, err
	}
	if err == ErrObjectNotFound || len(connections) == 0 {
		return relayer, nil
	}
	relayer.ChannelID = channels[0]
	return relayer, nil
}
