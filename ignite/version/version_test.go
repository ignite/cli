package version_test

import (
	"testing"

	"github.com/blang/semver/v4"
	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/v29/ignite/pkg/cosmosver"
	"github.com/ignite/cli/v29/ignite/version"
)

func TestAssertSupportedCosmosSDKVersion(t *testing.T) {
	testCases := []struct {
		name    string
		version cosmosver.Version
		errMsg  string
	}{
		{
			"invalid",
			cosmosver.Version{Version: "invalid"},
			"Your chain has been scaffolded with an older version of Cosmos SDK: invalid",
		},
		{
			"too old",
			cosmosver.Version{Version: "v0.45.0", Semantic: semver.MustParse("0.45.0")},
			"Your chain has been scaffolded with an older version of Cosmos SDK: v0.45.0",
		},
		{
			"v0.47.3",
			cosmosver.Version{Version: "v0.47.3", Semantic: semver.MustParse("0.47.3")},
			"Your chain has been scaffolded with an older version of Cosmos SDK: v0.47.3",
		},
		{
			"v0.50",
			cosmosver.Version{Version: "v0.50.1", Semantic: semver.MustParse("0.50.1")},
			"",
		},
		{
			"v0.50 fork",
			cosmosver.Version{Version: "v0.50.1-rollkit-v0.11.6-no-fraud-proofs", Semantic: semver.MustParse("0.50.1-rollkit-v0.11.6-no-fraud-proofs")},
			"",
		},
		{
			"v0.53",
			cosmosver.Version{Version: "v0.53.0", Semantic: semver.MustParse("0.53.0")},
			"",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := version.AssertSupportedCosmosSDKVersion(tc.version)
			if tc.errMsg == "" {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tc.errMsg)
			}
		})
	}
}
