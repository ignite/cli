package cosmosver_test

import (
	"testing"

	"github.com/ignite/cli/v28/ignite/pkg/cosmosver"
	"github.com/stretchr/testify/require"
)

func TestDetect(t *testing.T) {
	_, err := cosmosver.Detect(".")
	require.Error(t, err)

	v, err := cosmosver.Detect("testdata/chain")
	require.NoError(t, err)
	require.Equal(t, "v0.47.3", v.Version)

	v, err = cosmosver.Detect("testdata/chain-sdk-fork")
	require.NoError(t, err)
	require.Equal(t, "", v)
	require.Equal(t, "v0.50.1-rollkit-v0.11.6-no-fraud-proofs", v.Version)
}
