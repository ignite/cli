package networktypes

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	yaml "github.com/goccy/go-yaml"
	launchtypes "github.com/tendermint/spn/x/launch/types"
)

type (
	NetworkType string

	// ChainLaunch represents the launch of a chain on SPN
	ChainLaunch struct {
		ID                     uint64      `json:"ID"`
		ConsumerRevisionHeight int64       `json:"ConsumerRevisionHeight"`
		ChainID                string      `json:"ChainID"`
		SourceURL              string      `json:"SourceURL"`
		SourceHash             string      `json:"SourceHash"`
		GenesisURL             string      `json:"GenesisURL"`
		GenesisHash            string      `json:"GenesisHash"`
		GenesisConfig          string      `json:"GenesisConfig"`
		LaunchTime             time.Time   `json:"LaunchTime"`
		ProjectID              uint64      `json:"ProjectID"`
		LaunchTriggered        bool        `json:"LaunchTriggered"`
		Network                NetworkType `json:"Network"`
		Reward                 string      `json:"Reward,omitempty"`
		AccountBalance         sdk.Coins   `json:"AccountBalance"`
		Metadata               interface{} `json:"Metadata"`
	}
)

const (
	NetworkTypeMainnet NetworkType = "mainnet"
	NetworkTypeTestnet NetworkType = "testnet"
)

func (n NetworkType) String() string {
	return string(n)
}

// ToChainLaunch converts a chain launch data from SPN and returns a ChainLaunch object
func ToChainLaunch(chain launchtypes.Chain) ChainLaunch {
	var launchTime time.Time
	if chain.LaunchTriggered {
		launchTime = chain.LaunchTime
	}

	network := NetworkTypeTestnet
	if chain.IsMainnet {
		network = NetworkTypeMainnet
	}

	launch := ChainLaunch{
		ID:                     chain.LaunchID,
		ConsumerRevisionHeight: chain.ConsumerRevisionHeight,
		ChainID:                chain.GenesisChainID,
		SourceURL:              chain.SourceURL,
		SourceHash:             chain.SourceHash,
		LaunchTime:             launchTime,
		ProjectID:              chain.CampaignID,
		LaunchTriggered:        chain.LaunchTriggered,
		Network:                network,
		AccountBalance:         chain.AccountBalance,
	}

	err := yaml.Unmarshal(chain.Metadata, &launch.Metadata)
	if err != nil {
		// an error shouldn't happen
		// in case it occurs, we consider metadata as invalid and dismiss those
		launch.Metadata = nil
	}

	// check if custom genesis URL is provided.
	if customGenesisURL := chain.InitialGenesis.GetGenesisURL(); customGenesisURL != nil {
		launch.GenesisURL = customGenesisURL.Url
		launch.GenesisHash = customGenesisURL.Hash
	}

	// check if custom config genesis if provided.
	if customGenesisConfig := chain.InitialGenesis.GetGenesisConfig(); customGenesisConfig != nil {
		launch.GenesisConfig = customGenesisConfig.File
	}

	return launch
}
