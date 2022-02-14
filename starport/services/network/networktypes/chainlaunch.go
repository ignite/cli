package networktypes

import launchtypes "github.com/tendermint/spn/x/launch/types"

type (
	NetworkType string

	// ChainLaunch represents the launch of a chain on SPN
	ChainLaunch struct {
		ID          uint64      `json:"ID"`
		ChainID     string      `json:"ChainID"`
		SourceURL   string      `json:"SourceURL"`
		SourceHash  string      `json:"SourceHash"`
		GenesisURL  string      `json:"GenesisURL"`
		GenesisHash string      `json:"GenesisHash"`
		LaunchTime  int64       `json:"LaunchTime"`
		CampaignID  uint64      `json:"CampaignID"`
		Network     NetworkType `json:"Network"`
		Reward      string      `json:"Reward,omitempty"`
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
	var launchTime int64
	if chain.LaunchTriggered {
		launchTime = chain.LaunchTimestamp
	}

	network := NetworkTypeTestnet
	if chain.IsMainnet {
		network = NetworkTypeMainnet
	}

	launch := ChainLaunch{
		ID:         chain.LaunchID,
		ChainID:    chain.GenesisChainID,
		SourceURL:  chain.SourceURL,
		SourceHash: chain.SourceHash,
		LaunchTime: launchTime,
		CampaignID: chain.CampaignID,
		Network:    network,
	}

	// check if custom genesis URL is provided.
	if customGenesisURL := chain.InitialGenesis.GetGenesisURL(); customGenesisURL != nil {
		launch.GenesisURL = customGenesisURL.Url
		launch.GenesisHash = customGenesisURL.Hash
	}

	return launch
}
