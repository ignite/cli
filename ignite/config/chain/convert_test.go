package chain_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/ignite/config/chain"
	v0testdata "github.com/ignite/cli/ignite/config/chain/v0/testdata"
	"github.com/ignite/cli/ignite/config/testdata"
)

func TestConvertLatest(t *testing.T) {
	// Arrange
	cfgV0 := v0testdata.GetConfig(t)

	// Act
	cfgLatest, err := chain.ConvertLatest(cfgV0)

	// Assert
	require.NoError(t, err)
	require.Equal(t, chain.LatestVersion, cfgLatest.GetVersion())
	require.Equal(t, testdata.GetLatestConfig(t), cfgLatest)
}

func TestMigrateLatest(t *testing.T) {
	// Arrange
	current := bytes.NewReader(testdata.Versions[chain.LatestVersion-1])
	latest := bytes.Buffer{}
	want := string(testdata.Versions[chain.LatestVersion])

	// Act
	err := chain.MigrateLatest(current, &latest)

	// Assert
	require.NotEmpty(t, want, "testdata is missing the latest config version")
	require.NoError(t, err)
	require.Equal(t, want, latest.String())
}
