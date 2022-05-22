package network

import (
	"context"
	"encoding/base64"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ibcclienttypes "github.com/cosmos/ibc-go/v2/modules/core/02-client/types"
	ibcconntypes "github.com/cosmos/ibc-go/v2/modules/core/03-connection/types"
	ibcchanneltypes "github.com/cosmos/ibc-go/v2/modules/core/04-channel/types"
	"github.com/cosmos/ibc-go/v2/modules/core/exported"
	lightclienttypes "github.com/cosmos/ibc-go/v2/modules/light-clients/07-tendermint/types"
	"github.com/pkg/errors"
	spntypes "github.com/tendermint/spn/pkg/types"

	"github.com/ignite-hq/cli/ignite/pkg/cosmoserror"
	"github.com/ignite-hq/cli/ignite/services/network/networktypes"
)

// Node is node builder.
type Node struct {
	cosmos          CosmosClient
	stakingQuery    stakingtypes.QueryClient
	ibcClientQuery  ibcclienttypes.QueryClient
	ibcConnQuery    ibcconntypes.QueryClient
	ibcChannelQuery ibcchanneltypes.QueryClient
}

// ErrChainClientNotExist returned when specified client does not exist for a chain.
var ErrChainClientNotExist = errors.New("client id not found for chain")

// NewNodeClient creates a new client for node API
func NewNodeClient(cosmos CosmosClient) (Node, error) {
	cosmos.Context().InterfaceRegistry.RegisterImplementations(
		(*exported.ClientState)(nil),
		&lightclienttypes.ClientState{},
	)
	return Node{
		cosmos:          cosmos,
		stakingQuery:    stakingtypes.NewQueryClient(cosmos.Context()),
		ibcClientQuery:  ibcclienttypes.NewQueryClient(cosmos.Context()),
		ibcConnQuery:    ibcconntypes.NewQueryClient(cosmos.Context()),
		ibcChannelQuery: ibcchanneltypes.NewQueryClient(cosmos.Context()),
	}, nil
}

// consensus Fetches the consensus state with the validator set
func (n Node) consensus(ctx context.Context, client CosmosClient, height int64) (networktypes.Reward, error) {
	consensusState, err := client.ConsensusInfo(ctx, height)
	if err != nil {
		return networktypes.Reward{}, err
	}
	spnConsensusStatue := spntypes.NewConsensusState(
		consensusState.Timestamp,
		consensusState.NextValidatorsHash,
		consensusState.Root,
	)

	validators := make([]spntypes.Validator, len(consensusState.ValidatorSet.Validators))
	for i, validator := range consensusState.ValidatorSet.Validators {
		validators[i] = spntypes.NewValidator(
			base64.StdEncoding.EncodeToString(validator.PubKey.GetEd25519()),
			validator.ProposerPriority,
			validator.VotingPower,
		)
	}

	return networktypes.Reward{
		ConsensusState: spnConsensusStatue,
		ValidatorSet:   spntypes.NewValidatorSet(validators...),
		RevisionHeight: uint64(height),
	}, nil
}

// FindClientID Fetches the client id by channel id and check open connections
func (n Node) FindClientID(ctx context.Context, chainID string) (relayer networktypes.Relayer, err error) {
	relayer.ChainID = chainID

	clientStates, err := n.clientStates(ctx)
	for _, state := range clientStates {
		if state.ChainID == chainID {
			relayer.ClientID = state.ClientID

			conns, err := n.clientConnections(ctx, state.ClientID)
			if err != nil {
				return relayer, err
			} else if len(conns) == 0 {
				return relayer, nil
			}
			relayer.ConnectionID = conns[0]

			channels, err := n.connectionChannels(ctx, conns[0])
			if err != nil {
				return relayer, err
			} else if len(channels) == 0 {
				return relayer, nil
			}
			relayer.Channel = channels[0]

			return relayer, nil
		}
	}
	return relayer, ErrChainClientNotExist
}

// connectionChannels fetches the chain connection channels by connection id
func (n Node) connectionChannels(ctx context.Context, connectionID string) (channels []networktypes.Channel, err error) {
	res, err := n.ibcChannelQuery.ConnectionChannels(ctx, &ibcchanneltypes.QueryConnectionChannelsRequest{
		Connection: connectionID,
	})
	if cosmoserror.Unwrap(err) == cosmoserror.ErrNotFound {
		return channels, nil
	} else if err != nil {
		return nil, err
	}
	for _, channel := range res.Channels {
		channels = append(channels, networktypes.ToChannel(*channel))
	}
	return
}

// clientConnections fetches the chain client connections by client id
func (n Node) clientConnections(ctx context.Context, clientID string) ([]string, error) {
	res, err := n.ibcConnQuery.ClientConnections(ctx, &ibcconntypes.QueryClientConnectionsRequest{
		ClientId: clientID,
	})
	if cosmoserror.Unwrap(err) == cosmoserror.ErrNotFound {
		return []string{}, nil
	} else if err != nil {
		return nil, err
	}
	return res.ConnectionPaths, err
}

// clientStates fetches the chain client states
func (n Node) clientStates(ctx context.Context) (states []networktypes.ClientState, err error) {
	res, err := n.ibcClientQuery.ClientStates(ctx, &ibcclienttypes.QueryClientStatesRequest{})
	if err != nil {
		return nil, err
	}
	for _, state := range res.ClientStates {
		clientState, ok := networktypes.ToClientState(state)
		if !ok {
			continue
		}
		states = append(states, clientState)
	}
	return
}

// stakingParams fetches the staking module params
func (n Node) stakingParams(ctx context.Context) (stakingtypes.Params, error) {
	res, err := n.stakingQuery.Params(ctx, &stakingtypes.QueryParamsRequest{})
	if err != nil {
		return stakingtypes.Params{}, err
	}
	return res.Params, nil
}

// RewardsInfo Fetches the consensus state with the validator set and the unbounding time
func (n Node) RewardsInfo(ctx context.Context) (networktypes.Reward, string, int64, error) {
	status, err := n.cosmos.Status(ctx)
	if err != nil {
		return networktypes.Reward{}, "", 0, err
	}
	lastBlockHeight := status.SyncInfo.LatestBlockHeight

	info, err := n.consensus(ctx, n.cosmos, lastBlockHeight)
	if err != nil {
		return networktypes.Reward{}, "", 0, err
	}

	stakingParams, err := n.stakingParams(ctx)
	if err != nil {
		return networktypes.Reward{}, "", 0, err
	}
	return info,
		status.NodeInfo.Network,
		int64(stakingParams.UnbondingTime.Seconds()),
		nil
}
