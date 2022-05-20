package network

import (
	"context"
	"encoding/base64"
	"fmt"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ibcclienttypes "github.com/cosmos/ibc-go/v2/modules/core/02-client/types"
	"github.com/cosmos/ibc-go/v2/modules/core/exported"
	lightclienttypes "github.com/cosmos/ibc-go/v2/modules/light-clients/07-tendermint/types"
	spntypes "github.com/tendermint/spn/pkg/types"

	"github.com/ignite-hq/cli/ignite/services/network/networktypes"
)

// Node is node builder.
type Node struct {
	cosmos         CosmosClient
	stakingQuery   stakingtypes.QueryClient
	ibcClientQuery ibcclienttypes.QueryClient
}

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

// FindClientID fetches the client id from states for a given chain id
func (n Node) FindClientID(ctx context.Context, chainID string) (string, error) {
	states, err := n.states(ctx)
	if err != nil {
		return "", err
	}
	for _, state := range states {
		clientState, ok := state.ClientState.GetCachedValue().(*lightclienttypes.ClientState)
		if !ok && clientState != nil {
			continue
		}
		if clientState.ChainId == chainID {
			return state.ClientId, nil
		}
	}
	return "", fmt.Errorf("client id state not found for chain %s", chainID)
}

// states fetches the chain ibc states
func (n Node) states(ctx context.Context) (ibcclienttypes.IdentifiedClientStates, error) {
	res, err := n.ibcClientQuery.ClientStates(ctx, &ibcclienttypes.QueryClientStatesRequest{})
	if err != nil {
		return ibcclienttypes.IdentifiedClientStates{}, err
	}
	return res.ClientStates, nil
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
