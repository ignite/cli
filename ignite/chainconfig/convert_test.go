package chainconfig_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite-hq/cli/ignite/chainconfig"
	"github.com/ignite-hq/cli/ignite/chainconfig/testdata"
)

func TestConvertLatest(t *testing.T) {
	// Arrange
	cfgV0 := testdata.GetLatestConfig(t)

	// Act
	cfgLatest, err := chainconfig.ConvertLatest(cfgV0)

	// Assert
	require.NoError(t, err)
	require.Equal(t, chainconfig.LatestVersion, cfgLatest.GetVersion())
}

func TestMigrateLatest(t *testing.T) {
	// Arrange
	current := bytes.NewReader(testdata.Versions[chainconfig.LatestVersion-1])
	latest := bytes.Buffer{}
	want := string(testdata.Versions[chainconfig.LatestVersion])

	// Act
	err := chainconfig.MigrateLatest(current, &latest)

	// Assert
	require.NotEmpty(t, want, "testdata is missing the latest config version")
	require.NoError(t, err)
	require.Equal(t, want, latest.String())
}
