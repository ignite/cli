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

	// RewardIBCInfo holds IBC info to relay packets to claim rewards.
	RewardIBCInfo struct {
		ChainID      string
		ClientID     string
		ConnectionID string
		ChannelID    string
	}
)
