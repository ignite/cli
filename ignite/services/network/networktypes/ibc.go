package networktypes

import spntypes "github.com/tendermint/spn/pkg/types"

// IBCInfo is node client info.
type IBCInfo struct {
	ConsensusState spntypes.ConsensusState
	ValidatorSet   spntypes.ValidatorSet
	RevisionHeight uint64
}
