package networktypes

import (
	spntypes "github.com/tendermint/spn/pkg/types"
)

// Reward is node reward info.
type Reward struct {
	ConsensusState spntypes.ConsensusState
	ValidatorSet   spntypes.ValidatorSet
	RevisionHeight uint64
}
