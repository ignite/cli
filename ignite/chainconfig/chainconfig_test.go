package chainconfig_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ignite/cli/ignite/chainconfig"
	"github.com/ignite/cli/ignite/chainconfig/config"
)

func TestCheckVersion(t *testing.T) {
	// Arrange
	cfg := bytes.NewBufferString(
		fmt.Sprintf("version: %d", chainconfig.LatestVersion),
	)

	// Act
	err := chainconfig.CheckVersion(cfg)

	// Assert
	require.NoError(t, err)
}

func TestCheckVersionWithOutdatedVersion(t *testing.T) {
	// Arrange
	cfg := bytes.NewBufferString("version: 0")
	wantError := chainconfig.VersionError{}

	// Act
	err := chainconfig.CheckVersion(cfg)

	// Assert
	require.ErrorAs(t, err, &wantError)
	require.Equal(t, wantError.Version, config.Version(0))
}
