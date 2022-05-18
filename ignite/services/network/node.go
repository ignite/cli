package network

import (
	"context"
	"encoding/base64"
	"fmt"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ibctypes "github.com/cosmos/ibc-go/v2/modules/core/02-client/types"
	ibctmtypes "github.com/cosmos/ibc-go/v2/modules/light-clients/07-tendermint/types"
	spntypes "github.com/tendermint/spn/pkg/types"
	"github.com/tendermint/spn/x/monitoringp/types"
)

// Node is node builder.
type Node struct {
	cosmos         CosmosClient
	stakingQuery   stakingtypes.QueryClient
	ibcClientQuery ibctypes.QueryClient
}

// IBCInfo is node client info.
type IBCInfo struct {
	ConsensusState spntypes.ConsensusState
	ValidatorSet   spntypes.ValidatorSet
	UnbondingTime  int64
	Height         uint64
}

// NewNodeClient creates a new client for node API
func NewNodeClient(cosmos CosmosClient) (Node, error) {
	return Node{
		cosmos:         cosmos,
		stakingQuery:   stakingtypes.NewQueryClient(cosmos.Context()),
		ibcClientQuery: ibctypes.NewQueryClient(cosmos.Context()),
	}, nil
}

// IBCInfo Fetches the consensus state, validator set and the staking parameters
func (n Node) IBCInfo(ctx context.Context) (IBCInfo, error) {
	status, err := n.cosmos.Status(ctx)
	if err != nil {
		return IBCInfo{}, err
	}
	lastBlockHeight := status.SyncInfo.LatestBlockHeight

	consensusState, err := n.cosmos.ConsensusInfo(ctx, lastBlockHeight)
	if err != nil {
		return IBCInfo{}, err
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

	stakingParams, err := n.stakingParams(ctx)
	if err != nil {
		return IBCInfo{}, err
	}

	return IBCInfo{
		ConsensusState: spnConsensusStatue,
		ValidatorSet:   spntypes.NewValidatorSet(validators...),
		UnbondingTime:  int64(stakingParams.UnbondingTime.Seconds()),
		Height:         uint64(lastBlockHeight),
	}, nil
}

// FindClientID fetches the client id from states for a given chain id
func (n Node) FindClientID(ctx context.Context, chainID string) (string, error) {
	states, err := n.states(ctx)
	if err != nil {
		return "", err
	}
	for _, state := range states {
		tmClientState, ok := state.ClientState.GetCachedValue().(*ibctmtypes.ClientState)
		if !ok {
			return "", types.ErrInvalidClient
		}
		if tmClientState.ChainId == chainID {
			return state.ClientId, nil
		}
	}
	return "", fmt.Errorf("client id state not found for chain %s", chainID)
}

// states fetches the chain ibc states
func (n Node) states(ctx context.Context) (ibctypes.IdentifiedClientStates, error) {
	res, err := n.ibcClientQuery.ClientStates(ctx, &ibctypes.QueryClientStatesRequest{})
	if err != nil {
		return ibctypes.IdentifiedClientStates{}, err
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
