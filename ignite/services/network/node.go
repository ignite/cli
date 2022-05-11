package network

import (
	"context"
	"encoding/base64"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	spntypes "github.com/tendermint/spn/pkg/types"

	"github.com/ignite-hq/cli/ignite/services/network/networktypes"
)

type (
	// Node is node builder.
	Node struct {
		cosmos       CosmosClient
		stakingQuery stakingtypes.QueryClient
	}
)

func NewNodeClient(cosmos CosmosClient) (Node, error) {
	return Node{
		cosmos:       cosmos,
		stakingQuery: stakingtypes.NewQueryClient(cosmos.Context()),
	}, nil
}

// NodeInfo Fetches the consensus state with the validator set
func NodeInfo(ctx context.Context, client CosmosClient) (networktypes.IBCInfo, error) {
	status, err := client.Status(ctx)
	if err != nil {
		return networktypes.IBCInfo{}, err
	}
	lastBlockHeight := status.SyncInfo.LatestBlockHeight

	consensusState, err := client.ConsensusInfo(ctx, lastBlockHeight)
	if err != nil {
		return networktypes.IBCInfo{}, err
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

	return networktypes.IBCInfo{
		ConsensusState: spnConsensusStatue,
		ValidatorSet:   spntypes.NewValidatorSet(validators...),
		RevisionHeight: uint64(lastBlockHeight),
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

// IBCInfo Fetches the consensus state with the validator set and the unbounding time
func (n Node) IBCInfo(ctx context.Context) (networktypes.IBCInfo, int64, error) {
	info, err := NodeInfo(ctx, n.cosmos)
	if err != nil {
		return networktypes.IBCInfo{}, 0, err
	}

	stakingParams, err := n.StakingParams(ctx)
	if err != nil {
		return networktypes.IBCInfo{}, 0, err
	}
	return info, int64(stakingParams.UnbondingTime.Seconds()), nil
}
