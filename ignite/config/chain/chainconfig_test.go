package chain_test

import (
	"bytes"
	"fmt"
	"github.com/ignite/cli/ignite/config/chain"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCheckVersion(t *testing.T) {
	// Arrange
	cfg := bytes.NewBufferString(
		fmt.Sprintf("version: %d", LatestVersion),
	)

	// Act
	err := CheckVersion(cfg)

	// Assert
	require.NoError(t, err)
}

func TestCheckVersionWithOutdatedVersion(t *testing.T) {
	// Arrange
	cfg := bytes.NewBufferString("version: 0")
	wantError := chain.VersionError{}

	// Act
	err := CheckVersion(cfg)

	// Assert
	require.ErrorAs(t, err, &wantError)
	require.Equal(t, wantError.Version, Version(0))
}
