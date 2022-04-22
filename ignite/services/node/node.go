package node

import (
	"context"
	"encoding/base64"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	commitmenttypes "github.com/cosmos/ibc-go/v2/modules/core/23-commitment/types"
	ibctmtypes "github.com/cosmos/ibc-go/v2/modules/light-clients/07-tendermint/types"
	spntypes "github.com/tendermint/spn/pkg/types"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/ignite-hq/cli/ignite/pkg/cosmosaccount"
	"github.com/ignite-hq/cli/ignite/pkg/cosmosclient"
)

type (
	//go:generate mockery --name CosmosClient --case underscore
	CosmosClient interface {
		Account(accountName string) (cosmosaccount.Account, error)
		Address(accountName string) (sdktypes.AccAddress, error)
		Context() client.Context
		BroadcastTx(accountName string, msgs ...sdktypes.Msg) (cosmosclient.Response, error)
		BroadcastTxWithProvision(accountName string, msgs ...sdktypes.Msg) (gas uint64, broadcast func() (cosmosclient.Response, error), err error)
		Status(ctx context.Context) (*ctypes.ResultStatus, error)
	}

	// Info is node client info.
	Info struct {
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

func New(
	ctx context.Context,
	nodeAPI string,
) (*Node, error) {
	c, err := cosmosclient.New(ctx, cosmosclient.WithNodeAddress(nodeAPI))
	if err != nil {
		return nil, err
	}
	return &Node{
		cosmos:       c,
		stakingQuery: stakingtypes.NewQueryClient(c.Context()),
	}, nil
}

func (n Node) Info(ctx context.Context) (Info, error) {
	// choose the last block height to fetch data on the launched blockchain
	status, err := n.cosmos.Status(ctx)
	if err != nil {
		return Info{}, err
	}
	lastBlockHeight := status.SyncInfo.LatestBlockHeight

	consensusState, err := n.ConsensusState(lastBlockHeight)
	if err != nil {
		return Info{}, err
	}
	spnConsensusStatue := spntypes.NewConsensusState(
		consensusState.Timestamp.Format(time.RFC3339Nano),
		consensusState.NextValidatorsHash.String(),
		base64.StdEncoding.EncodeToString(consensusState.Root.Hash),
	)

	header, err := n.TendermintHeader(lastBlockHeight)
	if err != nil {
		return Info{}, err
	}
	validators := make([]spntypes.Validator, len(header.ValidatorSet.Validators))
	for i, validator := range header.ValidatorSet.Validators {
		validators[i] = spntypes.NewValidator(
			validator.PubKey.String(),
			validator.ProposerPriority,
			validator.VotingPower,
		)
	}

	stakingParams, err := n.StakingParams(ctx)
	if err != nil {
		return Info{}, err
	}

	return Info{
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

// ConsensusState returns the appropriate tendermint consensus state by given height
func (n Node) ConsensusState(height int64) (*ibctmtypes.ConsensusState, error) {
	node, err := n.cosmos.Context().GetNode()
	if err != nil {
		return &ibctmtypes.ConsensusState{}, err
	}

	commit, err := node.Commit(context.Background(), &height)
	if err != nil {
		return &ibctmtypes.ConsensusState{}, err
	}

	var (
		page       = 1
		count      = 10_000
		nextHeight = height + 1
	)
	nextVals, err := node.Validators(context.Background(), &nextHeight, &page, &count)
	if err != nil {
		return &ibctmtypes.ConsensusState{}, err
	}

	state := &ibctmtypes.ConsensusState{
		Timestamp:          commit.Time,
		Root:               commitmenttypes.NewMerkleRoot(commit.AppHash),
		NextValidatorsHash: tmtypes.NewValidatorSet(nextVals.Validators).Hash(),
	}

	return state, nil
}

// TendermintHeader returns the appropriate tendermint header by given height
func (n Node) TendermintHeader(height int64) (ibctmtypes.Header, error) {
	node, err := n.cosmos.Context().GetNode()
	if err != nil {
		return ibctmtypes.Header{}, err
	}

	commit, err := node.Commit(context.Background(), &height)
	if err != nil {
		return ibctmtypes.Header{}, err
	}

	var (
		page  = 1
		count = 10_000
	)
	validators, err := node.Validators(context.Background(), &height, &page, &count)
	if err != nil {
		return ibctmtypes.Header{}, err
	}

	protoValset, err := tmtypes.NewValidatorSet(validators.Validators).ToProto()
	if err != nil {
		return ibctmtypes.Header{}, err
	}

	return ibctmtypes.Header{
		SignedHeader: commit.SignedHeader.ToProto(),
		ValidatorSet: protoValset,
	}, nil
}
