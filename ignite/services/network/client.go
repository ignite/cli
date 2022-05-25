package network

import (
	"context"

	"github.com/pkg/errors"
	monitoringctypes "github.com/tendermint/spn/x/monitoringc/types"

	"github.com/ignite-hq/cli/ignite/pkg/cosmoserror"
	"github.com/ignite-hq/cli/ignite/services/network/networktypes"
)

// ErrChainClientNotExist returned when specified client does not exist for a chain.
var ErrChainClientNotExist = errors.New("client id not found for chain")

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

// VerifiedClientIDs fetches the verified client ids from SPN by launch id
func (n Network) VerifiedClientIDs(ctx context.Context, launchID uint64) ([]string, error) {
	res, err := n.monitoringcQuery.
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
	clientStates, err := n.VerifiedClientIDs(ctx, launchID)
	if err != nil {
		return relayer, err
	}
	if len(clientStates) == 0 {
		return relayer, ErrChainClientNotExist
	}
	clientID := clientStates[0]
	relayer.ConnectionID, relayer.Channel, err = n.GetConnectionChannel(ctx, clientID)
	return relayer, err
}
