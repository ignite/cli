package config_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/ignite/config"
	chainconfig "github.com/ignite/cli/ignite/config/chain"
	"github.com/ignite/cli/ignite/config/testdata"
)

func TestReadConfigVersion(t *testing.T) {
	// Arrange
	r := strings.NewReader("version: 42")
	want := chainconfig.Version(42)

	// Act
	version, err := config.ReadConfigVersion(r)

	// Assert
	require.NoError(t, err)
	require.Equal(t, want, version)
}

func TestParse(t *testing.T) {
	// Arrange: Initialize a reader with the previous version
	ver := config.LatestVersion - 1
	r := bytes.NewReader(testdata.Versions[ver])

	// Act
	cfg, err := config.Parse(r)

	// Assert
	require.NoError(t, err)

	// Assert: Parse must return the latest version
	require.Equal(t, config.LatestVersion, cfg.Version)
	require.Equal(t, testdata.GetLatestConfig(t), cfg)
}

func TestParseWithCurrentVersion(t *testing.T) {
	// Arrange
	r := bytes.NewReader(testdata.Versions[config.LatestVersion])

	// Act
	cfg, err := config.Parse(r)

	// Assert
	require.NoError(t, err)
	require.Equal(t, config.LatestVersion, cfg.Version)
	require.Equal(t, testdata.GetLatestConfig(t), cfg)
}

func TestParseWithUnknownVersion(t *testing.T) {
	// Arrange
	version := chainconfig.Version(9999)
	r := strings.NewReader(fmt.Sprintf("version: %d", version))

	var want *config.UnsupportedVersionError

	// Act
	_, err := config.Parse(r)

	// Assert
	require.ErrorAs(t, err, &want)
	require.NotNil(t, want)
	require.Equal(t, want.Version, version)
}

func TestParseNetworkWithCurrentVersion(t *testing.T) {
	// Arrange
	r := bytes.NewReader(testdata.NetworkConfig)

	// Act
	cfg, err := config.ParseNetwork(r)

	// Assert
	require.NoError(t, err)

	// Assert: Parse must return the latest version
	require.Equal(t, config.LatestVersion, cfg.Version)
	require.Equal(t, testdata.GetLatestNetworkConfig(t).Accounts, cfg.Accounts)
	require.Equal(t, testdata.GetLatestNetworkConfig(t).Genesis, cfg.Genesis)
}

func TestParseNetworkWithInvalidData(t *testing.T) {
	// Arrange
	r := bytes.NewReader(testdata.Versions[config.LatestVersion])

	// Act
	_, err := config.ParseNetwork(r)

	// Assert error
	require.True(
		t,
		strings.Contains(
			err.Error(),
			"config is not valid: no validators can be used in config for network genesis",
		),
	)
}
