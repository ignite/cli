package networktypes_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	launchtypes "github.com/tendermint/spn/x/launch/types"

	"github.com/ignite/cli/ignite/services/network/networktypes"
)

func TestToChainLaunch(t *testing.T) {
	tests := []struct {
		name     string
		fetched  launchtypes.Chain
		expected networktypes.ChainLaunch
	}{
		{
			name: "chain with default genesis",
			fetched: launchtypes.Chain{
				LaunchID:       1,
				GenesisChainID: "foo-1",
				SourceURL:      "foo.com",
				SourceHash:     "0xaaa",
				HasCampaign:    true,
				CampaignID:     1,
				InitialGenesis: launchtypes.NewDefaultInitialGenesis(),
			},
			expected: networktypes.ChainLaunch{
				ID:              1,
				ChainID:         "foo-1",
				SourceURL:       "foo.com",
				SourceHash:      "0xaaa",
				GenesisURL:      "",
				GenesisHash:     "",
				LaunchTriggered: false,
				CampaignID:      1,
				Network:         "testnet",
			},
		},
		{
			name: "launched chain with custom genesis url and no campaign",
			fetched: launchtypes.Chain{
				LaunchID:        1,
				GenesisChainID:  "bar-1",
				SourceURL:       "bar.com",
				SourceHash:      "0xbbb",
				LaunchTriggered: true,
				LaunchTimestamp: 100,
				InitialGenesis: launchtypes.NewGenesisURL(
					"genesisfoo.com",
					"0xccc",
				),
			},
			expected: networktypes.ChainLaunch{
				ID:              1,
				ChainID:         "bar-1",
				SourceURL:       "bar.com",
				SourceHash:      "0xbbb",
				GenesisURL:      "genesisfoo.com",
				GenesisHash:     "0xccc",
				LaunchTriggered: true,
				LaunchTime:      100,
				CampaignID:      0,
				Network:         "testnet",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.EqualValues(t, tt.expected, networktypes.ToChainLaunch(tt.fetched))
		})
	}
}
