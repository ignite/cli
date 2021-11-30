package starportcmd_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	launchtypes "github.com/tendermint/spn/x/launch/types"
	starportcmd "github.com/tendermint/starport/starport/cmd"
)

func TestLaunchSummaries(t *testing.T) {
	tests := []struct {
		name     string
		chains   []launchtypes.Chain
		wantSums []starportcmd.LaunchSummary
	}{
		{
			name: "chain summaries",
			chains: []launchtypes.Chain{
				{
					LaunchID:       1,
					GenesisChainID: "foo-1",
					SourceURL:      "foo.com",
					HasCampaign:    true,
					CampaignID:     3,
				},
				{
					LaunchID:       2,
					GenesisChainID: "bar-1",
					SourceURL:      "bar.com",
					HasCampaign:    false,
				},
			},
			wantSums: []starportcmd.LaunchSummary{
				{
					LaunchID:   "1",
					ChainID:    "foo-1",
					Source:     "foo.com",
					CampaignID: "3",
				},
				{
					LaunchID:   "2",
					ChainID:    "bar-1",
					Source:     "bar.com",
					CampaignID: "no campaign",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSums := starportcmd.LaunchSummaries(tt.chains)
			require.Equal(t, tt.wantSums, gotSums)
		})
	}
}
