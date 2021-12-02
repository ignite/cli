package networktypes

import launchtypes "github.com/tendermint/spn/x/launch/types"

// ChainLaunch represents the launch of a chain on SPN
type ChainLaunch struct {
	ID          uint64
	ChainID     string
	SourceURL   string
	SourceHash  string
	GenesisURL  string
	GenesisHash string
	CampaignID  uint64
}

// ParseChainLaunch parses a chain launch data from SPN and returns a ChainLaunch object
func ParseChainLaunch(chain launchtypes.Chain) ChainLaunch {
	launch := ChainLaunch{
		ID:         chain.LaunchID,
		ChainID:    chain.GenesisChainID,
		SourceURL:  chain.SourceURL,
		SourceHash: chain.SourceHash,
		CampaignID: chain.CampaignID,
	}

	// check if custom genesis URL is provided.
	if customGenesisURL := chain.InitialGenesis.GetGenesisURL(); customGenesisURL != nil {
		launch.GenesisURL = customGenesisURL.Url
		launch.GenesisHash = customGenesisURL.Hash
	}

	return launch
}
