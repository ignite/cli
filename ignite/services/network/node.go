package network

import (
	"context"
	"encoding/base64"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	spntypes "github.com/tendermint/spn/pkg/types"
)

type (
	// IBCInfo is node client info.
	IBCInfo struct {
		ConsensusState spntypes.ConsensusState
		ValidatorSet   spntypes.ValidatorSet
		UnbondingTime  int64
		Height         uint64
	}

	// Node is node builder.
	Node struct {
		cosmos       CosmosClient
		stakingQuery stakingtypes.QueryClient
	}
)

func NewNode(cosmos CosmosClient) (*Node, error) {
	return &Node{
		cosmos:       cosmos,
		stakingQuery: stakingtypes.NewQueryClient(cosmos.Context()),
	}, nil
}

// IBCInfo Fetches the consensus state, validator set and the staking parameters
func (n Node) IBCInfo(ctx context.Context) (IBCInfo, error) {
	status, err := n.cosmos.Status(ctx)
	if err != nil {
		return IBCInfo{}, err
	}
	lastBlockHeight := status.SyncInfo.LatestBlockHeight

	consensusState, err := n.cosmos.IBCInfo(ctx, lastBlockHeight)
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
			base64.StdEncoding.EncodeToString(validator.PubKey.Bytes()),
			validator.ProposerPriority,
			validator.VotingPower,
		)
	}

	stakingParams, err := n.StakingParams(ctx)
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

// StakingParams fetches the staking module params
func (n Node) StakingParams(ctx context.Context) (stakingtypes.Params, error) {
	res, err := n.stakingQuery.Params(ctx, &stakingtypes.QueryParamsRequest{})
	if err != nil {
		return stakingtypes.Params{}, err
	}
	return res.Params, nil
}
