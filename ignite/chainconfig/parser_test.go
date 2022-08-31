package chainconfig_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/ignite/chainconfig"
	"github.com/ignite/cli/ignite/chainconfig/config"
	"github.com/ignite/cli/ignite/chainconfig/testdata"
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

func TestParseWithUnknownVersion(t *testing.T) {
	// Arrange
	version := config.Version(9999)
	r := strings.NewReader(fmt.Sprintf("version: %d", version))

	var want *chainconfig.UnsupportedVersionError

	// Act
	_, err := chainconfig.Parse(r)

	// Assert
	require.ErrorAs(t, err, &want)
	require.NotNil(t, want)
	require.Equal(t, want.Version, version)
}
