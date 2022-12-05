package chain_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/ignite/config/chain"
	"github.com/ignite/cli/ignite/config/chain/version"
)

func TestCheckVersion(t *testing.T) {
	// Arrange
	cfg := bytes.NewBufferString(
		fmt.Sprintf("version: %d", chain.LatestVersion),
	)

	// Act
	err := chain.CheckVersion(cfg)

	// Assert
	require.NoError(t, err)
}

func TestCheckVersionWithOutdatedVersion(t *testing.T) {
	// Arrange
	cfg := bytes.NewBufferString("version: 0")
	wantError := chain.VersionError{}

	// Act
	err := chain.CheckVersion(cfg)

	// Assert
	require.ErrorAs(t, err, &wantError)
	require.Equal(t, wantError.Version, version.Version(0))
}
