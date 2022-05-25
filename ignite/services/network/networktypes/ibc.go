package networktypes

import (
	spntypes "github.com/tendermint/spn/pkg/types"
)

type (
	// Reward is node reward info.
	Reward struct {
		ConsensusState spntypes.ConsensusState
		ValidatorSet   spntypes.ValidatorSet
		RevisionHeight uint64
	}
	// Relayer is the relayer connection info.
	Relayer struct {
		ChainID      string
		ClientID     string
		ConnectionID string
		ChannelID    string
	}
)
