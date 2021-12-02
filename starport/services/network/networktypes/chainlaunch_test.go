package networktypes_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/spn/x/launch/types"
	launchtypes "github.com/tendermint/spn/x/launch/types"
	"github.com/tendermint/starport/starport/services/network/networktypes"
)

func TestParseChainLaunch(t *testing.T) {
	type args struct {
		chain types.Chain
	}
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
				ID:          1,
				ChainID:     "foo-1",
				SourceURL:   "foo.com",
				SourceHash:  "0xaaa",
				CampaignID:  1,
				GenesisURL:  "",
				GenesisHash: "",
			},
		},
		{
			name: "chain with custom genesis url and no campaign",
			fetched: launchtypes.Chain{
				LaunchID:       1,
				GenesisChainID: "bar-1",
				SourceURL:      "bar.com",
				SourceHash:     "0xbbb",
				InitialGenesis: launchtypes.NewGenesisURL(
					"genesisfoo.com",
					"0xccc",
				),
			},
			expected: networktypes.ChainLaunch{
				ID:          1,
				ChainID:     "bar-1",
				SourceURL:   "bar.com",
				SourceHash:  "0xbbb",
				CampaignID:  0,
				GenesisURL:  "genesisfoo.com",
				GenesisHash: "0xccc",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.EqualValues(t, tt.expected, networktypes.ParseChainLaunch(tt.fetched))
		})
	}
}
