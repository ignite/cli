package chainconfig_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/ignite-hq/cli/ignite/chainconfig"
	"github.com/ignite-hq/cli/ignite/chainconfig/config"
	"github.com/ignite-hq/cli/ignite/chainconfig/testdata"
	"github.com/stretchr/testify/require"
)

func TestReadConfigVersion(t *testing.T) {
	// Arrange
	r := strings.NewReader("version: 42")
	want := config.Version(42)

	// Act
	version, err := chainconfig.ReadConfigVersion(r)

	// Assert
	require.NoError(t, err)
	require.Equal(t, want, version)
}

func TestParse(t *testing.T) {
	// Arrange: Initialize a reader with the previous version
	ver := chainconfig.LatestVersion - 1
	r := bytes.NewReader(testdata.Versions[ver])

	// Act
	cfg, err := chainconfig.Parse(r)

	// Assert
	require.NoError(t, err)

	// Assert: Parse must return the latest version
	require.Equal(t, chainconfig.LatestVersion, cfg.Version)
	require.EqualValues(t, testdata.GetLatestConfig(t), cfg)
}

func TestParseWithCurrentVersion(t *testing.T) {
	// Arrange
	r := bytes.NewReader(testdata.Versions[chainconfig.LatestVersion])

	// Act
	cfg, err := chainconfig.Parse(r)

	// Assert
	require.NoError(t, err)
	require.Equal(t, chainconfig.LatestVersion, cfg.Version)
	require.EqualValues(t, testdata.GetLatestConfig(t), cfg)
}
