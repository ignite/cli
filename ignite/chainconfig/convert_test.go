package chainconfig_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite-hq/cli/ignite/chainconfig"
	v0testdata "github.com/ignite-hq/cli/ignite/chainconfig/v0/testdata"
)

func TestConvertLatest(t *testing.T) {
	// Arrange
	cfgV0 := v0testdata.GetConfigV0()

	// Act
	cfgLatest, err := chainconfig.ConvertLatest(cfgV0)

	// Assert
	require.NoError(t, err)
	require.Equal(t, chainconfig.LatestVersion, cfgLatest.GetVersion())
}
