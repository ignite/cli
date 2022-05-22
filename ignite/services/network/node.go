package network

import (
	"context"
	"encoding/base64"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ibcclienttypes "github.com/cosmos/ibc-go/v2/modules/core/02-client/types"
	"github.com/cosmos/ibc-go/v2/modules/core/exported"
	lightclienttypes "github.com/cosmos/ibc-go/v2/modules/light-clients/07-tendermint/types"
	"github.com/pkg/errors"
	spntypes "github.com/tendermint/spn/pkg/types"

	"github.com/ignite-hq/cli/ignite/services/network/networktypes"
)

// Node is node builder.
type Node struct {
	cosmos         CosmosClient
	stakingQuery   stakingtypes.QueryClient
	ibcClientQuery ibcclienttypes.QueryClient
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
		cosmos:         cosmos,
		stakingQuery:   stakingtypes.NewQueryClient(cosmos.Context()),
		ibcClientQuery: ibcclienttypes.NewQueryClient(cosmos.Context()),
	}, nil
}

// RewardsInfo Fetches the consensus state with the validator set
func RewardsInfo(ctx context.Context, client CosmosClient, height int64) (networktypes.Reward, error) {
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

func findClientID(ctx context.Context, client CosmosClient, chainID string) (string, error) {
	ibcClientQuery := ibcclienttypes.NewQueryClient(client.Context())
	client.Context().InterfaceRegistry.RegisterImplementations(
		(*exported.ClientState)(nil),
		&lightclienttypes.ClientState{},
	)
	res, err := ibcClientQuery.ClientStates(ctx, &ibcclienttypes.QueryClientStatesRequest{})
	if err != nil {
		return "", err
	}
	for _, state := range res.ClientStates {
		clientState, ok := state.ClientState.GetCachedValue().(*lightclienttypes.ClientState)
		if !ok && clientState != nil {
			continue
		}
		if clientState.ChainId == chainID {
			return state.ClientId, nil
		}
	}
	return "", ErrChainClientNotExist
}

// FindClientID find an IBC client id by chain
func (n Node) FindClientID(ctx context.Context, chainID string) (string, error) {
	return findClientID(ctx, n.cosmos, chainID)
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

	info, err := RewardsInfo(ctx, n.cosmos, lastBlockHeight)
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
