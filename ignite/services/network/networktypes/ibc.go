package networktypes

import (
	ibcclienttypes "github.com/cosmos/ibc-go/v2/modules/core/02-client/types"
	ibcchanneltypes "github.com/cosmos/ibc-go/v2/modules/core/04-channel/types"
	lightclienttypes "github.com/cosmos/ibc-go/v2/modules/light-clients/07-tendermint/types"
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
		Channel      Channel
	}
	// ClientState is the relayer client state.
	ClientState struct {
		ClientID string
		ChainID  string
	}
	// Channel is the relayer connection channel.
	Channel struct {
		State                 string
		Ordering              string
		ConnectionHops        []string
		Version               string
		PortID                string
		ChannelID             string
		CounterpartyPortID    string
		CounterpartyChannelID string
	}
)

// ToClientState converts ibc client state
func ToClientState(state ibcclienttypes.IdentifiedClientState) (ClientState, bool) {
	clientState, ok := state.ClientState.GetCachedValue().(*lightclienttypes.ClientState)
	if !ok || clientState == nil {
		return ClientState{}, false
	}
	return ClientState{
		ClientID: state.ClientId,
		ChainID:  clientState.ChainId,
	}, true
}

// ToChannel converts ibc channel
func ToChannel(channel ibcchanneltypes.IdentifiedChannel) Channel {
	return Channel{
		State:                 ibcchanneltypes.State_name[int32(channel.State)],
		Ordering:              ibcchanneltypes.Order_name[int32(channel.Ordering)],
		ConnectionHops:        channel.ConnectionHops,
		Version:               channel.Version,
		PortID:                channel.PortId,
		ChannelID:             channel.ChannelId,
		CounterpartyPortID:    channel.Counterparty.PortId,
		CounterpartyChannelID: channel.Counterparty.ChannelId,
	}
}
