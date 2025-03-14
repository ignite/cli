package cosmosclient

import (
	"context"
	"encoding/base64"
	"time"

	"github.com/cometbft/cometbft/libs/bytes"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmtypes "github.com/cometbft/cometbft/types"
)

// ConsensusInfo is the validator consensus info.
type ConsensusInfo struct {
	Timestamp          string                 `json:"Timestamp"`
	Root               string                 `json:"Root"`
	NextValidatorsHash string                 `json:"NextValidatorsHash"`
	ValidatorSet       *cmtproto.ValidatorSet `json:"ValidatorSet"`
}

// ConsensusInfo returns the appropriate tendermint consensus state by given height
// and the validator set for the next height.
func (c Client) ConsensusInfo(ctx context.Context, height int64) (ConsensusInfo, error) {
	node, err := c.Context().GetNode()
	if err != nil {
		return ConsensusInfo{}, err
	}

	commit, err := node.Commit(ctx, &height)
	if err != nil {
		return ConsensusInfo{}, err
	}

	var (
		page  = 1
		count = 10_000
	)
	validators, err := node.Validators(ctx, &height, &page, &count)
	if err != nil {
		return ConsensusInfo{}, err
	}

	protoValset, err := tmtypes.NewValidatorSet(validators.Validators).ToProto()
	if err != nil {
		return ConsensusInfo{}, err
	}

	heightNext := height + 1
	validatorsNext, err := node.Validators(ctx, &heightNext, &page, &count)
	if err != nil {
		return ConsensusInfo{}, err
	}

	hash := tmtypes.NewValidatorSet(validatorsNext.Validators).Hash()

	return ConsensusInfo{
		Timestamp:          commit.Time.Format(time.RFC3339Nano),
		NextValidatorsHash: bytes.HexBytes(hash).String(),
		Root:               base64.StdEncoding.EncodeToString(commit.AppHash),
		ValidatorSet:       protoValset,
	}, nil
}
